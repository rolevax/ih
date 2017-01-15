package dao

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"bitbucket.org/rolevax/sakilogy-server/model"
)

type Dao struct {
	db		*sql.DB
}

func New() *Dao {
	var dao Dao

	db, err := sql.Open("mysql",
		"sakilogy:@k052a9@tcp(127.0.0.1:3306)/sakilogy")
	if err != nil {
		log.Fatal(err)
	}
	dao.db = db

	return &dao
}

func (dao *Dao) Close() {
	dao.db.Close()
}

func (dao *Dao) Auth(userAuth *model.UserAuth) *model.User {
	var user model.User
	var password string

	err := dao.db.QueryRow(
		`select id, username, nickname, password
		from users where username = ?`, userAuth.Username).
		Scan(&user.Id, &user.Username, &user.Nickname, &password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Dao.Auth", err)
	}

	if userAuth.Password != password {
		return nil
	}

	user.Conn = userAuth.Conn

	return &user
}


