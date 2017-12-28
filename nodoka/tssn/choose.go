package tssn

import (
	"log"
	"math/rand"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
)

var availIds = []model.Gid{
	710113, 710114, 710115,
	712411, 712412, 712413,
	712611, 712613,
	712714, 712715,
	712915,
	713311, 713314,
	713301,
	713811, 713815,
	714915,
	715212,
	990001, 990002, 990003, 990011,
	990024,
}

type choices struct {
	gidcs [4][3]model.Gid
}

func newChoices(ruleId model.RuleId) *choices {
	c := &choices{}
	switch ruleId {
	case model.RuleFourDoges:
		// all zero
	case model.RuleClassic1In2:
		perms := rand.Perm(len(availIds))[:8]
		for i, n := range perms {
			who := i / 2
			what := i%2 + 1 // leave first 0
			c.gidcs[who][what] = availIds[n]
		}
	}
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
