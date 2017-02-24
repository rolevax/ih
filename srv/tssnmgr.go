package srv

import (
)

type tssnMgr struct {
	rec			map[uid]*tssn
	btStat		[4]int
	reg			chan *tssn
	unreg		chan *tssn
	hasUser		chan *msgTssnMgrHasUser
	ctEachBt	chan chan [4]int
	choose		chan *msgTssnChoose
	ready		chan uid
	action		chan *msgTssnMgrAction
}

func newTssnMgr() *tssnMgr {
	tm := new(tssnMgr)

	tm.rec = make(map[uid]*tssn)
	tm.reg = make(chan *tssn)
	tm.unreg = make(chan *tssn)
	tm.hasUser = make(chan *msgTssnMgrHasUser)
	tm.ctEachBt = make(chan chan [4]int)
	tm.choose = make(chan *msgTssnChoose)
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
		case ch := <-tm.ctEachBt:
			ch <- tm.btStat
		case msg := <-tm.choose:
			tm.handleChoose(msg)
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

func (tm *tssnMgr) Choose(uid uid, gidx int) {
	msg := newMsgTssnChoose(uid, gidx)
	tm.choose <- msg
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

func (tm *tssnMgr) CtEachBt() [4]int {
	ch := make(chan [4]int)
	tm.ctEachBt <- ch
	return <-ch
}

func (tm *tssnMgr) handleReg(tssn *tssn) {
	for w := 0; w < 4; w++ {
		tm.rec[tssn.uids[w]] = tssn
	}
	tm.btStat[tssn.bookType.index()]++
}

func (tm *tssnMgr) handleUnreg(tssn *tssn) {
	for w := 0; w < 4; w++ {
		delete(tm.rec, tssn.uids[w])
	}
	tm.btStat[tssn.bookType.index()]--
}

func (tm *tssnMgr) handleReady(uid uid) {
	if tssn, ok := tm.rec[uid]; ok {
		tssn.Ready(uid)
	}
}

func (tm *tssnMgr) handleChoose(msg *msgTssnChoose) {
	if tssn, ok := tm.rec[msg.uid]; ok {
		tssn.Choose(msg)
	}
}

func (tm *tssnMgr) handleAction(msg *msgTssnMgrAction) {
	if tssn, ok := tm.rec[msg.uid]; ok {
		tssn.Action(msg.uid, msg.act)
	}
}


