package ussn

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/nodoka"
)

var (
	rec   map[model.Uid]*ussn = make(map[model.Uid]*ussn)
	water []string            // optimize it later
)

func Init() {
	props := actor.FromFunc(Receive)
	pid, err := actor.SpawnNamed(props, "Umgr")
	if err != nil {
		log.Fatalln(err)
	}
	nodoka.Umgr = pid
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *nodoka.MuSc:
		handleSc(msg.To, msg.Msg, ctx.Sender())
	case *nodoka.MuKick:
		handleKick(msg.Uid, msg.Reason)
	case *nodoka.MuUpdateInfo:
		handleUpdateInfo(msg.Uid)
	case *cpReg:
		handleReg(msg.add, msg.ussn)
	case *cpWater:
		w := &pcWater{ct: len(rec), water: make([]string, len(water))}
		copy(w.water, water)
		ctx.Respond(w)
	default:
		log.Fatalf("Umgr.Recv: unexpected %T\n", msg)
	}
}

func handleReg(add bool, ussn *ussn) {
	if ussn.user == nil {
		return // login failure, entry won't present
	}

	if add {
		if prev, ok := rec[ussn.user.Id]; ok {
			prev.p.Tell(errors.New("kick by force login"))
		} else {
			// log only on non-force
			addWater(ussn.user.Username, "上线")
		}
		rec[ussn.user.Id] = ussn
	} else {
		if prev, ok := rec[ussn.user.Id]; ok && prev == ussn {
			delete(rec, ussn.user.Id)
			addWater(ussn.user.Username, "下线")
		}
	}
}

func handleSc(to model.Uid, msg interface{}, sender *actor.PID) {
	if to.IsBot() {
		botSc(to, msg, sender)
		return
	}

	if ussn, ok := rec[to]; ok {
		if sender != nil {
			ussn.p.Request(&pcSc{msg: msg}, sender)
		} else {
			ussn.p.Tell(&pcSc{msg: msg})
		}
	} else {
		if sender != nil {
			sender.Tell(fmt.Errorf("ussn %d not online", to))
		}
	}
}

func handleUpdateInfo(uid model.Uid) {
	if uid.IsBot() {
		return
	}

	if ussn, ok := rec[uid]; ok {
		ussn.p.Tell(&pcUpdateInfo{})
	}
}

func handleKick(uid model.Uid, reason string) {
	if uid.IsBot() {
		log.Println("Umgr: kicing a bot", uid)
		return
	}

	if ussn, ok := rec[uid]; ok {
		ussn.p.Tell(fmt.Errorf("kick as %v", reason))
	}
}

func addWater(username, what string) {
	wat := fmt.Sprintf(
		"%s %s %s",
		time.Now().Format("15:04"),
		username,
		what,
	)
	water = append(water, wat)
	if len(water) > 12 {
		water = water[1:]
	}
}
