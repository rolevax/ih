package srv

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

type dao struct {
	db		*sql.DB
}

func newDao() *dao {
	dao := new(dao)

	db, err := sql.Open("mysql",
		"sakilogy:@k052a9@tcp(127.0.0.1:3306)/sakilogy")
	if err != nil {
		log.Fatalln(err)
	}

	if db.Ping() != nil {
		log.Fatalln("ping DB failed", err)
	}

	dao.db = db

	return dao
}

func (dao *dao) close() {
	dao.db.Close()
}

func (dao *dao) login(login *login) *user {
	var user user

	err := dao.db.QueryRow(
		`select user_id, username
		from users where username=? && password=?`,
		login.Username, login.Password).
		Scan(&user.Id, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("dao.login", err)
	}

	return &user
}

func (dao *dao) signUp(sign *login) *user {
	var exist bool
	err := dao.db.QueryRow(
		"select exists(select 1 from users where username=?)",
		sign.Username).Scan(&exist)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	if exist {
		return nil
	}

	_, err = dao.db.Exec(
		"insert into users (username, password) values (?,?)",
		sign.Username, sign.Password)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	return dao.login(sign)
}

func (dao *dao) getUser(uid uid) *user {
	var user user

	err := dao.db.QueryRow(
		`select user_id, username 
		from users where user_id=?`, uid).
		Scan(&user.Id, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("dao.GetUser", err)
	}

	return &user
}

