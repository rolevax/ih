package srv

import (
	"log"
	"time"
	"errors"
	"encoding/json"
	"math/rand"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

const actTimeOut = 12 * time.Second
const readyTimeOut = 20 * time.Second

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tssn struct {
	ready		chan uid
	action		chan *msgTssnAction
	done		chan struct{}
	uids		[4]uid
	readys		[4]bool
	onlines		[4]bool
	nonces		[4]int
	actTimer	*time.Timer
	timeOutCts	[4]int
}

func loopTssn(uids [4]uid) {
	tssn := startTssn(uids)
	defer close(tssn.done)
	sing.TssnMgr.Reg(tssn)
	defer sing.TssnMgr.Unreg(tssn)

	girlIds := genIds()
	table := saki.NewTableSession(
		girlIds[0], girlIds[1], girlIds[2], girlIds[3])
	defer saki.DeleteTableSession(table)

	tssn.notifyLoad(&girlIds, table)

	readyTimer := time.NewTimer(readyTimeOut)
	hardTimer := time.NewTimer(2 * time.Hour)
	defer hardTimer.Stop()

	for tssn.anyOnline() && !table.GameOver() {
		select {
		case uid := <-tssn.ready:
			tssn.handleReady(uid, table)
		case mta := <-tssn.action:
			tssn.handleAction(mta.uid, mta.act, table)
		case <-tssn.actTimer.C:
			tssn.sweepAll(table)
		case <-readyTimer.C:
			if !tssn.allReady() {
				for i := 0; i < 4; i++ {
					tssn.readys[i] = true
				}
				tssn.start(table)
			}
		case <-hardTimer.C:
			return
		}
	}
}

func startTssn(uids [4]uid) *tssn {
	tssn := new(tssn)

	tssn.ready = make(chan uid)
	tssn.action = make(chan *msgTssnAction)
	tssn.done = make(chan struct{})
	tssn.uids = uids
	tssn.actTimer = time.NewTimer(actTimeOut)
	if !tssn.actTimer.Stop() {
		select {
		case <-tssn.actTimer.C:
		default:
		}
	}
	for i := 0; i < 4; i++ {
		tssn.onlines[i] = true // regard as good by default
	}

	return tssn
}

func (tssn *tssn) Ready(uid uid) {
	select {
	case tssn.ready <- uid:
	case <-tssn.done:
	}
}

type msgTssnAction struct {
	uid		uid
	act		*reqAction
}

func (tssn *tssn) Action(uid uid, act *reqAction) {
	select {
	case tssn.action <- &msgTssnAction{uid, act}:
	case <-tssn.done:
	}
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
	users := sing.Dao.GetUsers(&tssn.uids)
	for i, user := range users {
		if user == nil {
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
	tssn.handleMails(mails, table)
}

func (tssn *tssn) handleAction(uid uid, act *reqAction,
							   table saki.TableSession, ) {
	i, _ := tssn.findUser(uid)
	if act.Nonce != tssn.nonces[i] {
		log.Println("expired nonce", act.Nonce, "by", uid);
		return
	}
	tssn.timeOutCts[i] = 0
	mails := table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails, table)
}

func (tssn *tssn) sweepOne(table saki.TableSession, i int) {
	mails := table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails, table)
}

func (tssn *tssn) sweepAll(table saki.TableSession) {
	var targets int
	mails := table.SweepAll(&targets)
	for w := uint(0); w < 4; w++ {
		if (targets & (1 << w)) != 0 {
			tssn.timeOutCts[w]++
			if tssn.timeOutCts[w] == 3 {
				tssn.timeOutCts[w] = 0
				tssn.onlines[w] = false
				sing.UssnMgr.Logout(tssn.uids[w])
			}
		}
	}
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails, table)
}

func (tssn *tssn) handleMails(mails saki.MailVector,
							  table saki.TableSession) {
	size := int(mails.Size())
	if size > 0 {
		var nonceInced [4]bool
		for i := 0; i < size; i++ {
			w := mails.Get(i).GetTo()
			if w != -1 && !nonceInced[w] {
				tssn.nonces[w]++
				nonceInced[w] = true
			}
		}
		tssn.resetActTimer()
	}

	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		str := mails.Get(i).GetMsg()
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(str), &msg); err != nil {
			log.Fatalln("unmarshal c++ str", err)
		}
		if toWhom == -1 {
			tssn.handleSystemMail(msg, table)
		} else {
			tssn.sendUserMail(toWhom, msg, table)
		}
	}
}

func (tssn *tssn) sendUserMail(who int, msg map[string]interface{},
							   table saki.TableSession) {
	msg["Nonce"] = tssn.nonces[who]

	err := tssn.sendPeer(who, msg)
	if err != nil && msg["Event"] == "activated" {
		if tssn.anyOnline() && !table.GameOver() {
			tssn.sweepOne(table, who)
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

func (tssn *tssn) resetActTimer() {
	if !tssn.actTimer.Stop() {
		select {
		case <-tssn.actTimer.C:
		default: // prevent blocked by double-draining
		}
	}
	tssn.actTimer.Reset(actTimeOut)
}

func (tssn *tssn) handleSystemMail(msg map[string]interface{},
						           table saki.TableSession) {
	switch (msg["Type"]) {
	case "table-end-stat":
		var ordered [4]uid
		ranks := msg["Rank"].([]interface{})
		for r := 0; r < 4; r++ {
			ordered[r] = tssn.uids[int(ranks[r].(float64))]
		}
		statRank(&ordered)
	case "riichi-auto":
		time.Sleep(500 * time.Millisecond)
		who := int(msg["Who"].(float64))
		act := reqAction{tssn.nonces[who], "SPIN_OUT", "-1"}
		tssn.handleAction(tssn.uids[who], &act, table)
	default:
		log.Fatalln("unknown system mail", msg)
	}
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
		715212,
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

