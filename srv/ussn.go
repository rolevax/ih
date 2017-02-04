package srv

import (
	"log"
	"net"
	"bufio"
	"errors"
	"strings"
	"encoding/json"
)

type ussn struct {
	user	user
	conn	net.Conn
	read	chan []byte
	write	chan *msgUssnWrite
	logout	chan error
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
		case breq := <-ussn.read:
			ussn.switchRead(breq)
		case muw := <-ussn.write:
			muw.chErr <- ussn.send(muw.msg)
		case err := <-ussn.logout:
			log.Println(ussn.user.Id, "----", err)
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
	ussn.logout = make(chan error)

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
	ussn.write <- muw
	return <-muw.chErr
}

func (ussn *ussn) Logout(err error) {
	ussn.logout <- err
}

func (ussn *ussn) readLoop() {
	for {
		breq, err := bufio.NewReader(ussn.conn).ReadBytes('\n')
		if err != nil {
			ussn.Logout(err)
			return
		}

		log.Print(ussn.user.Id, " ---> ", string(breq))
		ussn.read <- breq
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
