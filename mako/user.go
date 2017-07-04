package mako

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/mjpancake/ih/ako/model"
)

func Login(username, password string) (*model.User, error) {
	user := &model.User{}

	if db == nil {
		log.Fatalln("db is nil")
	}

	err := db.QueryRow(
		`select user_id, username, level, pt, rating
		from users where username=? && password=?`,
		username, hash(password)).
		Scan(&user.Id, &user.Username, &user.Level,
			&user.Pt, &user.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		log.Fatalln("db.Login", err)
	}

	return user, nil
}

func SignUp(username, password string) (*model.User, error) {
	if !checkName(username) {
		return nil, errors.New("用户名不可用")
	}

	var exist bool
	err := db.QueryRow(
		"select exists(select 1 from users where username=?)",
		username).Scan(&exist)

	if err != nil {
		log.Fatalln("db.SignUp", err)
	}

	if exist {
		return nil, errors.New("用户名已存在")
	}

	_, err = db.Exec(
		"insert into users (username, password) values (?,?)",
		username, hash(password))

	if err != nil {
		log.Fatalln("db.SignUp", err)
	}

	return Login(username, password)
}

func GetUser(uid model.Uid) *model.User {
	var user model.User

	err := db.QueryRow(
		`select user_id, username, level, pt, rating
		from users where user_id=?`, uid).
		Scan(&user.Id, &user.Username, &user.Level,
			&user.Pt, &user.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("cb.GetUser", err)
	}

	return &user
}

func GetUsers(uids *[4]model.Uid) [4]*model.User {
	var users [4]*model.User

	rows, err := db.Query(
		`select user_id, username, level, pt, rating
		from users where user_id in (?,?,?,?)`,
		uids[0], uids[1], uids[2], uids[3])
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.Id, &user.Username,
			&user.Level, &user.Pt, &user.Rating)
		if err != nil {
			log.Fatalln(err)
		}
		for w := 0; w < 4; w++ {
			if uids[w] == user.Id {
				users[w] = user
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return users
}

func checkName(name string) bool {
	if strings.HasPrefix(name, "ⓝ") {
		return false
	}
	return true
}

func hash(password string) string {
	sha := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(sha[:])
}
