package srv

import (
	"log"
	"time"
	"errors"
	"math/rand"
	"encoding/json"
	"github.com/mjpancake/mjpancake-server/saki"
)

const answerTimeOut = 15 * time.Second

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tssnState int

const (
	tssnWaitChoose	tssnState = iota
	tssnWaitReady
	tssnWaitAction
)

type tssn struct {
	bookType		bookType
	state			tssnState
	ready			chan uid
	choose			chan *msgTssnChoose
	action			chan *msgTssnAction
	done			chan struct{}
	uids			[4]uid
	gids			[4]gid
	gidcs			[8]gid
	waits			[4]bool
	onlines			[4]bool
	nonces			[4]int
	answerTimer		*time.Timer
	table			saki.TableSession
}

func startTssn(bt bookType, uids [4]uid) *tssn {
	tssn := new(tssn)

	tssn.bookType = bt
	tssn.state = tssnWaitChoose
	tssn.ready = make(chan uid)
	tssn.choose = make(chan *msgTssnChoose)
	tssn.action = make(chan *msgTssnAction)
	tssn.done = make(chan struct{})
	tssn.uids = uids
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

func loopTssn(bt bookType, uids [4]uid) {
	log.Println("TSSN ++++", bt, uids)
	defer log.Println("TSSN ----", bt, uids)
	tssn := startTssn(bt, uids)
	defer close(tssn.done)
	sing.TssnMgr.Reg(tssn)
	defer sing.TssnMgr.Unreg(tssn)
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
			tssn.handleChoose(msg.uid, msg.gidx)
		case uid := <-tssn.ready:
			tssn.handleReady(uid)
		case mta := <-tssn.action:
			tssn.handleAction(mta.uid, mta.act)
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

type msgTssnChoose struct {
	uid		uid
	gidx	int
}

func newMsgTssnChoose(uid uid, gidx int) *msgTssnChoose {
	msg := new(msgTssnChoose)
	msg.uid = uid
	msg.gidx = gidx
	return msg
}

func (tssn *tssn) Choose(msg *msgTssnChoose) {
	select {
	case tssn.choose <- msg:
	case <-tssn.done:
	}
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

func (tssn *tssn) handleChoose(uid uid, gidx int) {
	if i, ok := tssn.findUser(uid); ok {
		if tssn.state != tssnWaitChoose {
			log.Println("tssn.handleChoose wrong state", uid)
			tssn.kick(i)
			return
		}

		tssn.waits[i] = false
		tssn.gids[i] = tssn.gidcs[2 * i + gidx]
		if !tssn.hasWait() {
			tssn.notifyChosen()
		}
	} else {
		log.Fatalln("loopTssn uid not found")
	}
}

func (tssn *tssn) handleReady(uid uid) {
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

func (tssn *tssn) handleAction(uid uid, act *reqAction) {
	i, _ := tssn.findUser(uid)
	if tssn.state != tssnWaitAction {
		log.Println("tssn.handleAction wrong state", uid)
		tssn.kick(i)
		return
	}

	if act.ActStr == "RESUME" {
		tssn.onlines[i] = true
	} else if act.Nonce != tssn.nonces[i] {
		log.Println(uid, "nonce", act.Nonce, "want", tssn.nonces[i]);
		return
	}
	tssn.waits[i] = false
	mails := tssn.table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) handleAnswerTimeout() {
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

func (tssn *tssn) notifyLoad() {
	users := sing.Dao.GetUsers(&tssn.uids)
	for i, user := range users {
		if user == nil {
			log.Fatalln("tssn.nofityLoad:", tssn.uids[i], "not in DB")
		}
	}

	msg := struct {
		Type		string
		Users		[4]*user
		TempDealer	int
		Choices		[8]gid
	}{"start", users, 0, tssn.gidcs}

	for i, uid := range tssn.uids {
		msg.TempDealer = (4 - i) % 4
		tssn.waits[i] = true
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
		g0, g1 := cs[0], cs[1]
		cs[0], cs[1] = cs[2], cs[3]
		cs[2], cs[3] = cs[4], cs[5]
		cs[4], cs[5] = cs[6], cs[7]
		cs[6], cs[7] = g0, g1
	}

	tssn.resetAnswerTimer()
}

func (tssn *tssn) notifyChosen() {
	tssn.state = tssnWaitReady

	msg := struct {
		Type		string
		GirlIds		[4]gid
	}{"chosen", tssn.gids}

	for w := 0; w < 4; w++ {
		tssn.waits[w] = true
		tssn.sendPeer(w, msg)

		gs := &msg.GirlIds;
		gs[0], gs[1], gs[2], gs[3] = gs[1], gs[2], gs[3], gs[0]
	}

	tssn.resetAnswerTimer()
}

func (tssn *tssn) hasWait() bool {
	r := &tssn.waits
	return r[0] || r[1] || r[2] || r[3]
}

func (tssn *tssn) anyOnline() bool {
	o := &tssn.onlines
	return o[0] || o[1] || o[2] || o[3]
}

func (tssn *tssn) findUser(uid uid) (int, bool) {
	for i, u := range tssn.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (tssn *tssn) start() {
	log.Println("TSSN ****", tssn.uids[0], tssn.gids)
	tssn.state = tssnWaitAction
	tssn.table = saki.NewTableSession(
		int(tssn.gids[0]),int(tssn.gids[1]),
		int(tssn.gids[2]), int(tssn.gids[3]))

	mails := tssn.table.Start()
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) sweepOne(i int) {
	mails := tssn.table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) sweepAll() {
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

func (tssn *tssn) kick(w int) {
	tssn.onlines[w] = false
	sing.UssnMgr.Logout(tssn.uids[w])
}

func (tssn *tssn) handleMails(mails saki.MailVector) {
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

func (tssn *tssn) sendUserMail(who int, msg map[string]interface{}) {
	msg["Nonce"] = tssn.nonces[who]
	if msg["Event"] == "resume" {
		users := sing.Dao.GetUsers(&tssn.uids)
		for i, user := range users {
			if user == nil {
				log.Fatalln("tssn.send-resume:", tssn.uids[i], "not in DB")
			}
		}
		right := (who + 1) % 4
		cross := (who + 2) % 4
		left := (who + 3) % 4
		rotated := [4]*user{users[who], users[right],
							users[cross], users[left]}
		msg["Args"].(map[string]interface{})["users"] = rotated
	}

	err := tssn.sendPeer(who, msg)
	if err != nil && msg["Event"] == "activated" {
		if tssn.anyOnline() && !tssn.table.GameOver() {
			tssn.sweepOne(who)
		}
	}
}

func (tssn *tssn) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := sing.UssnMgr.Peer(tssn.uids[i], msg)
		if err != nil {
			tssn.kick(i)
		}
		return err
	}
	return errors.New("peer not in session")
}

func (tssn *tssn) resetAnswerTimer() {
	if !tssn.answerTimer.Stop() {
		select {
		case <-tssn.answerTimer.C:
		default: // prevent blocked by double-draining
		}
	}
	tssn.answerTimer.Reset(answerTimeOut)
}

type systemEndTableStat struct {
	Ranks		[4]int
	Points		[4]int
}

func (tssn *tssn) handleSystemMail(msg map[string]interface{},
	msgStr string) {
	switch (msg["Type"]) {
	case "round-start-log":
		fmt := "%v %v\n" +
			   "\tr=%v e=%v d=%v al=%v depo=%v seed=%v"
		log.Printf(fmt,
				   tssn.uids, tssn.gids,
				   msg["round"], msg["extra"], msg["dealer"],
				   msg["allLast"], msg["deposit"],
				   uint(msg["seed"].(float64)))
	case "table-end-stat":
		var stat systemEndTableStat
		err := json.Unmarshal([]byte(msgStr), &stat)
		if err != nil {
			log.Fatalln("table-end-stat unmarshal", err)
		}
		sing.Dao.UpdateUserGirl(tssn.bookType, tssn.uids, tssn.gids, &stat)
		for w := 0; w < 4; w++ {
			sing.UssnMgr.UpdateInfo(tssn.uids[w])
		}
	case "riichi-auto":
		time.Sleep(1000 * time.Millisecond)
		who := int(msg["Who"].(float64))
		act := reqAction{tssn.nonces[who], "SPIN_OUT", "-1"}
		tssn.handleAction(tssn.uids[who], &act)
	default:
		log.Fatalln("unknown system mail", msg)
	}
}

func (tssn *tssn) genIds() {
	avails := sing.Dao.GetRankedGids()

	switch tssn.bookType.index() {
	case 0:
		avails = avails[len(avails) - 14:]
	case 1:
		avails = avails[0:10]
	default: // nobody may enter here by now, though
		avails = avails[0:10]
	}

	perm := rand.Perm(len(avails))
	for i := 0; i < 8; i++ { // 2-choose-1, thus 8 in total
		tssn.gidcs[i] = avails[perm[i]]
	}
	for i := 0; i < 4; i++ {
		tssn.gids[i] = tssn.gidcs[2 * i] // default values
	}
}

