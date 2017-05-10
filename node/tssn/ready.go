package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/model"
)

func (tssn *tssn) Ready(ctx actor.Context) {
	makeOnNext := func(ctx actor.Context) func() {
		return func() {
			ctx.SetBehavior(tssn.Happy)
			tssn.start()
		}
	}

	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		tssn.handleReadyTimeout()
	case *pcReady:
		tssn.handleReady(msg.Uid, makeOnNext(ctx))
	case *ccReady:
		tssn.handleReady(msg.Uid, makeOnNext(ctx))
	default:
		log.Printf("tssn.Ready unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) notifyChosen() {
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
}

func (tssn *tssn) handleReady(uid model.Uid, onNext func()) {
	if i, ok := tssn.findUser(uid); ok {
		tssn.waits[i] = false
		if !tssn.hasWait() {
			onNext()
		}
	} else {
		log.Fatalf("tssn.handleReady uid %d not found\n", uid)
	}
}

func (tssn *tssn) handleReadyTimeout() {
	for w := 0; w < 4; w++ {
		if tssn.waits[w] {
			tssn.p.Tell(&ccReady{Uid: tssn.uids[w]})
			tssn.kick(w, "ready timeout")
		}
	}
}
