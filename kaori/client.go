package main

import (
	"log"
	"net"

	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/hayari"
)

type client struct {
	conn net.Conn
}

var cl *client

func (cl *client) send(msg interface{}) {
	if cl == nil {
		log.Println("offline")
		return
	}

	hayari.Write(cl.conn, cs.ToJson(msg))
}

func login(username, password string) error {
	if cl != nil {
		logout()
	}

	cl = &client{}

	conn, err := net.Dial("tcp", "127.0.0.1:6171")
	if err != nil {
		return err
	}
	cl.conn = conn

	reqLogin := &cs.Auth{
		Version:  "0.8.3",
		Username: username,
		Password: password,
	}
	cl.send(reqLogin)

	go readLoop(conn)

	return nil
}

func logout() {
	if cl != nil {
		cl.conn.Close()
		cl = nil
	}
}

func lookAround() {
	cl.send(&cs.LookAround{})
}

func roomCreate() {
	cl.send(&cs.RoomCreate{
		AiNum:  model.Ai3,
		Bans:   []model.Gid{},
		AiGids: []model.Gid{0, 0, 0},
	})
}

func roomJoin(rid int) {
	cl.send(&cs.RoomJoin{
		RoomId: model.Rid(rid),
	})
}

func roomQuit() {
	cl.send(&cs.RoomQuit{})
}

func getReplayList() {
	cl.send(&cs.GetReplayList{})
}

func getReplay(id uint) {
	cl.send(&cs.GetReplay{
		ReplayId: id,
	})
}

func readLoop(conn net.Conn) {
	for {
		b, err := hayari.Read(conn)
		if err != nil {
			log.Println("LOGGED OUT", err)
			return
		}
		onRecv(b)
	}
}
