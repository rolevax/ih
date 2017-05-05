package tbus

import "github.com/mjpancake/hisa/model"

var (
	ChHasUser  chan *MsgHasUser = make(chan *MsgHasUser)
	ChCtEachBt chan chan [4]int = make(chan chan [4]int)
	ChChoose   chan *MsgChoose  = make(chan *MsgChoose)
	ChReady    chan model.Uid   = make(chan model.Uid)
	ChAction   chan *MsgAction  = make(chan *MsgAction)
)

func HasUser(uid model.Uid) bool {
	msg := newMsgHasUser(uid)
	ChHasUser <- msg
	return <-msg.ChRes
}

func Choose(uid model.Uid, gidx int) {
	msg := &MsgChoose{uid, gidx}
	ChChoose <- msg
}

func Ready(uid model.Uid) {
	ChReady <- uid
}

func Action(uid model.Uid, act *model.CsAction) {
	msg := &MsgAction{uid, act}
	ChAction <- msg
}

func CtEachBt() [4]int {
	ch := make(chan [4]int)
	ChCtEachBt <- ch
	return <-ch
}
