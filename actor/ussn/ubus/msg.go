package ubus

import "github.com/mjpancake/hisa/model"

type MsgPeer struct {
	To    model.Uid
	Msg   interface{}
	ChErr chan error
}

func newMsgPeer(to model.Uid, msg interface{}) *MsgPeer {
	mp := &MsgPeer{
		To:    to,
		Msg:   msg,
		ChErr: make(chan error),
	}
	return mp
}
