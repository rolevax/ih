package main

import (
	"log"
	"net"
	"os"

	"github.com/mjpancake/hisa/actor"
	"github.com/mjpancake/hisa/db"
)

const Version = "0.8.0"

func main() {
	port := "6171"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	serve(port)
}

func serve(port string) {
	db.AddAcceptingVersion(Version)
	actor.StartAll()

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("hisa server", Version, "listen", port)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("Accept", err)
		} else {
			actor.Accept(conn)
		}
	}
}
