package book

import (
	"github.com/mjpancake/hisa/actor/tssn"
	"github.com/mjpancake/hisa/actor/tssn/tbus"
	"github.com/mjpancake/hisa/model"
)

var (
	book    chan *msgBookMgrBook
	unbook  chan model.Uid
	ctBooks chan chan [4]BookState
	states  [4]BookState
)

func init() {
	book = make(chan *msgBookMgrBook)
	unbook = make(chan model.Uid)
	ctBooks = make(chan chan [4]BookState)
}

func Loop() {
	for {
		select {
		case msg := <-book:
			handleBook(msg.uid, msg.bookType)
		case uid := <-unbook:
			handleUnbook(uid)
		case ch := <-ctBooks:
			ch <- states
		}
	}
}

type msgBookMgrBook struct {
	uid      model.Uid
	bookType model.BookType
}

func Book(uid model.Uid, bookType model.BookType) {
	msg := msgBookMgrBook{uid, bookType}
	book <- &msg
}

func Unbook(uid model.Uid) {
	unbook <- uid
}

func CtBooks() [4]BookState {
	ch := make(chan [4]BookState)
	ctBooks <- ch
	return <-ch
}

func handleBook(uid model.Uid, bookType model.BookType) {
	state := &states[bookType.Index()]

	for i := 0; i < state.Wait; i++ {
		if state.Waits[i] == uid {
			return
		}
	}

	if tbus.HasUser(uid) {
		return
	}

	state.Waits[state.Wait] = uid
	state.Wait++
	if state.Wait == 4 {
		handleStart(bookType)
	}
}

func handleUnbook(uid model.Uid) {
	for i := range states {
		states[i].removeIfAny(uid)
	}
}

func handleStart(bt model.BookType) {
	state := &states[bt.Index()]
	for _, uid := range states[bt.Index()].Waits {
		handleUnbook(uid)
	}
	go tssn.LoopTssn(bt, state.Waits)
}
