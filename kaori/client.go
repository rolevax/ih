package main

import (
	"log"
	"net"

	"github.com/chzyer/readline"
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/hayari"
)

const (
	ClientVer = "0.9.1"
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

func login(rl *readline.Instance, username, password string) error {
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
		Version:  ClientVer,
		Username: username,
		Password: password,
	}
	cl.send(reqLogin)

	go readLoop(conn, rl)

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

func matchJoin(ruleId int) {
	cl.send(&cs.MatchJoin{
		RuleId: model.RuleId(ruleId),
	})
}

func getReplayList() {
	cl.send(&cs.GetReplayList{})
}

func getReplay(id uint) {
	cl.send(&cs.GetReplay{
		ReplayId: id,
	})
}

func readLoop(conn net.Conn, rl *readline.Instance) {
	for {
		b, err := hayari.Read(conn)
		if err != nil {
			rl.SetPrompt("\033[31mOFFLINE»\033[0m ")
			log.Println("LOGGED OUT", err)
			return
		}
		onRecv(b, rl)
	}
}
