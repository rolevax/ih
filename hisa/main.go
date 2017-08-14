package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/nodoka/book"
	"github.com/rolevax/ih/nodoka/tssn"
	"github.com/rolevax/ih/nodoka/ussn"
)

const Version = "0.8.3"

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
	mako.AddAcceptingVersion(Version)
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
