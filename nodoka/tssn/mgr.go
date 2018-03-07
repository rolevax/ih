package tssn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/nodoka"
)

var (
	rec     map[model.Uid]*tssn = make(map[model.Uid]*tssn)
	tableCt int
)

func Init() {
	props := actor.FromFunc(Receive)
	pid, err := actor.SpawnNamed(props, "Tmgr")
	if err != nil {
		log.Fatalln(err)
	}
	nodoka.Tmgr = pid
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *nodoka.MtHasUser:
		handleHasUser(msg.Uid, ctx.Respond)
	case *nodoka.MtCtTables:
		handleCtTables(ctx.Respond)
	case *nodoka.MtChoose:
		handleChoose(msg)
	case *nodoka.MtSeat:
		handleSeat(msg)
	case *nodoka.MtAction:
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

func handleCtTables(resp func(interface{})) {
	resp(tableCt)
}

func handleReg(add bool, tssn *tssn) {
	if add {
		for w := 0; w < 4; w++ {
			// TODO if exists, reject in some way
			rec[tssn.match.Users[w].Id] = tssn
		}
		tableCt++
	} else {
		for w := 0; w < 4; w++ {
			delete(rec, tssn.match.Users[w].Id)
		}
		tableCt--
	}
}

func handleChoose(msg *nodoka.MtChoose) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcChoose{msg})
	}
}

func handleSeat(msg *nodoka.MtSeat) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcSeat{msg})
	}
}

func handleAction(msg *nodoka.MtAction) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.p.Tell(&pcAction{msg})
	}
}
