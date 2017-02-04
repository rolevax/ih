package srv

import (
)

type tssnMgr struct {
	rec		map[uid]*tssn
	reg		chan *tssn
	unreg	chan *tssn
}

func newTssnMgr() *tssnMgr {
	tm := new(tssnMgr)

	tm.rec = make(map[uid]*tssn)
	tm.reg = make(chan *tssn)
	tm.unreg = make(chan *tssn)

	return tm
}

func (tm *tssnMgr) Loop() {
	for {
		select {
		case tssn := <-tm.reg:
			tm.handleReg(tssn)
		case tssn := <-tm.unreg:
			tm.handleUnreg(tssn)
		}
	}
}

func (tm *tssnMgr) Reg(tssn *tssn) {
	tm.reg <- tssn
}

func (tm *tssnMgr) Unreg(tssn *tssn) {
	tm.unreg <- tssn
}

func (tm *tssnMgr) HasUser(uid uid) bool {
	_, ok := tm.rec[uid]
	return ok
}

func (tm *tssnMgr) Ready(uid uid) {
	if tssn, ok := tm.rec[uid]; ok {
		tssn.Ready(uid)
	}
}

func (tm *tssnMgr) Action(uid uid, act *reqAction) {
	if tssn, ok := tm.rec[uid]; ok {
		tssn.Action(uid, act)
	}
}

func (tm *tssnMgr) CtUser() int {
	return len(tm.rec)
}

func (tm *tssnMgr) handleReg(tssn *tssn) {
	for w := 0; w < 4; w++ {
		tm.rec[tssn.uids[w]] = tssn
	}
}

func (tm *tssnMgr) handleUnreg(tssn *tssn) {
	for w := 0; w < 4; w++ {
		delete(tm.rec, tssn.uids[w])
	}
}


