package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/ih/ako/model"
	"github.com/mjpancake/ih/ako/sc"
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
	case *pcChoose:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Ready get pcChoose")
	case *pcReady:
		tssn.handleReady(msg.Uid, makeOnNext(ctx))
	case *pcAction:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Ready get pcAction")
	case *ccReady:
		tssn.handleReady(msg.Uid, makeOnNext(ctx))
	default:
		log.Fatalf("tssn.Ready unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) notifyChosen() {
	msg := sc.NewChosen(tssn.gids)

	for w := 0; w < 4; w++ {
		// bots are always ready, wait your sister wait
		tssn.waits[w] = tssn.uids[w].IsHuman()
	}

	for w := 0; w < 4; w++ {
		if tssn.waits[w] {
			tssn.sendPeer(w, msg)
		}
		msg.RightPers()
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
