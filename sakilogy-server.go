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
	case "fetch-ann":
		reply := struct {
			Type	string
			Ann		string
			Login	bool
		} { "fetch-ann", "[公告]服务器正在测试", true }
		jsonb, _ := json.Marshal(reply)
		if _, err := conn.Write(append(jsonb, '\n')); err != nil {
			log.Println("E main:handle write ann", err)
		} else {
			log.Println(conn.RemoteAddr(), "<--- announcement")
		}
		conn.Close()
	case "login":
		userAuth := model.UserAuth{req.Username, req.Password, conn}
		conns.Auth <- &userAuth
	default:
		log.Println("E main.handle unkown request", req.Type)
		conn.Close()
	}
}



