package tssn

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/mjpancake/hisa/actor/tssn/tbus"
	"github.com/mjpancake/hisa/actor/ussn/ubus"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/saki"
)

const answerTimeOut = 15 * time.Second

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tssnState int

const (
	tssnWaitChoose tssnState = iota
	tssnWaitReady
	tssnWaitAction
)

type ss struct {
	bookType    model.BookType
	state       tssnState
	ready       chan model.Uid
	choose      chan *tbus.MsgChoose
	action      chan *tbus.MsgAction
	done        chan struct{}
	uids        [4]model.Uid
	users       [4]*model.User
	gids        [4]model.Gid
	gidcs       [12]model.Gid
	waits       [4]bool
	onlines     [4]bool
	nonces      [4]int
	answerTimer *time.Timer
	table       saki.TableSession
}

func startTssn(bt model.BookType, uids [4]model.Uid) *ss {
	tssn := new(ss)

	tssn.bookType = bt
	tssn.state = tssnWaitChoose
	tssn.ready = make(chan model.Uid)
	tssn.choose = make(chan *tbus.MsgChoose)
	tssn.action = make(chan *tbus.MsgAction)
	tssn.done = make(chan struct{})
	tssn.uids = uids
	tssn.users = db.GetUsers(&tssn.uids)
	for i, user := range tssn.users {
		if user == nil {
			log.Fatalln("startTssn:", tssn.uids[i], "not in DB")
		}
	}
	tssn.answerTimer = time.NewTimer(answerTimeOut)
	if !tssn.answerTimer.Stop() {
		select {
		case <-tssn.answerTimer.C:
		default:
		}
	}
	for i := 0; i < 4; i++ {
		tssn.onlines[i] = true // regard as good by default
	}
	tssn.genIds()

	return tssn
}

func LoopTssn(bt model.BookType, uids [4]model.Uid) {
	log.Println("TSSN ++++", bt, uids)
	defer log.Println("TSSN ----", bt, uids)
	tssn := startTssn(bt, uids)
	defer close(tssn.done)
	Reg(tssn)
	defer Unreg(tssn)
	defer func() {
		if tssn.table != nil {
			saki.DeleteTableSession(tssn.table)
		}
	}()

	tssn.notifyLoad()

	hardTimer := time.NewTimer(2 * time.Hour)
	defer hardTimer.Stop()

	for tssn.anyOnline() {
		select {
		case msg := <-tssn.choose:
			tssn.handleChoose(msg.Uid, msg.Gidx)
		case uid := <-tssn.ready:
			tssn.handleReady(uid)
		case mta := <-tssn.action:
			tssn.handleAction(mta.Uid, mta.Act)
		case <-tssn.answerTimer.C:
			tssn.handleAnswerTimeout()
		case <-hardTimer.C:
			log.Println("TSSN hard timer")
			return
		}

		if tssn.table != nil && tssn.table.GameOver() {
			return
		}
	}
}

func (tssn *ss) Choose(msg *tbus.MsgChoose) {
	select {
	case tssn.choose <- msg:
	case <-tssn.done:
	}
}

func (tssn *ss) Ready(uid model.Uid) {
	select {
	case tssn.ready <- uid:
	case <-tssn.done:
	}
}

type msgTssnAction struct {
	uid model.Uid
	act *model.CsAction
}

func (tssn *ss) Action(uid model.Uid, act *model.CsAction) {
	select {
	case tssn.action <- &tbus.MsgAction{uid, act}:
	case <-tssn.done:
	}
}

func (tssn *ss) handleChoose(uid model.Uid, gidx int) {
	if i, ok := tssn.findUser(uid); ok {
		if tssn.state != tssnWaitChoose {
			log.Println("tssn.handleChoose wrong state", uid)
			tssn.kick(i)
			return
		}

		tssn.waits[i] = false
		cpu := len(tssn.gidcs) / 4 // choice per user
		tssn.gids[i] = tssn.gidcs[i*cpu+gidx]
		if !tssn.hasWait() {
			tssn.notifyChosen()
		}
	} else {
		log.Fatalln("loopTssn uid not found")
	}
}

func (tssn *ss) handleReady(uid model.Uid) {
	if i, ok := tssn.findUser(uid); ok {
		if tssn.state != tssnWaitReady {
			log.Println("tssn.handleReady wrong state", uid)
			tssn.kick(i)
			return
		}

		tssn.waits[i] = false
		if !tssn.hasWait() {
			tssn.start()
		}
	} else {
		log.Fatalln("loopTssn uid not found")
	}
}

