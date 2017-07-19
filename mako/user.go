package mako

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/go-pg/pg"
	"github.com/mjpancake/ih/ako/model"
)

func Login(username, password string) (*model.User, error) {
	user := &model.User{}

	if db == nil {
		log.Fatalln("mako.Login: db is nil")
	}

	err := db.Model(user).
		Where("username=?", username).
		Where("password=?", hash(password)).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		log.Fatalln("mako.Login", err)
	}

	return user, nil
}

func SignUp(username, password string) error {
	if !checkName(username) {
		return errors.New("用户名不可用")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("db.SignUp", err)
	}

	var exist bool
	_, err = tx.QueryOne(
		&exist,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=?)",
		username,
	)

	if err != nil {
		tx.Rollback()
		log.Fatalln("db.SignUp", err)
	}

	if exist {
		tx.Rollback()
		return errors.New("用户名已存在")
	}

	var uid model.Uid
	_, err = tx.QueryOne(
		&uid,
		"INSERT INTO users(username, password) VALUES (?,?) RETURNING user_id",
		username, hash(password),
	)

	if err != nil {
		tx.Rollback()
		log.Fatalln("db.SignUp", err)
	}

	_, err = tx.Exec(
		"INSERT INTO user_girl(user_id, girl_id) VALUES (?, 0)",
		uid,
	)

	if err != nil {
		tx.Rollback()
		log.Fatalln("db.SignUp", err)
	}

	return nil
}

func GetUser(uid model.Uid) *model.User {
	user := &model.User{
		Id: uid,
	}

	err := db.Select(user)

	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		log.Fatalln("mako.GetUser", err)
	}

	return user
}

func GetUsers(uids *[4]model.Uid) [4]*model.User {
	var users [4]*model.User

	_, err := db.Query(
		&users,
		`SELECT user_id, username, level, pt, rating
		FROM users WHERE user_id in (?)
		ORDER BY user_id=? DESC,
		         user_id=? DESC,
		         user_id=? DESC,
		         user_id=? DESC`,
		pg.In(uids), uids[0], uids[1], uids[2], uids[3],
	)
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