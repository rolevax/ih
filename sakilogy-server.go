package main

import (
	"log"
	"net"
	"bufio"
	"encoding/json"
	"bitbucket.org/rolevax/sakilogy-server/srv"
	"bitbucket.org/rolevax/sakilogy-server/dao"
	"bitbucket.org/rolevax/sakilogy-server/model"
)

func main() {
	dao := dao.New()
	defer dao.Close()

	conns := srv.NewConns(dao)
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

func handle(conn net.Conn, conns *srv.Conns) {
	breq, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("E main:handle", err)
		conn.Close()
		return
	}

	var req struct {
		Type		string
		Username	string
		Password	string
	}

	if err := json.Unmarshal(breq, &req); err != nil {
		log.Println("E main:handle", err)
		conn.Close()
		return
	}

	switch req.Type {
	case "login":
		login := model.Login{req.Username, req.Password, conn}
		conns.Login <- &login
	case "sign-up":
		sign := model.Login{req.Username, req.Password, conn}
		conns.SignUp <- &sign
	default:
		log.Println("E main.handle unkown request", req.Type)
		conn.Close()
	}
}



