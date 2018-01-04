package tssn

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
)

var (
	availIds = []model.Gid{}

	girlCosts = map[model.Gid]int{
		0:      0,
		710113: 2000, 710114: 3900, 710115: 8000,
		712411: 5200, 712412: 5200, 712413: 3900,
		712611: 8000, 712613: 7700,
		712714: 2000, 712715: 3900,
		712915: 3900,
		713311: 7700, 713314: 3900,
		713301: 3900,
		713811: 3900, 713815: 7700,
		714915: 7700,
		715212: 8000,
		990001: 5200, 990002: 11600, 990003: 7700, 990011: 12000,
		990024: 5200,
	}
)

func init() {
	for gid, _ := range girlCosts {
		if gid == 0 {
			continue
		}

		availIds = append(availIds, gid)
	}
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
		cost := girlCosts[tssn.gids[i]]
		if tssn.match.Users[i].Food-cost < 0 {
			tssn.gids[i] = 0 // use doge
			tssn.kick(i, "insufficient food")
		} else {
			tssn.gids[i] = tssn.choices.gidcs[i][gidx]
		}
		tssn.addFoodChange(i, &model.FoodChange{
			Delta:  -cost,
			Reason: fmt.Sprintf("%v吃掉", tssn.gids[i]),
		})

		tssn.waits[i] = false
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
