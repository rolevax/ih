package actor

import (
	"net"

	"github.com/mjpancake/hisa/actor/book"
	"github.com/mjpancake/hisa/actor/tssn"
	"github.com/mjpancake/hisa/actor/ussn"
)

func StartAll() {
	go ussn.Loop()
	go tssn.Loop()
	go book.Loop()
}

func Accept(conn net.Conn) {
	go ussn.LoopUssn(conn)
}
