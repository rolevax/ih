package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
)

type choices struct {
	gidcs [4][3]model.Gid
}

func newChoices() *choices {
	c := &choices{}
	// FUCK
	return c
}

func (tssn *tssn) Choose(ctx actor.Context) {
	makeOnNext := func(ctx actor.Context) func() {
		return func() {
			ctx.SetBehavior(tssn.Seat)
			tssn.notifySeat()
		}
	}

	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		tssn.handleChooseTimeout()
	case *pcChoose:
		tssn.handleChoose(msg.Uid, msg.Gidx, makeOnNext(ctx))
	case *pcSeat:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Choose get pcSeat")
	case *pcAction:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Choose get pcAction")
	case *ccChoose:
		tssn.handleChoose(msg.Uid, 0, makeOnNext(ctx))
	default:
		log.Fatalf("tssn.Choose unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) notifyChoose() {
	msg := &sc.TableInit{
		MatchResult: *tssn.match,
	}

	for i, _ := range tssn.waits {
		// bots have no choose process, use default 0 gidx
		tssn.waits[i] = tssn.match.Users[i].Id.IsHuman()
	}

	for i, uid := range tssn.match.Uids() {
		msg.Choices = tssn.choices.gidcs[i]
		if tssn.waits[i] {
			err := tssn.sendPeer(i, msg)
			if err != nil {
				tssn.p.Tell(&ccChoose{Uid: uid})
			}
		}
		msg = msg.RightPers()
	}
}

func (tssn *tssn) handleChoose(uid model.Uid, gidx int, onNext func()) {
	if i, ok := tssn.findUser(uid); ok {
		tssn.waits[i] = false
		tssn.gids[i] = tssn.choices.gidcs[i][gidx]
		if !tssn.hasWait() {
			onNext()
		}
	} else {
		log.Fatalf("tssn.handleChoose uid %d not found\n", uid)
	}
}

func (tssn *tssn) handleChooseTimeout() {
	for w := 0; w < 4; w++ {
		if tssn.waits[w] {
			tssn.p.Tell(&ccChoose{Uid: tssn.match.Users[w].Id})
			tssn.kick(w, "choose timeout")
		}
	}
}
