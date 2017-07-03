package book

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/ih/ako/model"
	"github.com/mjpancake/ih/nodoka"
	"github.com/mjpancake/ih/nodoka/tssn"
)

var states [model.BookTypeKinds]BookState

func Init() {
	props := actor.FromFunc(Receive)
	pid, err := actor.SpawnNamed(props, "Bmgr")
	if err != nil {
		log.Fatalln(err)
	}
	nodoka.Bmgr = pid
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *nodoka.MbBook:
		handleBook(msg.Uid, msg.BookType)
	case *nodoka.MbUnbook:
		handleUnbook(msg.Uid)
	case *nodoka.MbCtBooks:
		handleCtBooks(ctx.Respond)
	default:
		log.Fatalf("Bmgr.Recv unexpected %T\n", msg)
	}
}

func handleBook(uid model.Uid, bookType model.BookType) {
	state := &states[bookType.Index()]

	for i := 0; i < state.Wait; i++ {
		if state.Waits[i] == uid {
			return
		}
	}

	playing, err := (&nodoka.MtHasUser{Uid: uid}).Req()
	if err != nil {
		log.Println("Bmgr:", err)
		return
	}
	if playing {
		return
	}

	state.Waits[state.Wait] = uid
	state.Wait++
	if state.Wait == bookType.NeedUser() {
		state.fillByAi()
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
	uids := state.Waits // copy
	for _, uid := range uids {
		handleUnbook(uid)
	}
	tssn.Start(bt, uids)
}

func handleCtBooks(resp func(interface{})) {
	var cts [model.BookTypeKinds]int
	for i, v := range states {
		cts[i] = v.Wait
	}
	resp(cts)
}
