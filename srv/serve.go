package srv

import (
	"log"
	"net"
)

const Version = "0.6.5"

var sing struct {
	Dao			*dao
	UssnMgr		*ussnMgr
	BookMgr		*bookMgr
	TssnMgr		*tssnMgr
}

func Serve() {
	sing.Dao = newDao()
	sing.UssnMgr = newUssnMgr()
	sing.BookMgr = newBookMgr()
	sing.TssnMgr = newTssnMgr()
	defer sing.Dao.Close()
	go sing.UssnMgr.Loop()
	go sing.BookMgr.Loop()
	go sing.TssnMgr.Loop()

	ln, err := net.Listen("tcp", ":6171")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("sakilogy-server", Version, "listen 6171")
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("E accept", err)
		} else {
			go loopUssn(conn)
		}
	}
}


