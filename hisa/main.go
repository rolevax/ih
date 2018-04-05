package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rolevax/ih/hitomi"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/nodoka/book"
	"github.com/rolevax/ih/nodoka/tssn"
	"github.com/rolevax/ih/nodoka/ussn"
	"github.com/rolevax/ih/ryuuka"
)

const Version = "0.9.4"

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	prefix := time.Now().Format("01/02 15:04:05")
	return fmt.Print(prefix, " ", string(bytes))
}

type initArgs struct {
	port       string
	redisAddr  string
	dbAddr     string
	ryuukaAddr string
	tokiAddr   string
	dictPath   string
}

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	if flag.Parsed() {
		log.Fatalln("unexpected flag parse before main()")
	}

	port := flag.String("port", "6171", "port to listen")
	redis := flag.String("redis", "localhost:6379", "redis server addr")
	db := flag.String("db", "localhost:5432", "pg db server addr")
	ryuuka := flag.String("ryuuka", "localhost:6172", "2nd addr to listen")
	toki := flag.String("toki", "localhost:8900", "toki server addr")
	dict := flag.String("dict", "/srv/dict.txt", "sersitive word dictionary")

	flag.Parse()

	serve(&initArgs{
		port:       *port,
		redisAddr:  *redis,
		dbAddr:     *db,
		ryuukaAddr: *ryuuka,
		tokiAddr:   *toki,
		dictPath:   *dict,
	})
}

func serve(args *initArgs) {
	mako.InitRedis(args.redisAddr)
	mako.InitDb(args.dbAddr)
	mako.AddAcceptingVersion(Version)

	hitomi.Init(args.dictPath)

	ryuuka.Init(args.ryuukaAddr, args.tokiAddr)

	ussn.Init()
	tssn.Init()
	book.Init()

	ln, err := net.Listen("tcp", ":"+args.port)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("hisa server", Version, "listen", args.port)
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
