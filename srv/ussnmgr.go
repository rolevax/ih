package srv

import (
	"errors"
)

type ussnMgr struct {
	rec		map[uid]*ussn
	reg		chan *ussn
	unreg	chan *ussn
    peer    chan *msgUssnMgrPeer
	update	chan *user
	logout	chan uid
	ctUser	chan chan int
}

func newUssnMgr() *ussnMgr {
	um := new(ussnMgr)

	um.rec = make(map[uid]*ussn)
	um.reg = make(chan *ussn)
	um.unreg = make(chan *ussn)
    um.peer = make(chan *msgUssnMgrPeer)
	um.update = make(chan *user)
	um.logout = make(chan uid)
	um.ctUser = make(chan chan int)

	return um
}

type msgUssnMgrPeer struct {
	to		uid
	msg		interface{}
	chErr	chan error
}

func newMsgUssnMgrPeer(to uid, msg interface{}) *msgUssnMgrPeer {
	mp := new(msgUssnMgrPeer)
	mp.to = to
	mp.msg = msg
	mp.chErr = make(chan error)
	return mp
}

func (um *ussnMgr) Loop() {
	for {
		select {
		case ussn := <-um.reg:
			um.handleReg(ussn)
		case ussn := <-um.unreg:
			um.handleUnreg(ussn)
		case mump := <-um.peer:
			um.handlePeer(mump)
		case user := <-um.update:
			um.handleUpdate(user)
		case uid := <-um.logout:
			um.handleLogout(uid)
		case ch := <-um.ctUser:
			ch <- len(um.rec)
		}
	}
}

func (um *ussnMgr) Reg(ussn *ussn) {
	um.reg <- ussn
}

func (um *ussnMgr) Unreg(ussn *ussn) {
	um.unreg <- ussn
}

func (um *ussnMgr) Peer(uid uid, msg interface{}) error {
	mp := newMsgUssnMgrPeer(uid, msg)
	um.peer <- mp
	return <-mp.chErr
}

func (um *ussnMgr) UpdateInfo(user *user) {
	um.update <- user
}

func (um *ussnMgr) Logout(uid uid) {
	um.logout <- uid
}

func (um *ussnMgr) CtUser() int {
	ch := make(chan int)
	um.ctUser <- ch
	return <-ch
}

func (um *ussnMgr) handleReg(ussn *ussn) {
	if prev, ok := um.rec[ussn.user.Id]; ok {
		prev.Logout(errors.New("kick login"))
	}
	um.rec[ussn.user.Id] = ussn
}

func (um *ussnMgr) handleUnreg(ussn *ussn) {
	if prev, ok := um.rec[ussn.user.Id]; ok && prev == ussn {
		delete(um.rec, ussn.user.Id)
	}
}

func (um *ussnMgr) handlePeer(mump *msgUssnMgrPeer) {
	if ussn, ok := um.rec[mump.to]; ok {
		go func() { // no block, thus go
			mump.chErr <- ussn.Write(mump.msg)
		}()
	} else {
		mump.chErr <- errors.New("ussn not in rec")
	}
}

func (um *ussnMgr) handleUpdate(user *user) {
	if ussn, ok := um.rec[user.Id]; ok {
		// no block, thus go
		go ussn.UpdateInfo(user)
	}
}

func (um *ussnMgr) handleLogout(uid uid) {
	if ussn, ok := um.rec[uid]; ok {
		// no block, thus go
		go ussn.Logout(errors.New("server kick"))
	}
}

