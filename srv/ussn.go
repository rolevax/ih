package srv

import (
	"log"
	"net"
	"bufio"
	"errors"
	"strings"
	"time"
	"encoding/json"
)

const idleTimeOut = 15 * time.Minute

type ussn struct {
	user		user
	conn		net.Conn
	read		chan []byte
	write		chan *msgUssnWrite
	done		chan struct{}
	logout		chan error
	idleTimer	*time.Timer
}

func loopUssn(conn net.Conn) {
	defer conn.Close()

	ussn, err := startUssn(conn)
	if err != nil {
		reject(conn, newRespAuthFail(err.Error()))
		return
	}
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
			ussn.resetIdleTimer()
			ussn.switchRead(breq)
		case muw := <-ussn.write:
			muw.chErr <-ussn.send(muw.msg)
		case <-ussn.idleTimer.C:
			ussn.Logout(errors.New("idle timeout"))
		case err := <-ussn.logout:
			log.Println(ussn.user.Id, "----", err)
			close(ussn.done)
			return
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
	ussn.send(newRespAuthOk(&ussn.user))

	ussn.read = make(chan []byte)
	ussn.write = make(chan *msgUssnWrite)
	ussn.done = make(chan struct{})
	ussn.logout = make(chan error, 1) // sendable from same goroutine
	ussn.idleTimer = time.NewTimer(idleTimeOut)

	go ussn.readLoop()

	return ussn, nil
}

func authUssn(conn net.Conn) (*ussn, error) {
	breq, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var req reqAuth
	if err := json.Unmarshal(breq, &req); err != nil {
		return nil, err
	}

	if req.Version != Version {
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
	for {
		breq, err := bufio.NewReader(ussn.conn).ReadBytes('\n')
		if err != nil {
			ussn.Logout(err)
			return
		}

		select {
		case ussn.read <- breq:
			log.Print(ussn.user.Id, " ---> ", string(breq))
		case <-ussn.done:
			return
		}
	}
}

func (ussn *ussn) switchRead(breq []byte) {
	var req reqTypeOnly
	if err := json.Unmarshal(breq, &req); err != nil {
		log.Fatalln("ussn.readLoop", err)
	}
	t := req.Type

	switch {
	case t == "look-around":
		ussn.sendLookAround()
	case t == "book":
		sing.BookMgr.Book(ussn.user.Id)
	case t == "unbook":
		sing.BookMgr.Unbook(ussn.user.Id)
	case t == "ready":
		sing.TssnMgr.Ready(ussn.user.Id)
	case strings.HasPrefix(t, "t-"):
		var act reqAction
		if err := json.Unmarshal(breq, &act); err != nil {
			ussn.Logout(err)
			return
		}
		sing.TssnMgr.Action(ussn.user.Id, &act)
	}
}

func (ussn *ussn) send(msg interface{}) error {
    var jsonb []byte
    if str, ok := msg.(string); ok {
        jsonb = []byte(str)
    } else {
        var err error
        jsonb, err = json.Marshal(msg)
        if err != nil {
            log.Fatalln("ussn.send", err)
        }
    }

	_, err := ussn.conn.Write(append(jsonb, '\n'))
	if err != nil {
		ussn.Logout(err)
	} else {
		log.Println(ussn.user.Id, "<---", string(jsonb))
	}
	return err
}

func (ussn *ussn) sendLookAround() {
	bookable := !sing.TssnMgr.HasUser(ussn.user.Id)
	connCt := sing.UssnMgr.CtUser()

	playCt := sing.TssnMgr.CtUser()
	bookCt := sing.BookMgr.CtBook()
	ussn.send(newRespLookAround(bookable, connCt, bookCt, playCt))
}

func (ussn *ussn) resetIdleTimer() {
	if !ussn.idleTimer.Stop() {
		select {
		case <-ussn.idleTimer.C:
		default: // prevent blocked by double-draining
		}
	}
	ussn.idleTimer.Reset(idleTimeOut)
}

