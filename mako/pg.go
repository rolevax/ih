package mako

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
)

var db *pg.DB

func init() {
	db = pg.Connect(&pg.Options{
		Network:  "tcp",
		Addr:     "db:5432", // hostname in docker net
		User:     "postgres",
		Password: "",
		Database: "postgres",
	})

	for i := 0; i < 5; i++ {
		err := testConn()
		if err == nil {
			log.Println("mako.init pg: Ok")
			return
		}

		log.Println("mako.init pg:", err)
		time.Sleep(5 * time.Second)
	}

	log.Fatal("mako.init pg: tried too many times")
}

func testConn() error {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 0*COUNT(*)+42 FROM replays")
	if err != nil {
		return err
	}
	if n != 42 {
		return fmt.Errorf("evil db, 42 became %d", n)
	}
	return nil
}
