package srv

import (
	"log"
	"net"
	"bufio"
	"io"
	"strings"
	"encoding/json"
	"bitbucket.org/rolevax/sakilogy-server/model"
	"bitbucket.org/rolevax/sakilogy-server/dao"
)

type Conns struct {
	Login	chan *model.Login
	Logout	chan model.Uid
	Start	chan [4]model.Uid
    Peer    chan *Mail
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
	conns.Peer = make(chan *Mail)

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
				str := "用户名或密码错误"
				conns.reject(login.Conn, newLoginFailMsg(str))
			}
		case uid := <-conns.Logout:
			conns.logout(uid)
		case uids := <-conns.Start:
			conns.tables.Create <- uids
		case mail := <-conns.Peer:
            conns.send(mail.To, mail.Msg)
		}
	}
}

func (conns *Conns) add(user *model.User, conn net.Conn) {
	// prevent dup login
	if _, ok := conns.users[user.Id]; ok {
		str := "该用户已登录"
		conns.reject(conn, newLoginFailMsg(str));
		return
	}

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

		log.Print(uid, " ---> ", string(breq))
		var req struct {Type string}
		if err := json.Unmarshal(breq, &req); err != nil {
			log.Fatal("E Conns.readLoop", err)
			return
		}
		conns.switchRead(uid, req.Type, breq)
	}
}

func (conns *Conns) switchRead(uid model.Uid, t string, breq []byte) {
	switch {
	case t == "look-around":
		conns.sendLookAround(uid)
	case t == "book":
		conns.books.Book <- uid
	case t == "unbook":
		conns.books.Unbook <- uid
	case t == "ready":
		conns.tables.Ready <- uid
	case strings.HasPrefix(t, "t-"):
		act := Action{Uid: uid}
		if err := json.Unmarshal(breq, &act); err != nil {
			log.Println("E Conns.switchRead", err)
			return
		}
		conns.tables.Action <- &act
	}
}

func (conns *Conns) send(uid model.Uid, msg interface{}) {
	conn, found := conns.conns[uid]
	if !found {
		log.Println("E Conns.send user", uid, "not found")
		return
	}

    var jsonb []byte
    if str, ok := msg.(string); ok {
        jsonb = []byte(str)
    } else {
        var err error
        jsonb, err = json.Marshal(msg)
        if err != nil {
            log.Fatalln("Conns.send", err)
        }
    }

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("Conns.send", err)
	} else {
		log.Println(uid, "<---", string(jsonb))
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

func (conns *Conns) sendLookAround(uid model.Uid) {
	connCt := len(conns.conns)
	playCt := 4 * conns.tables.SessionCount();
	idleCt := connCt - playCt;
	bookCt := conns.books.BookCount()

	msg := struct {
		Type	string
		Conn	int
		Idle	int
		Book	int
		Play	int
	}{"look-around", connCt, idleCt, bookCt, playCt}
	conns.send(uid, msg)
}



/// messages

func newLoginFailMsg(str string) interface{} {
	return struct {
		Type	string
		Ok		bool
		Reason	string
	}{"auth", false, str}
}

func newLoginOkMsg(user *model.User) interface{} {
	return struct {
		Type	string
		Ok		bool
		User	*model.User
	}{"auth", true, user}
}

