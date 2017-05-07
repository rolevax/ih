package ussn

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/mjpancake/hisa/actor/book"
	"github.com/mjpancake/hisa/actor/tssn/tbus"
	"github.com/mjpancake/hisa/actor/ussn/ubus"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/netio"
)

type ss struct {
	user   model.User
	conn   net.Conn
	read   chan []byte
	write  chan *msgUssnWrite
	update chan struct{}
	done   chan struct{}
	logout chan error
}

func LoopUssn(conn net.Conn) {
	defer conn.Close()

	ss, err := startUssn(conn)
	if err != nil {
		reject(conn, model.NewScAuthFail(err.Error()))
		return
	}
	log.Println(ss.user.Id, "++++", conn.RemoteAddr())
	Reg(ss)
	defer Unreg(ss)
	defer book.Unbook(ss.user.Id)

	for {
		select {
		case <-ss.done: // in prior
			return
		default:
		}

		select {
		case breq := <-ss.read:
			ss.handleRead(breq)
		case muw := <-ss.write:
			muw.chErr <- ss.handleWrite(muw.msg)
		case <-ss.update:
			ss.handleUpdateInfo()
		case err := <-ss.logout:
			ss.handleLogout(err)
		case <-ss.done:
			return
		}
	}
}

func startUssn(conn net.Conn) (*ss, error) {
	ss, err := authUssn(conn)
	if err != nil {
		return nil, err
	}

	ss.read = make(chan []byte)
	ss.write = make(chan *msgUssnWrite)
	ss.update = make(chan struct{})
	ss.done = make(chan struct{})
	ss.logout = make(chan error)

	ss.conn = conn
	ss.handleWrite(model.NewScAuthOk(&ss.user))
	ss.handleUpdateInfo()

	go ss.readLoop()

	return ss, nil
}

func authUssn(conn net.Conn) (*ss, error) {
	breq, err := netio.ReadAuth(conn)
	if err != nil {
		return nil, err
	}

	var req model.CsAuth
	if err := json.Unmarshal(breq, &req); err != nil {
		return nil, err
	}

	if !db.AcceptVersion(req.Version) {
		return nil, errors.New("客户端版本过旧")
	}

	switch req.Type {
	case "login":
		user, err := db.Login(req.Username, req.Password)
		if err != nil {
			return nil, err
		} else {
			return &ss{user: *user}, nil
		}
	case "sign-up":
		user, err := db.SignUp(req.Username, req.Password)
		if err != nil {
			return nil, err
		} else {
			return &ss{user: *user}, nil
		}
	default:
		return nil, errors.New("invalid auth req")
	}
}

func reject(conn net.Conn, msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("auth reject", err)
	}

	if err := netio.Write(conn, jsonb); err != nil {
		log.Println("auth reject", err)
	} else {
		log.Println(conn.RemoteAddr(), "<---", string(jsonb))
	}
}

type msgUssnWrite struct {
	msg   interface{}
	chErr chan error
}

func newMsgUssnWrite(msg interface{}) *msgUssnWrite {
	muw := new(msgUssnWrite)
	muw.msg = msg
	muw.chErr = make(chan error)
	return muw
}

func (ss *ss) Write(msg interface{}) error {
	muw := newMsgUssnWrite(msg)
	select {
	case ss.write <- muw:
		return <-muw.chErr
	case <-ss.done:
		return errors.New("ss done")
	}
}

func (ss *ss) UpdateInfo() {
	select {
	case ss.update <- struct{}{}:
	case <-ss.done:
	}
}

func (ss *ss) Logout(err error) {
	if err == nil {
		log.Fatalln("logout with nil err")
	}
	select {
	case ss.logout <- err:
	case <-ss.done:
	}
}

func (ss *ss) readLoop() {
	for {
		breq, err := netio.Read(ss.conn)
		if err != nil {
			ss.Logout(err) // ok, not in ss main goroutine
			return
		}

		select {
		case ss.read <- breq:
			//log.Print(ss.user.Id, " ---> ", string(breq))
		case <-ss.done:
			return
		}
	}
}

func (ss *ss) handleRead(breq []byte) {
	var req model.ScTypeOnly
	if err := json.Unmarshal(breq, &req); err != nil {
		log.Fatalln("ss.readLoop", err)
	}
	t := req.Type

	switch {
	case t == "look-around":
		ss.handleLookAround()
	case t == "heartbeat":
		// do nothing
	case t == "book":
		var req model.CsBook
		if err := json.Unmarshal(breq, &req); err != nil {
			ss.handleLogout(err)
			return
		}
		if !req.BookType.Valid() {
			ss.handleLogout(errors.New("invalid bktype"))
			return
		}
		book.Book(ss.user.Id, req.BookType)
	case t == "unbook":
		book.Unbook(ss.user.Id)
	case t == "ready":
		tbus.Ready(ss.user.Id)
	case t == "choose":
		var req model.CsChoose
		if err := json.Unmarshal(breq, &req); err != nil {
			ss.handleLogout(err)
			return
		}
		tbus.Choose(ss.user.Id, req.GirlIndex)
	case t == "t-action":
		var act model.CsAction
		if err := json.Unmarshal(breq, &act); err != nil {
			ss.handleLogout(err)
			return
		}
		tbus.Action(ss.user.Id, &act)
	case t == "get-replay":
		var req model.CsGetReplay
		if err := json.Unmarshal(breq, &req); err != nil {
			ss.handleLogout(err)
			return
		}
		time.Sleep(2 * time.Second)
		ss.handleGetReplay(req.ReplayId)
	default:
		ss.handleLogout(errors.New("invalid req: " + string(breq)))
	}
}

func (ss *ss) handleWrite(msg interface{}) error {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("ss.handleWrite marshal", err)
		return err
	}

	err = netio.Write(ss.conn, jsonb)
	if err != nil {
		ss.handleLogout(err)
	} else {
		//log.Println(ss.user.Id, "<---", string(jsonb))
	}
	return err
}

func (ss *ss) handleLogout(err error) {
	log.Println(ss.user.Id, "----", err)
	close(ss.done)
}

func (ss *ss) handleLookAround() {
	if tbus.HasUser(ss.user.Id) {
		msg := model.ScTypeOnly{"resume"}
		ss.handleWrite(msg)
	} else {
		connCt := ubus.CtUser()
		msg := model.NewScLookAround(connCt)
		pss := tbus.CtEachBt()
		bss := book.CtBooks()
		user := ss.user
		cBookable := user.Level >= 9
		bBookable := user.Level >= 13 && user.Rating >= 1800.0
		dBookable := !bBookable
		aBookable := user.Level >= 16 && user.Rating >= 2000.0
		msg.Books[0] = model.BookEntry{dBookable, bss[0].Wait, 4 * pss[0]}
		msg.Books[1] = model.BookEntry{cBookable, bss[1].Wait, 4 * pss[1]}
		msg.Books[2] = model.BookEntry{bBookable, bss[2].Wait, 4 * pss[2]}
		msg.Books[3] = model.BookEntry{aBookable, bss[3].Wait, 4 * pss[3]}
		ss.handleWrite(msg)
	}
}

func (ss *ss) handleUpdateInfo() {
	ss.user = *db.GetUser(ss.user.Id)
	ss.handleWrite(model.NewScUpdateUser(&ss.user, db.GetStats(ss.user.Id)))
}

func (ss *ss) handleGetReplay(replayId uint) {
	text, err := db.GetReplay(replayId)
	if err != nil {
		ss.handleLogout(err)
	}
	ss.handleWrite(model.NewScGetReplay(replayId, text))
}
