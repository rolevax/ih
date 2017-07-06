package mako

import (
	"log"

	"github.com/go-pg/pg"
)

var db *pg.DB

func init() {
	db = pg.Connect(&pg.Options{
		Network:  "unix",
		User:     "mako",
		Password: "",
		Database: "mako",
	})

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 42")
	if err != nil {
		log.Fatal("mako.init:", err)
	}
	if n != 42 {
		log.Fatal("mako.init: evil db")
	}
}
