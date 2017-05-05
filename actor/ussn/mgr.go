package ussn

import (
	"errors"

	"github.com/mjpancake/hisa/actor/ussn/ubus"
	"github.com/mjpancake/hisa/model"
)

var (
	rec   map[model.Uid]*ss = make(map[model.Uid]*ss)
	reg   chan *ss          = make(chan *ss)
	unreg chan *ss          = make(chan *ss)
)

func Loop() {
	for {
		select {
		case ussn := <-reg:
			handleReg(ussn)
		case ussn := <-unreg:
			handleUnreg(ussn)
		case mump := <-ubus.ChPeer:
			handlePeer(mump)
		case uid := <-ubus.ChUpdate:
			handleUpdate(uid)
		case uid := <-ubus.ChLogout:
			handleLogout(uid)
		case ch := <-ubus.ChCtUser:
			ch <- len(rec)
		}
	}
}

func Reg(ussn *ss) {
	reg <- ussn
}

func Unreg(ussn *ss) {
	unreg <- ussn
}

func handleReg(ussn *ss) {
	if prev, ok := rec[ussn.user.Id]; ok {
		prev.Logout(errors.New("kick login"))
	}
	rec[ussn.user.Id] = ussn
}

func handleUnreg(ussn *ss) {
	if prev, ok := rec[ussn.user.Id]; ok && prev == ussn {
		delete(rec, ussn.user.Id)
	}
}

func handlePeer(mp *ubus.MsgPeer) {
	if ussn, ok := rec[mp.To]; ok {
		go func() { // no block, thus go
			mp.ChErr <- ussn.Write(mp.Msg)
		}()
	} else {
		mp.ChErr <- errors.New("ussn not in rec")
	}
}

func handleUpdate(uid model.Uid) {
	if ussn, ok := rec[uid]; ok {
		// no block, thus go
		go ussn.UpdateInfo()
	}
}

func handleLogout(uid model.Uid) {
	if ussn, ok := rec[uid]; ok {
		// no block, thus go
		go ussn.Logout(errors.New("server kick"))
	}
}
