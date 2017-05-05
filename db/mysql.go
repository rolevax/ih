package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// just enjoying hard coding password
	var err error
	db, err = sql.Open("mysql",
		"sakilogy:@k052a9@tcp(127.0.0.1:3306)/sakilogy")
	if err != nil {
		log.Fatalln(err)
	}

	if db.Ping() != nil {
		log.Fatalln("ping DB failed", err)
	}
}
