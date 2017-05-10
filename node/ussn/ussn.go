package ussn

import (
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/netio"
	"github.com/mjpancake/hisa/node"
)

type ussn struct {
	p    *actor.PID
	user *model.User
	conn net.Conn
}

func Start(conn net.Conn) {
	ussn := &ussn{
		conn: conn,
	}
	props := actor.FromInstance(ussn)
	ussn.p = actor.Spawn(props)
}

func (ussn *ussn) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		if breq, err := netio.ReadAuth(ussn.conn); err != nil {
			ussn.handleReject(err)
		} else {
			ussn.p.Tell(breq)
		}
	case *actor.Stopping:
		ussn.bye()
	case *actor.Stopped:
	case *actor.Restarting:
	case error:
		ussn.handleError(msg)
	case []byte:
		if err := ussn.auth(msg); err != nil {
			ussn.handleReject(err)
		} else {
			ctx.SetBehavior(ussn.Work)
			ussn.hello()
		}
	default:
		log.Fatalf("ussn.Recv: unexpected %T\n", msg)
	}
}

func (ussn *ussn) Work(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
		ussn.bye()
	case *actor.Stopped:
	case *actor.Restarting:
	case error:
		ussn.handleError(msg)
	case []byte:
		ussn.handleRead(msg)
	case *pcSc:
		ussn.handleSc(msg.msg, makeResp(ctx))
	case *pcUpdateInfo:
		ussn.handleUpdateInfo()
	default:
		log.Fatalf("ussn.Work: unexpected %T\n", msg)
	}
}

func readLoop(conn net.Conn, succ func([]byte), fail func(error)) {
	for {
		breq, err := netio.Read(conn)
		if err != nil {
			if e, ok := err.(*net.OpError); ok {
				err = e.Err
			}
			fail(err)
			return
		}

		succ(breq)
	}
}

func (ussn *ussn) auth(breq []byte) error {
	var req model.CsAuth
	if err := json.Unmarshal(breq, &req); err != nil {
		return err
	}

	if !db.AcceptVersion(req.Version) {
		return errors.New("客户端版本过旧")
	}

	switch req.Type {
	case "login":
		u, err := db.Login(req.Username, req.Password)
		ussn.user = u
		return err
	case "sign-up":
		u, err := db.SignUp(req.Username, req.Password)
		ussn.user = u
		return err
	default:
		return errors.New("invalid auth req")
	}
}

func (ussn *ussn) hello() {
	log.Println(ussn.user.Id, "++++", ussn.conn.RemoteAddr())
	ussn.handleSc(model.NewScAuthOk(ussn.user), noResp)
	ussn.handleUpdateInfo()
	onRead := func(breq []byte) { ussn.p.Tell(breq) }
	onReadErr := func(err error) {
		node.Umgr.Tell(&node.MuKick{ussn.user.Id, err.Error()})
	}
	go readLoop(ussn.conn, onRead, onReadErr)
	node.Umgr.Tell(&cpReg{add: true, ussn: ussn})
}

func (ussn *ussn) bye() {
	ussn.conn.Close()
	node.Umgr.Tell(&cpReg{add: false, ussn: ussn})
}

func noResp(interface{}) {}

func makeResp(ctx actor.Context) func(interface{}) {
	return func(msg interface{}) {
		if ctx.Sender() != nil {
			ctx.Respond(msg)
		}
	}
}
