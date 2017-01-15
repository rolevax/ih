package srv

import (
	"log"
	"net"
	"bufio"
	"io"
	"encoding/json"
	"bitbucket.org/rolevax/sakilogy-server/model"
	"bitbucket.org/rolevax/sakilogy-server/dao"
)

type Conns struct {
	Login	chan *model.Login
	Logout	chan model.Uid
	Start	chan [4]model.Uid
	dao		*dao.Dao
	users	map[model.Uid]*model.User
	conns	map[model.Uid]net.Conn
	books	*Books
	tables	*Tables
}

func NewConns(dao *dao.Dao) *Conns {
	var conns Conns

	conns.Login = make(chan *model.Login)
	conns.Logout = make(chan model.Uid)
	conns.Start = make(chan [4]model.Uid)
	conns.dao = dao
	conns.users = make(map[model.Uid]*model.User)
	conns.conns = make(map[model.Uid]net.Conn)
	conns.books = NewBooks(&conns)
	conns.tables = NewTables(&conns)

	return &conns
}

func (conns *Conns) Loop() {
	go conns.books.Loop()
	go conns.tables.Loop()

	for {
		select {
		case login := <-conns.Login:
			user := conns.dao.Login(login)
			if user != nil {
				conns.add(user, login.Conn)
			} else {
				conns.reject(login.Conn, msgLoginFail)
			}
		case uid := <-conns.Logout:
			conns.logout(uid)
		case uids := <-conns.Start:
			log.Println("send create to table")
			conns.tables.Create <- uids
		}
	}
}

func (conns *Conns) add(user *model.User, conn net.Conn) {
	// all existing session kicked
	conns.users[user.Id] = user
	conns.conns[user.Id] = conn
	conns.send(user.Id, newLoginOkMsg(user))

	go conns.readLoop(user.Id)
}

func (conns *Conns) logout(uid model.Uid) {
	conn, found := conns.conns[uid]
	if found {
		conns.books.Unbook <- uid
		log.Println(uid, "----")
		conn.Close()
	}

	delete(conns.conns, uid)
	delete(conns.users, uid)
}

func (conns *Conns) readLoop(uid model.Uid) {
	conn := conns.conns[uid]
	for {
		breq, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				conns.Logout <- uid
			} else {
				log.Println("E Conns.readLoop", err)
			}
			return
		}

		log.Print(uid, "--->", string(breq))

		var req struct {
			Type		string
		}

		if err := json.Unmarshal(breq, &req); err != nil {
			log.Fatal("E Conns.readLoop", err)
			return
		}

		switch req.Type {
		case "book":
			conns.books.Book <- uid
		case "unbook":
			conns.books.Unbook <- uid
		}
	}
}

func (conns *Conns) send(uid model.Uid, msg interface{}) {
	conn, found := conns.conns[uid]
	if !found {
		log.Println("E Conns.send user", uid, "not found")
		return
	}

	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Conns.send", err)
	}

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("Conns.send", err)
	} else {
		log.Println(uid, " <--- ", string(jsonb))
	}
}

func (conns *Conns) reject(conn net.Conn, msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Conns.reject", err)
	}

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("Conns.reject", err)
	} else {
		log.Println(conn.RemoteAddr(), "<---", string(jsonb))
	}

	conn.Close()
}



/// messages

var msgLoginFail = struct {
	Type	string
	Ok		bool
	Reason	string
}{"auth", false, "用户名或密码错误"}

func newLoginOkMsg(user *model.User) interface{} {
	return struct {
		Type	string
		Ok		bool
		User	*model.User
	}{"auth", true, user}
}

