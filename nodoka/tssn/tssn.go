package tssn

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/nodoka"
)

const recvTimeout = 15 * time.Second

func init() {
	// for girl choice list
	rand.Seed(time.Now().UnixNano())
}

type tssn struct {
	p            *actor.PID
	match        *model.MatchResult
	choices      *choices
	gids         [4]model.Gid
	waits        [4]bool
	onlines      [4]bool
	nonces       [4]int
	foodChangess [4][]*model.FoodChange
	answerTimer  *time.Timer
	waitClient   bool
	gameOver     bool
}

func Start(mr *model.MatchResult) {
	tssn := &tssn{
		match:   mr,
		choices: newChoices(mr.RuleId),
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
		ctx.SetBehavior(tssn.Choose)
		tssn.notifyChoose()
		log.Println("TSSN ++++", tssn.match.Uids())
	default:
		log.Fatalf("tssn.Recv unexpected %T\n", ctx.Message())
	}
}

func (tssn *tssn) sendPeer(i int, msg interface{}) error {
	if tssn.onlines[i] {
		err := (&nodoka.MuSc{To: tssn.match.Users[i].Id, Msg: msg}).Req()
		if err != nil {
			tssn.kick(i, "write err")
		}
		return err
	}
	return fmt.Errorf(
		"tssn.sendPeer: %d not online",
		tssn.match.Users[i].Id,
	)
}

func (tssn *tssn) kick(uidx int, reason string) {
	tssn.onlines[uidx] = false
	nodoka.Umgr.Tell(&nodoka.MuKick{tssn.match.Users[uidx].Id, reason})
}

func (tssn *tssn) findUser(uid model.Uid) (int, bool) {
	for i, u := range tssn.match.Users {
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
	if !tssn.anyOnline() || tssn.gameOver {
		tssn.p.Stop()
	}
}

func (tssn *tssn) bye(ctx actor.Context) {
	for w := 0; w < 4; w++ {
		sc := sc.TableEnd{
			FoodChanges: tssn.foodChangess[w],
		}
		_ = tssn.sendPeer(w, sc)

		sumDelta := 0
		for _, fc := range sc.FoodChanges {
			sumDelta += fc.Delta
		}
		_ = mako.UpdateFood(tssn.match.Users[w].Id, sumDelta)
	}

	nodoka.Tmgr.Tell(&cpReg{add: false, tssn: tssn})
	ctx.SetBehavior(func(ctx actor.Context) {}) // clear bahavior

	log.Println("TSSN ----", tssn.match.Uids())
}

func (tssn *tssn) addFoodChange(w int, fc *model.FoodChange) {
	tssn.foodChangess[w] = append(tssn.foodChangess[w], fc)
}
