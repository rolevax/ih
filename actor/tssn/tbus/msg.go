package tbus

import "github.com/mjpancake/hisa/model"

type MsgHasUser struct {
	Uid   model.Uid
	ChRes chan bool
}

func newMsgHasUser(uid model.Uid) *MsgHasUser {
	msg := &MsgHasUser{
		Uid:   uid,
		ChRes: make(chan bool),
	}
	return msg
}

type MsgAction struct {
	Uid model.Uid
	Act *model.CsAction
}

type MsgChoose struct {
	Uid  model.Uid
	Gidx int
}
