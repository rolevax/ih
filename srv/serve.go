package srv

import (
	"log"
	"net"
	"bufio"
	"encoding/json"
)

func Serve() {
	dao := newDao()
	defer dao.close()

	conns := newConns(dao)
	go conns.Loop()

	ln, err := net.Listen("tcp", ":6171")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("listen 6171")
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("E accept", err)
		} else {
			handle(conn, conns)
		}
	}
}

func handle(conn net.Conn, conns *conns) {
	breq, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("E main:handle", err)
		conn.Close()
		return
	}

	var req reqAuth

	if err := json.Unmarshal(breq, &req); err != nil {
		log.Println("E main:handle", err)
		conn.Close()
		return
	}

	switch req.Type {
	case "login":
		login := login{req.Version, req.Username, req.Password, conn}
		conns.Login() <- &login
	case "sign-up":
		sign := login{req.Version, req.Username, req.Password, conn}
		conns.SignUp() <- &sign
	default:
		log.Println("E main.handle unkown request", req.Type)
		conn.Close()
	}
}