func (tssn *ss) handleAction(uid model.Uid, act *model.CsAction) {
	i, _ := tssn.findUser(uid)
	if tssn.state != tssnWaitAction {
		log.Println("tssn.handleAction wrong state", uid,
			act.ActStr, act.ActArg, tssn.state)
		tssn.kick(i)
		return
	}

	if act.ActStr == "RESUME" {
		tssn.onlines[i] = true
	} else if act.Nonce != tssn.nonces[i] {
		log.Println(uid, "nonce", act.Nonce, "want", tssn.nonces[i])
		return
	}
	tssn.waits[i] = false
	mails := tssn.table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *ss) handleAnswerTimeout() {
	prevWaits := tssn.waits // copy, prevent overwrites incoming

	switch tssn.state {
	case tssnWaitChoose:
		for w := 0; w < 4; w++ {
			if prevWaits[w] {
				tssn.handleChoose(tssn.uids[w], 0)
				tssn.kick(w)
			}
		}
	case tssnWaitReady:
		for w := 0; w < 4; w++ {
			if prevWaits[w] {
				tssn.handleReady(tssn.uids[w])
				tssn.kick(w)
			}
		}
	case tssnWaitAction:
		tssn.sweepAll()
	}
}

func (tssn *ss) notifyLoad() {
	users := tssn.users

	msg := struct {
		Type       string
		Users      [4]*model.User
		TempDealer int
		Choices    [len(tssn.gidcs)]model.Gid
	}{"start", users, 0, tssn.gidcs}

	for i, _ := range tssn.waits {
		tssn.waits[i] = true
	}

	for i, uid := range tssn.uids {
		msg.TempDealer = (4 - i) % 4
		err := tssn.sendPeer(i, msg)
		if err != nil {
			tssn.handleChoose(uid, 0)
		}

		// rotate perspectives
		u0 := msg.Users[0]
		msg.Users[0] = msg.Users[1]
		msg.Users[1] = msg.Users[2]
		msg.Users[2] = msg.Users[3]
		msg.Users[3] = u0

		cs := &msg.Choices
		cpu := len(cs) / 4 // choice per user
		for i := 0; i < cpu; i++ {
			tmp := cs[i]
			for w := 0; w < 3; w++ {
				cs[w*cpu+i] = cs[(w+1)*cpu+i]
			}
			cs[3*cpu+i] = tmp
		}
	}

	tssn.resetAnswerTimer()
}

func (tssn *ss) notifyChosen() {
	tssn.state = tssnWaitReady

	msg := struct {
		Type    string
		GirlIds [4]model.Gid
	}{"chosen", tssn.gids}

	for w := 0; w < 4; w++ {
		tssn.waits[w] = true
	}

	for w := 0; w < 4; w++ {
		tssn.sendPeer(w, msg)

		gs := &msg.GirlIds
		gs[0], gs[1], gs[2], gs[3] = gs[1], gs[2], gs[3], gs[0]
	}

	tssn.resetAnswerTimer()
}

func (tssn *ss) hasWait() bool {
	r := &tssn.waits
	return r[0] || r[1] || r[2] || r[3]
}

func (tssn *ss) anyOnline() bool {
	o := &tssn.onlines
	return o[0] || o[1] || o[2] || o[3]
}

func (tssn *ss) findUser(uid model.Uid) (int, bool) {
	for i, u := range tssn.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (tssn *ss) start() {
	log.Println("TSSN ****", tssn.uids[0], tssn.gids)
	tssn.state = tssnWaitAction
	tssn.table = saki.NewTableSession(
		int(tssn.gids[0]), int(tssn.gids[1]),
		int(tssn.gids[2]), int(tssn.gids[3]))

	mails := tssn.table.Start()
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *ss) sweepOne(i int) {
	mails := tssn.table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *ss) sweepAll() {
	var targets int
	mails := tssn.table.SweepAll(&targets)
	for w := uint(0); w < 4; w++ {
		if (targets & (1 << w)) != 0 {
			tssn.waits[w] = false
			tssn.kick(int(w))
		}
	}
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *ss) kick(w int) {
	tssn.onlines[w] = false
	ubus.Logout(tssn.uids[w])
}

func (tssn *ss) handleMails(mails saki.MailVector) {
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
		tssn.resetAnswerTimer()
	}

	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		str := mails.Get(i).GetMsg()
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(str), &msg); err != nil {
			log.Fatalln("unmarshal c++ str", err)
		}
		if toWhom == -1 {
			tssn.handleSystemMail(msg, str)
		} else {
			tssn.sendUserMail(toWhom, msg)
		}
	}
}

