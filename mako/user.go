package mako

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/rolevax/ih/ako/model"
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

	tx.Commit()

	return nil
}

// unused yet
func Activate(username, password, answer string) error {
	user, err := Login(username, password)
	if err != nil {
		return err
	}

	wrongs, err := checkAnswer(answer)
	if err != nil {
		return err // evil client
	}

	if len(wrongs) > 0 {
		qids := []string{}
		for _, qid := range wrongs {
			// question numbers start from 1
			qids = append(qids, strconv.Itoa(qid+1))
		}
		return fmt.Errorf("第%s题答错", strings.Join(qids, ","))
	}

	_, err = db.Model(user).Set("activated = TRUE").Update()
	if err != nil {
		log.Fatal("db.Activate", err)
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

func GetCPoints() []model.CPointEntry {
	var res []model.CPointEntry

	err := db.Model(&res).
		Where("c_point > 0").
		Order("c_point DESC").
		Select()
	if err != nil {
		log.Fatalln("mako.GetCPoints", err)
	}

	return res
}

func UpdateCPoint(username string, delta int) error {
	res, err := db.Model(&model.CPointEntry{}).
		Set("c_point=c_point+?", delta).
		Where("username=?", username).
		Update()

	if err == nil {
		aff := res.RowsAffected()
		if aff != 1 {
			err = fmt.Errorf("%d row(s) affected", aff)
		}
	}

	return err
}

func ClaimFood(uid model.Uid, gotAt *time.Time) error {
	res, err := db.Model(&model.User{}).
		Set("food=food+(50*c_point)").
		Set("got_food_at=?", gotAt).
		Where("user_id=?", uid).
		Update()

	if err == nil {
		aff := res.RowsAffected()
		if aff != 1 {
			err = fmt.Errorf("%d row(s) affected", aff)
		}
	}

	return err
}

func UpdateFood(uid model.Uid, delta int) error {
	res, err := db.Model(&model.User{}).
		Set("food=food+?", delta).
		Where("user_id=?", uid).
		Update()

	if err == nil {
		aff := res.RowsAffected()
		if aff != 1 {
			err = fmt.Errorf("%d row(s) affected", aff)
		}
	}

	return err
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
