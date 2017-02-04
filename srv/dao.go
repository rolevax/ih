package srv

import (
	"log"
	"database/sql"
	"errors"
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

func (dao *dao) Close() {
	dao.db.Close()
}

func (dao *dao) Login(username, password string) (*ussn, error) {
	ussn := new(ussn)

	err := dao.db.QueryRow(
		`select user_id, username
		from users where username=? && password=?`,
		username, password).
		Scan(&ussn.user.Id, &ussn.user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		log.Fatalln("dao.login", err)
	}

	return ussn, nil
}

func (dao *dao) SignUp(username, password string) (*ussn, error) {
	var exist bool
	err := dao.db.QueryRow(
		"select exists(select 1 from users where username=?)",
		username).Scan(&exist)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	if exist {
		return nil, errors.New("用户名已存在")
	}

	_, err = dao.db.Exec(
		"insert into users (username, password) values (?,?)",
		username, password)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	return dao.Login(username, password)
}

func (dao *dao) GetUser(uid uid) *user {
	user := new(user)

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

	return user
}

