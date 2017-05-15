package tssn

import (
	"log"
	"math/rand"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
)

func (tssn *tssn) Choose(ctx actor.Context) {
	makeOnNext := func(ctx actor.Context) func() {
		return func() {
			ctx.SetBehavior(tssn.Ready)
			tssn.notifyChosen()
		}
	}

	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		tssn.handleChooseTimeout()
	case *pcChoose:
		tssn.handleChoose(msg.Uid, msg.Gidx, makeOnNext(ctx))
	case *pcReady:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Choose get pcReady")
	case *pcAction:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Choose get pcAction")
	case *ccChoose:
		tssn.handleChoose(msg.Uid, msg.Gidx, makeOnNext(ctx))
	default:
		log.Fatalf("tssn.Choose unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) genIds() {
	avails := db.GetRankedGids()
	cpu := len(tssn.gidcs) / 4 // choice per user

	switch tssn.bookType.Index() {
	case 0:
		last14 := avails[len(avails)-14:]
		perm := rand.Perm(len(last14))
		for i := 0; i < len(tssn.gidcs); i++ {
			tssn.gidcs[i] = last14[perm[i]]
		}
	default: // not so many girls yet
		{
			top8 := avails[0:8]
			perm := rand.Perm(len(top8))
			for i := 0; i < 4; i++ {
				tssn.gidcs[(i%4)*cpu+i/4] = top8[perm[i]]
			}
		}
		{
			rest := avails[8:]
			perm := rand.Perm(len(rest))
			for i := 0; i < 8; i++ { // assume len(gidcs)==12
				tssn.gidcs[(i%4)*cpu+i/4+1] = rest[perm[i]]
			}
		}
	}

	for w := 0; w < 4; w++ {
		tssn.gids[w] = tssn.gidcs[w*cpu] // choose first as default
	}
}

func (tssn *tssn) notifyLoad() {
	users := tssn.users

	msg := struct {
		Type       string
		Users      [4]*model.User
		TempDealer int
		Choices    [len(tssn.gidcs)]model.Gid
	}{"start", users, 0, tssn.gidcs}

	for i, _ := range tssn.waits {
		tssn.waits[i] = true
	}

	for i, uid := range tssn.uids {
		msg.TempDealer = (4 - i) % 4
		err := tssn.sendPeer(i, msg)
		if err != nil {
			tssn.p.Tell(&ccChoose{Uid: uid, Gidx: 0})
		}

		// rotate perspectives
		u0 := msg.Users[0]
		msg.Users[0] = msg.Users[1]
		msg.Users[1] = msg.Users[2]
		msg.Users[2] = msg.Users[3]
		msg.Users[3] = u0

		cs := &msg.Choices
		cpu := len(cs) / 4 // choice per user
		for i := 0; i < cpu; i++ {
			tmp := cs[i]
			for w := 0; w < 3; w++ {
				cs[w*cpu+i] = cs[(w+1)*cpu+i]
			}
			cs[3*cpu+i] = tmp
		}
	}
}

func (tssn *tssn) handleChoose(uid model.Uid, gidx int, onNext func()) {
	if i, ok := tssn.findUser(uid); ok {
		tssn.waits[i] = false
		cpu := len(tssn.gidcs) / 4 // choice per user
		tssn.gids[i] = tssn.gidcs[i*cpu+gidx]
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
			tssn.p.Tell(&ccChoose{Uid: tssn.uids[w], Gidx: 0})
			tssn.kick(w, "choose timeout")
		}
	}
}
