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
	Auth	chan *model.UserAuth
	dao		*dao.Dao
	users	map[int]*model.User
	books	*Books
}

func NewConns(dao *dao.Dao) *Conns {
	var conns Conns

	conns.Auth = make(chan *model.UserAuth)
	conns.dao = dao
	conns.users = make(map[int]*model.User)
	conns.books = NewBooks()

	return &conns
}

func (conns *Conns) Loop() {
	for {
		select {
		case userAuth := <-conns.Auth:
			user := conns.dao.Auth(userAuth)
			if user != nil {
				conns.add(user)
			} else {
				conns.bye(userAuth.Conn, msgAuthFail)
			}
		}
	}
}

func (conns *Conns) add(user *model.User) {
	// all existing session kicked
	conns.users[user.Id] = user
	conns.send(user.Id, newAuthOkMsg(user))

	go conns.readLoop(user)
}

func (conns *Conns) readLoop(user *model.User) {
	conn := user.Conn
	defer conn.Close()
	for {
		breq, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println(user.Id, "----")
			} else {
				log.Println("E Conns.readLoop", err)
			}
			return
		}

		var req struct {
			Type		string
		}

		if err := json.Unmarshal(breq, &req); err != nil {
			log.Fatal("E Conns.readLoop", err)
			return
		}

		switch req.Type {
		case "book":
			log.Println("=== book!!! by", user.Username)
		}
	}
}

func (conns *Conns) send(uid int, msg interface{}) {
	user, found := conns.users[uid]
	if !found {
		log.Println("E Conns.send user", uid, "not found")
		return
	}
	conn := user.Conn

	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Conns.send", err)
	}

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("Conns.send", err)
	} else {
		log.Println(uid, "<---", string(jsonb))
	}
}

func (conns *Conns) bye(conn net.Conn, msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Conns.bye", err)
	}

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("Conns.bye", err)
	} else {
		log.Println(conn.RemoteAddr(), "<---", string(jsonb))
	}

	conn.Close()
}



/// messages

var msgAuthFail = struct {
	Type	string
	Ok		bool
	Reason	string
}{"auth", false, "用户名或密码错误"}

func newAuthOkMsg(user *model.User) interface{} {
	return struct {
		Type	string
		Ok		bool
		User	*model.User
	}{"auth", true, user}
}

