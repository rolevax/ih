package tssn

import (
	"github.com/mjpancake/hisa/actor/tssn/tbus"
	"github.com/mjpancake/hisa/model"
)

var (
	rec    map[model.Uid]*ss = make(map[model.Uid]*ss)
	btStat [4]int
	reg    chan *ss = make(chan *ss)
	unreg  chan *ss = make(chan *ss)
)

func Loop() {
	for {
		select {
		case tssn := <-reg:
			handleReg(tssn)
		case tssn := <-unreg:
			handleUnreg(tssn)
		case msg := <-tbus.ChHasUser:
			_, ok := rec[msg.Uid]
			msg.ChRes <- ok
		case ch := <-tbus.ChCtEachBt:
			ch <- btStat
		case msg := <-tbus.ChChoose:
			handleChoose(msg)
		case uid := <-tbus.ChReady:
			handleReady(uid)
		case msg := <-tbus.ChAction:
			handleAction(msg)
		}
	}
}

func Reg(tssn *ss) {
	reg <- tssn
}

func Unreg(tssn *ss) {
	unreg <- tssn
}

func handleReg(tssn *ss) {
	for w := 0; w < 4; w++ {
		rec[tssn.uids[w]] = tssn
	}
	btStat[tssn.bookType.Index()]++
}

func handleUnreg(tssn *ss) {
	for w := 0; w < 4; w++ {
		delete(rec, tssn.uids[w])
	}
	btStat[tssn.bookType.Index()]--
}

func handleReady(uid model.Uid) {
	if tssn, ok := rec[uid]; ok {
		tssn.Ready(uid)
	}
}

func handleChoose(msg *tbus.MsgChoose) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.Choose(msg)
	}
}

func handleAction(msg *tbus.MsgAction) {
	if tssn, ok := rec[msg.Uid]; ok {
		tssn.Action(msg.Uid, msg.Act)
	}
}
