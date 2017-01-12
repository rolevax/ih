package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

func main() {
    session := saki.NewTableSession()
	sv := session.Action(0, 3)
	fmt.Println("action result 0", sv.Get(0))

	ln, _ := net.Listen("tcp", ":6171")
	fmt.Println("sakilogy server listening at 6171")

	conn, _ := ln.Accept()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Read error: ", err.Error())
			break
		}
		fmt.Print("Message Received:", string(message))
		back := strings.ToUpper(message)
		conn.Write([]byte(back + "\n"))
	}
}

