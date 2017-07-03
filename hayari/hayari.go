package hayari

import (
	"net"
	"time"

	"github.com/rolevax/sp4g"
)

const readAuthTimeOut = 5 * time.Second
const idleTimeOut = 15 * time.Minute
const writeTimeOut = 10 * time.Second

func Read(conn net.Conn) ([]byte, error) {
	return ReadTime(conn, idleTimeOut)
}

func ReadAuth(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(readAuthTimeOut))
	return sp4g.ReadN(conn, 1024)
}

func ReadTime(conn net.Conn, out time.Duration) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(out))
	return sp4g.Read(conn)
}

func Write(conn net.Conn, data []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
	return sp4g.Write(conn, data)
}
