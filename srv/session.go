package srv

import (
	"log"
	"time"
	"strconv"
	"strings"
	"math/rand"
	"errors"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type session struct {
	ready		chan int
	action		chan *reqAction
	uids		[4]uid
	readys		[4]bool
	onlines		[4]bool
	tables		*tables
	nonce		int
	timer		*time.Timer
	timeOutCts	[4]int
}

func newSession(tables *tables, uids [4]uid) *session {
	s := new(session)

	s.ready = make(chan int)
	s.action = make(chan *reqAction)
	s.uids = uids
	s.tables = tables
	s.nonce = 0
	s.timer = time.NewTimer(1 * time.Second)
	if !s.timer.Stop() {
		select {
		case <-s.timer.C:
		default:
		}
	}
	for i := 0; i < 4; i++ {
		s.onlines[i] = true // regard as good by default
	}

	return s
}

func (s *session) Loop() {
	girlIds := genIds()
	s.notifyLoad(&girlIds)

	readyTimer := time.NewTimer(7 * time.Second)
	hardTimer := time.NewTimer(2 * time.Hour)

	table := saki.NewTableSession(
		girlIds[0], girlIds[1], girlIds[2], girlIds[3])
	defer saki.DeleteTableSession(table)

	for s.anyOnline() && !table.GameOver() {
		select {
		case i := <-s.ready:
			if !(0 <= i && i < 4) {
				log.Fatalln("session.loop() i", i)
			}
			s.readys[i] = true;
			if s.allReady() {
				s.start(table)
			}
		case act:= <-s.action:
			s.doAction(table, act)
		case <-s.timer.C:
			s.sweepAll(table)
		case <-readyTimer.C:
			if !s.allReady() {
				for i := 0; i < 4; i++ {
					s.readys[i] = true
				}
				s.start(table)
			}
		case <-hardTimer.C:
			break
		}
	}

	s.tables.EndSession() <- s
}

func (s *session) Ready() chan<- int {
	return s.ready
}

func (s *session) Action() chan<- *reqAction {
	return s.action
}

func (s *session) allReady() bool {
	return s.readys[0] && s.readys[1] && s.readys[2] && s.readys[3]
}

func (s *session) anyOnline() bool {
	ol := &s.onlines
	return ol[0] || ol[1] || ol[2] || ol[3]
}

func (s *session) notifyLoad(girlIds *[4]int) {
	var users [4]*user
	for i := range users {
		users[i] = s.tables.conns.dao.getUser(s.uids[i])
		if users[i] == nil {
			log.Fatal("session.nofityLoad:", s.uids[i], "not in DB")
		}
	}

	msg := struct {
		Type		string
		Users		[4]*user
		GirlIds		[4]int
		TempDealer	int
	}{"start", users, *girlIds, 0}

	for i, uid := range s.uids {
		msg.TempDealer = (4 - i) % 4
		s.sendPeer(i, &Mail{uid, msg})
		// rotate perspectives
		u0 := msg.Users[0]
		msg.Users[0] = msg.Users[1]
		msg.Users[1] = msg.Users[2]
		msg.Users[2] = msg.Users[3]
		msg.Users[3] = u0

		g0 := msg.GirlIds[0]
		msg.GirlIds[0] = msg.GirlIds[1]
		msg.GirlIds[1] = msg.GirlIds[2]
		msg.GirlIds[2] = msg.GirlIds[3]
		msg.GirlIds[3] = g0
	}
}

func (s *session) FindUser(uid uid) (int, bool) {
	for i, u := range s.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (s *session) start(table saki.TableSession) {
	mails := table.Start()
	defer saki.DeleteMailVector(mails)
	s.sendMail(mails, table)
}

func (s *session) doAction(table saki.TableSession, act *reqAction) {
	if act.Nonce != s.nonce {
		log.Println("expired nonce", act.Nonce, "by", act.uid);
		return
	}

	i, _ := s.FindUser(act.uid)
	s.timeOutCts[i] = 0
	mails := table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	s.sendMail(mails, table)
}

func (s *session) sweepOne(table saki.TableSession, i int) {
	mails := table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	s.sendMail(mails, table)
}

func (s *session) sweepAll(table saki.TableSession) {
	var targets int
	mails := table.SweepAll(&targets)
	for w := uint(0); w < 4; w++ {
		if (targets & (1 << w)) != 0 {
			s.timeOutCts[w]++
			if s.timeOutCts[w] >= 3 {
				s.onlines[w] = false
				s.tables.conns.Logout() <- s.uids[w]
			}
		}
	}
	defer saki.DeleteMailVector(mails)
	s.sendMail(mails, table)
}

func (s *session) sendMail(mails saki.MailVector, table saki.TableSession) {
	size := int(mails.Size())
	if size > 0 {
		s.nonce++
		if !s.timer.Stop() {
			select {
			case <-s.timer.C:
			default:
			}
		}
		s.timer.Reset(7 * time.Second)
	}

	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		msg := mails.Get(i).GetMsg()
		if msg == "auto" { // special mark
			time.Sleep(500 * time.Millisecond)
			act := reqAction{s.uids[toWhom], s.nonce, "SPIN_OUT", "-1"}
			s.doAction(table, &act)
		} else {
			msg = `{"Nonce":` + strconv.Itoa(s.nonce) + "," + msg[1:]
			mail := Mail{s.uids[toWhom], msg}
			err := s.sendPeer(toWhom, &mail)
			if err != nil && strings.Contains(msg, "t-activated") {
				if s.anyOnline() && !table.GameOver() {
					s.sweepOne(table, toWhom)
				}
			}
		}
	}
}

func (s *session) sendPeer(i int, mail *Mail) error {
	if s.onlines[i] {
		if _, ok := s.tables.conns.conns[mail.To]; ok {
			s.tables.conns.Peer() <- mail
			return nil
		}
		s.onlines[i] = false;
	}
	return errors.New("peer not in session")
}

func genIds() [4]int {
	avails := []int{
		710113, 710114, 710115,
		712411, 712412,
		712611,
		712714, 712715,
		712915,
		713311, 713314,
		713811, 713815,
		714915,
		713301,
		990001, 990002}
	for {
		i0 := rand.Intn(len(avails))
		i1 := rand.Intn(len(avails))
		if i1 == i0 {
			continue
		}
		i2 := rand.Intn(len(avails))
		if i2 == i0 || i2 == i1 {
			continue
		}
		i3 := rand.Intn(len(avails))
		if i3 == i0 || i3 == i1 || i3 == i2 {
			continue
		}
		return [4]int{avails[i0], avails[i1], avails[i2], avails[i3]}
	}
}

