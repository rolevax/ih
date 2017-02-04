package srv

import (
)

type bookMgr struct {
	book	chan uid
	unbook	chan uid
	waits	[4]uid
	wait	int
}

func newBookMgr() *bookMgr {
	bm := new(bookMgr)

	bm.book = make(chan uid)
	bm.unbook = make(chan uid)
	bm.wait = 0

	return bm;
}

func (bm *bookMgr) Loop() {
	for {
		select {
		case uid := <-bm.book:
			bm.handleBook(uid)
		case uid := <-bm.unbook:
			bm.handleUnbook(uid)
		}
	}
}

func (bm *bookMgr) Book(uid uid) {
	bm.book <- uid
}

func (bm *bookMgr) Unbook(uid uid) {
	bm.unbook <- uid
}

func (bm *bookMgr) CtBook() int {
	return bm.wait
}

func (bm *bookMgr) handleBook(uid uid) {
	for i := 0; i < bm.wait; i++ {
		if bm.waits[i] == uid {
			return
		}
	}
	if sing.TssnMgr.HasUser(uid) {
		return
	}

	bm.waits[bm.wait] = uid;
	bm.wait++
	if bm.wait == 4 {
		go loopTssn(bm.waits)
		bm.wait = 0
	}
}

func (bm *bookMgr) handleUnbook(uid uid) {
	i := 0
	for i < bm.wait && bm.waits[i] != uid {
		i++
	}

	if i == bm.wait {
		return
	}

	// swap to back, and pop back
	e := bm.wait - 1;
	bm.waits[i], bm.waits[e] = bm.waits[e], bm.waits[i];
	bm.wait--
}

