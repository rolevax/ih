package srv

import (
)

type tssnMgr struct {
	rec		map[uid]*tssn
	reg		chan *tssn
	unreg	chan *tssn
	hasUser	chan *msgTssnMgrHasUser
	ctUser	chan chan int
	ready	chan uid
	action	chan *msgTssnMgrAction
}

func newTssnMgr() *tssnMgr {
	tm := new(tssnMgr)

	tm.rec = make(map[uid]*tssn)
	tm.reg = make(chan *tssn)
	tm.unreg = make(chan *tssn)
	tm.hasUser = make(chan *msgTssnMgrHasUser)
	tm.ctUser = make(chan chan int)
	tm.ready = make(chan uid)
	tm.action = make(chan *msgTssnMgrAction)

	return tm
}

func (tm *tssnMgr) Loop() {
	for {
		select {
		case tssn := <-tm.reg:
			tm.handleReg(tssn)
		case tssn := <-tm.unreg:
			tm.handleUnreg(tssn)
		case mtmhu := <-tm.hasUser:
			_, ok := tm.rec[mtmhu.uid]
			mtmhu.chRes <- ok
		case ch := <-tm.ctUser:
			ch <- len(tm.rec)
		case uid := <-tm.ready:
			tm.handleReady(uid)
		case msg := <-tm.action:
			tm.handleAction(msg)
		}
	}
}

func (tm *tssnMgr) Reg(tssn *tssn) {
	tm.reg <- tssn
}

func (tm *tssnMgr) Unreg(tssn *tssn) {
	tm.unreg <- tssn
}

type msgTssnMgrHasUser struct {
	uid		uid
	chRes	chan bool
}

func newMsgTssnMgrHasUser(uid uid) *msgTssnMgrHasUser {
	msg := new(msgTssnMgrHasUser)
	msg.uid = uid
	msg.chRes = make(chan bool)
	return msg
}

func (tm *tssnMgr) HasUser(uid uid) bool {
	msg := newMsgTssnMgrHasUser(uid)
	tm.hasUser <- msg
	return <-msg.chRes
}

func (tm *tssnMgr) Ready(uid uid) {
	tm.ready <- uid
}

type msgTssnMgrAction struct {
	uid		uid
	act		*reqAction
}

func (tm *tssnMgr) Action(uid uid, act *reqAction) {
	msg := msgTssnMgrAction{uid, act}
	tm.action <- &msg
}

func (tm *tssnMgr) CtUser() int {
	ch := make(chan int)
	tm.ctUser <- ch
	return <-ch
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

func (tm *tssnMgr) handleReady(uid uid) {
	if tssn, ok := tm.rec[uid]; ok {
		tssn.Ready(uid)
	}
}

func (tm *tssnMgr) handleAction(msg *msgTssnMgrAction) {
	if tssn, ok := tm.rec[msg.uid]; ok {
		tssn.Action(msg.uid, msg.act)
	}
}


