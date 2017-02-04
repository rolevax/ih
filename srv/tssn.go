package srv

import (
	"log"
	"strconv"
	"strings"
	"time"
	"math/rand"
	"errors"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tssn struct {
	ready		chan uid
	action		chan *msgTssnAction
	uids		[4]uid
	readys		[4]bool
	onlines		[4]bool
	nonce		int
	timer		*time.Timer
	timeOutCts	[4]int
}

func loopTssn(uids [4]uid) {
	tssn := startTssn(uids)
	sing.TssnMgr.Reg(tssn)
	defer sing.TssnMgr.Unreg(tssn)

	girlIds := genIds()
	table := saki.NewTableSession(
		girlIds[0], girlIds[1], girlIds[2], girlIds[3])
	defer saki.DeleteTableSession(table)

	tssn.notifyLoad(&girlIds, table)

	readyTimer := time.NewTimer(7 * time.Second)
	hardTimer := time.NewTimer(2 * time.Hour)

	for tssn.anyOnline() && !table.GameOver() {
		select {
		case uid := <-tssn.ready:
			tssn.handleReady(uid, table)
		case mta := <-tssn.action:
			tssn.handleAction(mta.uid, mta.act, table)
		case <-tssn.timer.C:
			tssn.sweepAll(table)
		case <-readyTimer.C:
			if !tssn.allReady() {
				for i := 0; i < 4; i++ {
					tssn.readys[i] = true
				}
				tssn.start(table)
			}
		case <-hardTimer.C:
			break
		}
	}
}

func startTssn(uids [4]uid) *tssn {
	s := new(tssn)

	s.ready = make(chan uid)
	s.action = make(chan *msgTssnAction)
	s.uids = uids
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

func (tssn *tssn) Ready(uid uid) {
	tssn.ready <- uid
}

type msgTssnAction struct {
	uid		uid
	act		*reqAction
}

func (tssn *tssn) Action(uid uid, act *reqAction) {
	tssn.action <- &msgTssnAction{uid, act}
}

func (tssn *tssn) handleReady(uid uid, table saki.TableSession) {
	if i, ok := tssn.findUser(uid); ok {
		if !tssn.readys[i] { // trigger only on toggle
			tssn.readys[i] = true;
			if tssn.allReady() {
				tssn.start(table)
			}
		}
	} else {
		log.Fatalln("loopTssn uid not found")
	}
}

func (tssn *tssn) allReady() bool {
	r := &tssn.readys
	return r[0] && r[1] && r[2] && r[3]
}

func (tssn *tssn) anyOnline() bool {
	o := &tssn.onlines
	return o[0] || o[1] || o[2] || o[3]
}

func (tssn *tssn) notifyLoad(girlIds *[4]int, table saki.TableSession) {
	var users [4]*user
	for i := range users {
		users[i] = sing.Dao.GetUser(tssn.uids[i])
		if users[i] == nil {
			log.Fatalln("tssn.nofityLoad:", tssn.uids[i], "not in DB")
		}
	}

	msg := struct {
		Type		string
		Users		[4]*user
		GirlIds		[4]int
		TempDealer	int
	}{"start", users, *girlIds, 0}

	for i, uid := range tssn.uids {
		msg.TempDealer = (4 - i) % 4
		err := tssn.sendPeer(i, msg)
		if err != nil {
			tssn.handleReady(uid, table)
		}

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

func (tssn *tssn) findUser(uid uid) (int, bool) {
	for i, u := range tssn.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (tssn *tssn) start(table saki.TableSession) {
	mails := table.Start()
	defer saki.DeleteMailVector(mails)
	tssn.sendMails(mails, table)
}

func (tssn *tssn) handleAction(uid uid, act *reqAction,
							   table saki.TableSession, ) {
	if act.Nonce != tssn.nonce {
		log.Println("expired nonce", act.Nonce, "by", uid);
		return
	}

	i, _ := tssn.findUser(uid)
	tssn.timeOutCts[i] = 0
	mails := table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	tssn.sendMails(mails, table)
}

func (tssn *tssn) sweepOne(table saki.TableSession, i int) {
	mails := table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	tssn.sendMails(mails, table)
}

func (tssn *tssn) sweepAll(table saki.TableSession) {
	var targets int
	mails := table.SweepAll(&targets)
	for w := uint(0); w < 4; w++ {
		if (targets & (1 << w)) != 0 {
			tssn.timeOutCts[w]++
			if tssn.timeOutCts[w] >= 3 {
				tssn.onlines[w] = false
				sing.UssnMgr.Logout(tssn.uids[w])
			}
		}
	}
	defer saki.DeleteMailVector(mails)
	tssn.sendMails(mails, table)
}

func (tssn *tssn) sendMails(mails saki.MailVector,
							table saki.TableSession) {
	size := int(mails.Size())
	if size > 0 {
		tssn.nonce++
		if !tssn.timer.Stop() {
			select {
			case <-tssn.timer.C:
			default:
			}
		}
		tssn.timer.Reset(7 * time.Second)
	}

	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		msg := mails.Get(i).GetMsg()
		if msg == "auto" { // special mark
			time.Sleep(500 * time.Millisecond)
			act := reqAction{tssn.nonce, "SPIN_OUT", "-1"}
			tssn.handleAction(tssn.uids[toWhom], &act, table)
		} else {
			msg = `{"Nonce":` + strconv.Itoa(tssn.nonce) + "," + msg[1:]
			err := tssn.sendPeer(toWhom, msg)
			if err != nil && strings.Contains(msg, "t-activated") {
				if tssn.anyOnline() && !table.GameOver() {
					tssn.sweepOne(table, toWhom)
				}
			}
		}
	}
}

func (tssn *tssn) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := sing.UssnMgr.Peer(tssn.uids[i], msg)
		if err != nil {
			tssn.onlines[i] = false
		}
		return err
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

