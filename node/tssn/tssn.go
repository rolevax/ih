package tssn

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/node"
	"github.com/mjpancake/hisa/saki"
)

const recvTimeout = 15 * time.Second

func init() {
	// for girl choice list
	rand.Seed(time.Now().UnixNano())
}

type tssn struct {
	p           *actor.PID
	bookType    model.BookType
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

func Start(bt model.BookType, uids [4]model.Uid) {
	users := db.GetUsers(&uids)
	for i, user := range users {
		if user == nil {
			log.Fatalln("uid", uids[i], "not in DB")
		}
	}

	tssn := &tssn{
		bookType: bt,
		uids:     uids,
		users:    users,
		onlines:  [4]bool{true, true, true, true},
	}
	tssn.genIds()

	props := actor.FromInstance(tssn)
	pid, err := actor.SpawnPrefix(props, "tssn")
	if err != nil {
		log.Fatalln(err)
	}
	tssn.p = pid
	node.Tmgr.Tell(&cpReg{add: true, tssn: tssn})
}

func (tssn *tssn) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.SetReceiveTimeout(recvTimeout)
		ctx.SetBehavior(tssn.Choose)
		log.Println("TSSN ++++", tssn.bookType, tssn.uids)
		tssn.notifyLoad()
	default:
		log.Fatalf("tssn.Recv unexpected %T\n", ctx.Message())
	}
}

func (tssn *tssn) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := (&node.MuSc{To: tssn.uids[i], Msg: msg}).Req()
		if err != nil {
			tssn.kick(i, "write err")
		}
		return err
	}
	return fmt.Errorf("tssn.sendPeer: %d not online", tssn.uids[i])
}

func (tssn *tssn) kick(uidx int, reason string) {
	tssn.onlines[uidx] = false
	node.Umgr.Tell(&node.MuKick{tssn.uids[uidx], reason})
}

func (tssn *tssn) findUser(uid model.Uid) (int, bool) {
	for i, u := range tssn.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (tssn *tssn) hasWait() bool {
	r := &tssn.waits
	return r[0] || r[1] || r[2] || r[3]
}

func (tssn *tssn) anyOnline() bool {
	o := &tssn.onlines
	return o[0] || o[1] || o[2] || o[3]
}

func (tssn *tssn) checkGameOver() {
	if !tssn.anyOnline() || (tssn.table != nil && tssn.table.GameOver()) {
		tssn.p.Stop()
	}
}

func (tssn *tssn) bye(ctx actor.Context) {
	node.Tmgr.Tell(&cpReg{add: false, tssn: tssn})
	if tssn.table != nil {
		saki.DeleteTableSession(tssn.table)
	}
	ctx.SetBehavior(func(ctx actor.Context) {}) // clear bahavior
	log.Println("TSSN ----", tssn.bookType, tssn.uids)
}
