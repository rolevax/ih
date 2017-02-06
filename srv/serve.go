package srv

import (
	"log"
	"net"
)

const Version = "0.6.7"

func versionCheck(v string) bool {
	// still support for 0.6.6 by now...
	return v == Version || v == "0.6.6"
}

var sing struct {
	Dao			*dao
	UssnMgr		*ussnMgr
	BookMgr		*bookMgr
	TssnMgr		*tssnMgr
}

func Serve(port string) {
	sing.Dao = newDao()
	sing.UssnMgr = newUssnMgr()
	sing.BookMgr = newBookMgr()
	sing.TssnMgr = newTssnMgr()
	defer sing.Dao.Close()
	go sing.UssnMgr.Loop()
	go sing.BookMgr.Loop()
	go sing.TssnMgr.Loop()

	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("sakilogy-server", Version, "listen", port)
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


