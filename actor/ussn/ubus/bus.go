package ubus

import "github.com/mjpancake/hisa/model"

var (
	ChPeer   chan *MsgPeer  = make(chan *MsgPeer)
	ChUpdate chan model.Uid = make(chan model.Uid)
	ChLogout chan model.Uid = make(chan model.Uid)
	ChCtUser chan chan int  = make(chan chan int)
)

func Peer(uid model.Uid, msg interface{}) error {
	mp := newMsgPeer(uid, msg)
	ChPeer <- mp
	return <-mp.ChErr
}

func UpdateInfo(uid model.Uid) {
	ChUpdate <- uid
}

func Logout(uid model.Uid) {
	ChLogout <- uid
}

func CtUser() int {
	ch := make(chan int)
	ChCtUser <- ch
	return <-ch
}
