package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
)

func (tssn *tssn) Seat(ctx actor.Context) {
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
		tssn.handleSeatTimeout()
	case *pcSeat:
		tssn.handleSeat(msg.Uid, makeOnNext(ctx))
	case *pcAction:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Seat get pcAction")
	case *ccSeat:
		tssn.handleSeat(msg.Uid, makeOnNext(ctx))
	default:
		log.Fatalf("tssn.Seat unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) notifySeat() {
	// TODO random seat

	for w := 0; w < 4; w++ {
		// bots are always ready, wait your sister wait
		tssn.waits[w] = tssn.room.Users[w].Id.IsHuman()
	}

	msg := &sc.Seat{
		TempDealer: 0,
		Room:       *tssn.room,
	}
	for w := 0; w < 4; w++ {
		if tssn.waits[w] {
			tssn.sendPeer(w, msg)
		}
		msg = msg.RightPers()
	}
}

func (tssn *tssn) handleSeat(uid model.Uid, onNext func()) {
	if i, ok := tssn.findUser(uid); ok {
		tssn.waits[i] = false
		if !tssn.hasWait() {
			onNext()
		}
	} else {
		log.Fatalf("tssn.handleSeat uid %d not found\n", uid)
	}
}

func (tssn *tssn) handleSeatTimeout() {
	for w := 0; w < 4; w++ {
		if tssn.waits[w] {
			tssn.p.Tell(&ccSeat{Uid: tssn.room.Users[w].Id})
			tssn.kick(w, "seat timeout")
		}
	}
}