func (tssn *ss) sendUserMail(who int, msg map[string]interface{}) {
	msg["Nonce"] = tssn.nonces[who]
	if msg["Event"] == "resume" {
		right := (who + 1) % 4
		cross := (who + 2) % 4
		left := (who + 3) % 4
		rotated := [4]*model.User{
			tssn.users[who],
			tssn.users[right],
			tssn.users[cross],
			tssn.users[left],
		}
		msg["Args"].(map[string]interface{})["users"] = rotated
	}

	err := tssn.sendPeer(who, msg)
	if err != nil && msg["Event"] == "activated" {
		if tssn.anyOnline() && !tssn.table.GameOver() {
			tssn.sweepOne(who)
		}
	}
}

func (tssn *ss) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := ubus.Peer(tssn.uids[i], msg)
		if err != nil {
			tssn.kick(i)
		}
		return err
	}
	return errors.New("peer not in session")
}

func (tssn *ss) resetAnswerTimer() {
	if !tssn.answerTimer.Stop() {
		select {
		case <-tssn.answerTimer.C:
		default: // prevent blocked by double-draining
		}
	}
	tssn.answerTimer.Reset(answerTimeOut)
}

func (tssn *ss) handleSystemMail(msg map[string]interface{},
	msgStr string) {
	switch msg["Type"] {
	case "round-start-log":
		fmt := "%v %v\n" +
			"\tr=%v e=%v d=%v al=%v depo=%v seed=%v"
		log.Printf(fmt,
			tssn.uids, tssn.gids,
			msg["round"], msg["extra"], msg["dealer"],
			msg["allLast"], msg["deposit"],
			uint(msg["seed"].(float64)))
	case "table-end-stat":
		var stat model.EndTableStat
		err := json.Unmarshal([]byte(msgStr), &stat)
		if err != nil {
			log.Fatalln("table-end-stat unmarshal", err)
		}
		tssn.injectUsers(stat.Replay)
		db.UpdateUserGirl(tssn.bookType, tssn.uids, tssn.gids, &stat)
		for w := 0; w < 4; w++ {
			ubus.UpdateInfo(tssn.uids[w])
		}
	case "riichi-auto":
		time.Sleep(1000 * time.Millisecond)
		who := int(msg["Who"].(float64))
		act := &model.CsAction{
			Nonce:  tssn.nonces[who],
			ActStr: "SPIN_OUT",
			ActArg: "-1",
		}
		tssn.handleAction(tssn.uids[who], act)
	default:
		log.Fatalln("unknown system mail", msg)
	}
}

func (tssn *ss) genIds() {
	avails := db.GetRankedGids()
	cpu := len(tssn.gidcs) / 4 // choice per user

	switch tssn.bookType.Index() {
	case 0:
		last14 := avails[len(avails)-14:]
		perm := rand.Perm(len(last14))
		for i := 0; i < len(tssn.gidcs); i++ {
			tssn.gidcs[i] = last14[perm[i]]
		}
	default: // not so many girls yet
		{
			top8 := avails[0:8]
			perm := rand.Perm(len(top8))
			for i := 0; i < 4; i++ {
				tssn.gidcs[(i%4)*cpu+i/4] = top8[perm[i]]
			}
		}
		{
			rest := avails[8:]
			perm := rand.Perm(len(rest))
			for i := 0; i < 8; i++ { // assume len(gidcs)==12
				tssn.gidcs[(i%4)*cpu+i/4+1] = rest[perm[i]]
			}
		}
	}

	for w := 0; w < 4; w++ {
		tssn.gids[w] = tssn.gidcs[w*cpu] // choose first as default
	}
}

func (tssn *ss) injectUsers(replay map[string]interface{}) {
	var users [4]map[string]interface{}
	for w := 0; w < 4; w++ {
		user := make(map[string]interface{})
		user["Id"] = tssn.uids[w]
		user["Username"] = tssn.users[w].Username
		user["Level"] = tssn.users[w].Level
		user["Rating"] = tssn.users[w].Rating
		users[w] = user
	}
	replay["users"] = users
}
