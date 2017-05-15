package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/node"
)

var (
	rec    map[model.Uid]*tssn = make(map[model.Uid]*tssn)
	btStat [model.BookTypeKinds]int
)

func Init() {
	props := actor.FromFunc(Receive)
	pid, err := actor.SpawnNamed(props, "Tmgr")
	if err != nil {
		log.Fatalln(err)
	}
	node.Tmgr = pid
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *node.MtHasUser:
		handleHasUser(msg.Uid, ctx.Respond)
	case *node.MtCtPlays:
		handleCtPlays(ctx.Respond)
	case *node.MtChoose:
		handleChoose(msg)
	case *node.MtReady:
		handleReady(msg)
	case *node.MtAction:
		handleAction(msg)
	case *cpReg:
		handleReg(msg.add, msg.tssn)
	default:
		log.Fatalf("Tmgr.Recv: unexpected %T\n", msg)
	}
}

func handleHasUser(uid model.Uid, resp func(interface{})) {
	_, ok := rec[uid]
	resp(ok)
}

func handleCtPlays(resp func(interface{})) {
	resp(btStat) // pass by value
}

func handleReg(add bool, tssn *tssn) {
	if add {
		for w := 0; w < 4; w++ {
			rec[tssn.uids[w]] = tssn
		}
		btStat[tssn.bookType.Index()]++
	} else {
		for w := 0; w < 4; w++ {
			delete(rec, tssn.uids[w])
		}
		btStat[tssn.bookType.Index()]--
	}
}

func handleReady(msg *node.MtReady) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcReady{msg})
	}
}

func handleChoose(msg *node.MtChoose) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcChoose{msg})
	}
}

func handleAction(msg *node.MtAction) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcAction{msg})
	}
}
