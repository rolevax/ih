package srv

import (
	"log"
	"net"
)

const Version = "0.7.6"

var sing struct {
	Dao     *dao
	Rao     *rao
	UssnMgr *ussnMgr
	BookMgr *bookMgr
	TssnMgr *tssnMgr
}

func Serve(port string) {
	sing.Dao = newDao()
	sing.Rao = newRao()
	sing.UssnMgr = newUssnMgr()
	sing.BookMgr = newBookMgr()
	sing.TssnMgr = newTssnMgr()
	defer sing.Dao.Close()
	go sing.UssnMgr.Loop()
	go sing.BookMgr.Loop()
	go sing.TssnMgr.Loop()

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("hisa server", Version, "listen", port)
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
