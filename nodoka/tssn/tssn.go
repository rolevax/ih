package tssn

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/ih/ako/model"
	"github.com/mjpancake/ih/nodoka"
	"github.com/mjpancake/ih/saki"
)

const recvTimeout = 15 * time.Second

func init() {
	// for girl choice list
	rand.Seed(time.Now().UnixNano())
}

type tssn struct {
	p           *actor.PID
	room        *model.Room
	waits       [4]bool
	onlines     [4]bool
	nonces      [4]int
	answerTimer *time.Timer
	table       saki.TableSession
	waitClient  bool
}

func Start(room *model.Room) {
	if !room.Four() {
		log.Fatal("tssn.Start room not 4")
	}

	tssn := &tssn{
		room:    room,
		onlines: [4]bool{true, true, true, true},
	}

	props := actor.FromInstance(tssn)
	pid, err := actor.SpawnPrefix(props, "tssn")
	if err != nil {
		log.Fatalln(err)
	}
	tssn.p = pid
	nodoka.Tmgr.Tell(&cpReg{add: true, tssn: tssn})
}

func (tssn *tssn) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.SetReceiveTimeout(recvTimeout)
		ctx.SetBehavior(tssn.Seat)
		uids := []model.Uid{}
		for _, u := range tssn.room.Users {
			uids = append(uids, u.Id)
		}
		log.Println("TSSN ++++", uids)
		tssn.notifySeat()
	default:
		log.Fatalf("tssn.Recv unexpected %T\n", ctx.Message())
	}
}

func (tssn *tssn) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := (&nodoka.MuSc{To: tssn.room.Users[i].Id, Msg: msg}).Req()
		if err != nil {
			tssn.kick(i, "write err")
		}
		return err
	}
	return fmt.Errorf("tssn.sendPeer: %d not online", tssn.room.Users[i].Id)
}

func (tssn *tssn) kick(uidx int, reason string) {
	tssn.onlines[uidx] = false
	nodoka.Umgr.Tell(&nodoka.MuKick{tssn.room.Users[uidx].Id, reason})
}

func (tssn *tssn) findUser(uid model.Uid) (int, bool) {
	for i, u := range tssn.room.Users {
		if u.Id == uid {
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
	nodoka.Tmgr.Tell(&cpReg{add: false, tssn: tssn})
	if tssn.table != nil {
		saki.DeleteTableSession(tssn.table)
	}
	ctx.SetBehavior(func(ctx actor.Context) {}) // clear bahavior

	uids := []model.Uid{}
	for _, u := range tssn.room.Users {
		uids = append(uids, u.Id)
	}
	log.Println("TSSN ----", uids)
}
