package srv

import (
)

type bookMgr struct {
	book		chan *msgBookMgrBook
	unbook		chan uid
	ctBooks		chan chan [4]bookState
	states		[4]bookState
}

type bookState struct {
	waits		[4]uid
	wait		int
}

func (bs *bookState) removeIfAny(uid uid) {
	i := 0
	for i < bs.wait && bs.waits[i] != uid {
		i++
	}
	if i == bs.wait {
		return
	}
	// swap to back, and pop back
	e := bs.wait - 1;
	bs.waits[i], bs.waits[e] = bs.waits[e], bs.waits[i];
	bs.wait--
}

func newBookMgr() *bookMgr {
	bm := new(bookMgr)

	bm.book = make(chan *msgBookMgrBook)
	bm.unbook = make(chan uid)
	bm.ctBooks = make(chan chan [4]bookState)

	return bm;
}

func (bm *bookMgr) Loop() {
	for {
		select {
		case msg := <-bm.book:
			bm.handleBook(msg.uid, msg.bookType)
		case uid := <-bm.unbook:
			bm.handleUnbook(uid)
		case ch := <-bm.ctBooks:
			ch <- bm.states
		}
	}
}

type msgBookMgrBook struct {
	uid			uid
	bookType	bookType
}

func (bm *bookMgr) Book(uid uid, bookType bookType) {
	msg := msgBookMgrBook{uid, bookType}
	bm.book <- &msg
}

func (bm *bookMgr) Unbook(uid uid) {
	bm.unbook <- uid
}

func (bm *bookMgr) CtBooks() [4]bookState {
	ch := make(chan [4]bookState)
	bm.ctBooks <- ch
	return <-ch
}

func (bm *bookMgr) handleBook(uid uid, bookType bookType) {
	state := &bm.states[bookType.index()]

	for i := 0; i < state.wait; i++ {
		if state.waits[i] == uid {
			return
		}
	}
	if sing.TssnMgr.HasUser(uid) {
		return
	}

	state.waits[state.wait] = uid;
	state.wait++
	if state.wait == 4 {
		bm.handleStart(bookType)
	}
}

func (bm *bookMgr) handleUnbook(uid uid) {
	for i := range bm.states {
		bm.states[i].removeIfAny(uid)
	}
}

func (bm *bookMgr) handleStart(bt bookType) {
	state := &bm.states[bt.index()]
	for _, uid := range bm.states[bt.index()].waits {
		bm.handleUnbook(uid)
	}
	go loopTssn(bt, state.waits)
}



