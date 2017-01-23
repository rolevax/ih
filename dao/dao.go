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
		log.Fatalln(err)
	}

	if db.Ping() != nil {
		log.Fatalln("ping DB failed", err)
	}

	dao.db = db

	return &dao
}

func (dao *Dao) Close() {
	dao.db.Close()
}

func (dao *Dao) Login(login *model.Login) *model.User {
	var user model.User

	err := dao.db.QueryRow(
		`select user_id, username
		from users where username=? && password=?`,
		login.Username, login.Password).
		Scan(&user.Id, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("Dao.Login", err)
	}

	return &user
}

func (dao *Dao) SignUp(sign *model.Login) *model.User {
	var exist bool
	err := dao.db.QueryRow(
		"select exists(select 1 from users where username=?)",
		sign.Username).Scan(&exist)

	if err != nil {
		log.Fatalln("Dao.SignUp", err)
	}

	if exist {
		return nil
	}

	_, err = dao.db.Exec(
		"insert into users (username, password) values (?,?)",
		sign.Username, sign.Password)

	if err != nil {
		log.Fatalln("Dao.SignUp", err)
	}

	return dao.Login(sign)
}

func (dao *Dao) GetUser(uid model.Uid) *model.User {
	var user model.User

	err := dao.db.QueryRow(
		`select user_id, username 
		from users where user_id=?`, uid).
		Scan(&user.Id, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("Dao.GetUser", err)
	}

	return &user
}


