package srv

import (
	"errors"
)

type ussnMgr struct {
	rec		map[uid]*ussn
	reg		chan *ussn
	unreg	chan *ussn
    peer    chan *msgUssnMgrPeer
	logout	chan uid
}

func newUssnMgr() *ussnMgr {
	um := new(ussnMgr)

	um.rec = make(map[uid]*ussn)
	um.reg = make(chan *ussn)
	um.unreg = make(chan *ussn)
    um.peer = make(chan *msgUssnMgrPeer)
	um.logout = make(chan uid)

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
		case uid := <-um.logout:
			um.handleLogout(uid)
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

func (um *ussnMgr) Logout(uid uid) {
	um.logout <- uid
}

func (um *ussnMgr) CtUser() int {
	return len(um.rec)
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
		mump.chErr <- ussn.Write(mump.msg)
	} else {
		mump.chErr <- errors.New("ussn not in rec")
	}
}

func (um *ussnMgr) handleLogout(uid uid) {
	if ussn, ok := um.rec[uid]; ok {
		ussn.Logout(errors.New("server kick"))
	}
}

