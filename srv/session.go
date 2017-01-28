package srv

import (
	"log"
	"time"
	"strconv"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

type session struct {
	ready	chan int
	action	chan *reqAction
	uids	[4]uid
	readys	[4]bool
	tables	*tables
	nonce	int
	timer	*time.Timer
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

	return s
}

func (s *session) Loop() {
	s.notifyLoad()

	table := saki.NewTableSession()
	defer saki.DeleteTableSession(table)

	for !table.GameOver() {
		select {
		case i := <-s.ready:
			if !(0 <= i && i < 4) {
				log.Fatalln("session.loop() i", i)
			}
			s.readys[i] = true;
			if s.readys[0] && s.readys[1] && s.readys[2] && s.readys[3] {
				s.start(table)
			}
		case act:= <-s.action:
			s.doAction(table, act)
		case <-s.timer.C:
			s.sweep(table)
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

func (s *session) notifyLoad() {
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
	}{"start", users, [4]int{0,0,0,0}, 0}

	for i, uid := range s.uids {
		msg.TempDealer = (4 - i) % 4
		s.tables.conns.peer <- &Mail{uid, msg}
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
	mails := table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	s.sendMail(mails, table)
}

func (s *session) sweep(table saki.TableSession) {
	mails := table.Sweep()
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
			time.Sleep(300 * time.Millisecond)
			act := reqAction{s.uids[toWhom], s.nonce, "SPIN_OUT", "-1"}
			s.doAction(table, &act)
		} else {
			msg = `{"Nonce":` + strconv.Itoa(s.nonce) + "," + msg[1:]
			mail := Mail{s.uids[toWhom], msg}
			s.tables.conns.peer <- &mail
		}
	}
}

