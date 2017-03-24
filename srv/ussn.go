package srv

import (
	"log"
	"net"
	"bufio"
	"errors"
	"time"
	"encoding/json"
)

const readAuthTimeOut = 10 * time.Second
const idleTimeOut = 15 * time.Minute
const writeTimeOut = 10 * time.Second
const obayTimeOut = 5 * time.Second

type ussn struct {
	user		user
	conn		net.Conn
	read		chan []byte
	write		chan *msgUssnWrite
	update		chan *user
	done		chan struct{}
	logout		chan error
}

func loopUssn(conn net.Conn) {
	defer conn.Close()

	ussn, err := startUssn(conn)
	if err != nil {
		reject(conn, newRespAuthFail(err.Error()))
		return
	}
	log.Println(ussn.user.Id, "++++", conn.RemoteAddr())
	sing.UssnMgr.Reg(ussn)
	defer sing.UssnMgr.Unreg(ussn)
	defer sing.BookMgr.Unbook(ussn.user.Id)

	for {
		select {
		case <-ussn.done: // in prior
			return
		default:
		}

		select {
		case breq := <-ussn.read:
			ussn.handleRead(breq)
		case muw := <-ussn.write:
			muw.chErr <-ussn.handleWrite(muw.msg)
		case user:= <-ussn.update:
			ussn.handleUpdateInfo(user)
		case err := <-ussn.logout:
			ussn.handleLogout(err)
		case <-ussn.done:
			return
		}
	}
}

func startUssn(conn net.Conn) (*ussn, error) {
	ussn, err := authUssn(conn)
	if err != nil {
		return nil, err
	}

	ussn.conn = conn
	ussn.handleWrite(newRespAuthOk(&ussn.user))

	ussn.read = make(chan []byte)
	ussn.write = make(chan *msgUssnWrite)
	ussn.update = make(chan *user)
	ussn.done = make(chan struct{})
	ussn.logout = make(chan error)

	go ussn.readLoop()

	return ussn, nil
}

func authUssn(conn net.Conn) (*ussn, error) {
	conn.SetReadDeadline(time.Now().Add(readAuthTimeOut))
	breq, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var req reqAuth
	if err := json.Unmarshal(breq, &req); err != nil {
		return nil, err
	}

	if !sing.Rao.AcceptVersion(req.Version) {
		return nil, errors.New("客户端版本过旧")
	}

	switch req.Type {
	case "login":
		return sing.Dao.Login(req.Username, req.Password)
	case "sign-up":
		return sing.Dao.SignUp(req.Username, req.Password)
	default:
		return nil, errors.New("invalid auth req")
	}
}

func reject(conn net.Conn, msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("auth reject", err)
	}

	conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("auth reject", err)
	} else {
		log.Println(conn.RemoteAddr(), "<---", string(jsonb))
	}
}

type msgUssnWrite struct {
	msg		interface{}
	chErr	chan error
}

func newMsgUssnWrite(msg interface{}) *msgUssnWrite {
	muw := new(msgUssnWrite)
	muw.msg = msg
	muw.chErr = make(chan error)
	return muw
}

func (ussn *ussn) Write(msg interface{}) error {
	muw := newMsgUssnWrite(msg)
	select {
	case ussn.write <- muw:
		return <-muw.chErr
	case <-ussn.done:
		return errors.New("ussn done")
	}
}

func (ussn *ussn) UpdateInfo(user *user) {
	select {
	case ussn.update <- user:
	case <-ussn.done:
	}
}

func (ussn *ussn) Logout(err error) {
	if err == nil {
		log.Fatalln("logout with nil err")
	}
	select {
	case ussn.logout <- err:
	case <-ussn.done:
	}
}

func (ussn *ussn) readLoop() {
	reader := bufio.NewReader(ussn.conn)
	for {
		ussn.conn.SetReadDeadline(time.Now().Add(idleTimeOut))
		breq, err := reader.ReadBytes('\n')
		if err != nil {
			ussn.Logout(err) // ok, not in ussn main goroutine
			return
		}

		select {
		case ussn.read <- breq:
			//log.Print(ussn.user.Id, " ---> ", string(breq))
		case <-ussn.done:
			return
		}
	}
}

func (ussn *ussn) handleRead(breq []byte) {
	var req reqTypeOnly
	if err := json.Unmarshal(breq, &req); err != nil {
		log.Fatalln("ussn.readLoop", err)
	}
	t := req.Type

	switch {
	case t == "look-around":
		ussn.handleLookAround()
	case t == "heartbeat":
		// do nothing
	case t == "book":
		var req reqBook
		if err := json.Unmarshal(breq, &req); err != nil {
			ussn.handleLogout(err)
			return
		}
		if !req.BookType.valid() {
			ussn.handleLogout(errors.New("invalid bktype"))
			return
		}
		sing.BookMgr.Book(ussn.user.Id, req.BookType)
	case t == "unbook":
		sing.BookMgr.Unbook(ussn.user.Id)
	case t == "ready":
		sing.TssnMgr.Ready(ussn.user.Id)
	case t == "choose":
		var req reqChoose
		if err := json.Unmarshal(breq, &req); err != nil {
			ussn.handleLogout(err)
			return
		}
		sing.TssnMgr.Choose(ussn.user.Id, req.GirlIndex)
	case t == "t-action":
		var act reqAction
		if err := json.Unmarshal(breq, &act); err != nil {
			ussn.handleLogout(err)
			return
		}
		sing.TssnMgr.Action(ussn.user.Id, &act)
	default:
		ussn.handleLogout(errors.New("invalid req: " + string(breq)))
	}
}

func (ussn *ussn) handleWrite(msg interface{}) error {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("ussn.handleWrite marshal", err)
		return err
	}

	ussn.conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
	_, err = ussn.conn.Write(append(jsonb, '\n'))
	if err != nil {
		ussn.handleLogout(err)
	} else {
		//log.Println(ussn.user.Id, "<---", string(jsonb))
	}
	return err
}

func (ussn *ussn) handleLogout(err error) {
	log.Println(ussn.user.Id, "----", err)
	close(ussn.done)
}

func (ussn *ussn) handleLookAround() {
	if sing.TssnMgr.HasUser(ussn.user.Id) {
		msg := respTypeOnly{"resume"}
		ussn.handleWrite(msg)
	} else {
		connCt := sing.UssnMgr.CtUser()
		msg := newRespLookAround(connCt)
		pss := sing.TssnMgr.CtEachBt()
		bss := sing.BookMgr.CtBooks()
		user := ussn.user
		cBookable := user.Level >= 9
		bBookable := user.Level >= 13 && user.Rating >= 1800.0
		dBookable := !bBookable
		aBookable := user.Level >= 16 && user.Rating >= 2000.0
		msg.Books[0] = bookEntry{dBookable, bss[0].wait, 4 * pss[0]}
		msg.Books[1] = bookEntry{cBookable, bss[1].wait, 4 * pss[1]}
		msg.Books[2] = bookEntry{bBookable, bss[2].wait, 4 * pss[2]}
		msg.Books[2] = bookEntry{aBookable, bss[3].wait, 4 * pss[3]}
		ussn.handleWrite(msg)
	}
}

func (ussn *ussn) handleUpdateInfo(user *user) {
	ussn.user = *user
	ussn.handleWrite(newRespUpdateUser(user))
}

