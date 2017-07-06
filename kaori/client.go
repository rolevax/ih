package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/mjpancake/ih/ako/cs"
	"github.com/mjpancake/ih/hayari"
)

type client struct {
	conn net.Conn
}

var cl *client

func (cl *client) send(cs interface{}) {
	jsonb, _ := json.Marshal(cs)
	hayari.Write(cl.conn, jsonb)
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
		Version:  "0.8.2",
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

func getReplayList() {
	if cl == nil {
		fmt.Println("offline")
		return
	}

	cl.send(&struct{ Type string }{Type: "get-replay-list"})
}

func getReplay(id uint) {
	if cl == nil {
		fmt.Println("offline")
		return
	}

	cl.send(&struct {
		Type     string
		ReplayId uint
	}{
		Type:     "get-replay",
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
		log.Println(string(b))
	}
}
