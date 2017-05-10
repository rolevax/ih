package ussn

import (
	"errors"
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/node"
)

var rec map[model.Uid]*ussn = make(map[model.Uid]*ussn)

func Init() {
	props := actor.FromFunc(Receive)
	node.Umgr = actor.Spawn(props)
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *node.MuSc:
		handleSc(msg.To, msg.Msg, ctx.Sender())
		ctx.Respond(nil)
	case *node.MuKick:
		handleKick(msg.Uid, msg.Reason)
	case *node.MuUpdateInfo:
		handleUpdateInfo(msg.Uid)
	case *cpReg:
		handleReg(msg.add, msg.ussn)
	case *cpCtUser:
		ctx.Respond(len(rec))
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
		}
		rec[ussn.user.Id] = ussn
	} else {
		node.Bmgr.Tell(&node.MbUnbook{Uid: ussn.user.Id})
		if prev, ok := rec[ussn.user.Id]; ok && prev == ussn {
			delete(rec, ussn.user.Id)
		}
	}
}

func handleSc(to model.Uid, msg interface{}, sender *actor.PID) {
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
	if ussn, ok := rec[uid]; ok {
		ussn.p.Tell(&pcUpdateInfo{})
	}
}

func handleKick(uid model.Uid, reason string) {
	if ussn, ok := rec[uid]; ok {
		ussn.p.Tell(fmt.Errorf("kick as %v", reason))
	}
}
