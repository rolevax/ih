package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/node/book"
	"github.com/mjpancake/hisa/node/tssn"
	"github.com/mjpancake/hisa/node/ussn"
)

const Version = "0.8.2"

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	prefix := time.Now().Format("01/02 15:04:05")
	return fmt.Print(prefix, " ", string(bytes))
}

func main() {
	port := "6171"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	serve(port)
}

func serve(port string) {
	db.AddAcceptingVersion(Version)
	ussn.Init()
	tssn.Init()
	book.Init()

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
			ussn.Start(conn)
		}
	}
}
