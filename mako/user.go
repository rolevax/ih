package mako

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/hitomi"
)

func Login(username, password string) (*model.User, error) {
	str, err := rclient.Get(keyAuth(username)).Result()
	if err == redis.Nil {
		return nil, errors.New("用户，不存在的x")
	} else if err != nil {
		return nil, err
	}

	auth := &model.Auth{}
	err = json.Unmarshal([]byte(str), auth)
	if err != nil {
		return nil, err
	}

	return GetUser(auth.Uid)
}

func SignUp(username, password string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if !hitomi.CheckName(username) {
		return errors.New("用户名不可用")
	}

	var exist bool
	_, err = tx.QueryOne(
		&exist,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=?)",
		username,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if exist {
		tx.Rollback()
		return errors.New("用户名已存在")
	}

	// using raw query since password is absent from the model
	var uid model.Uid
	_, err = tx.QueryOne(
		&uid,
		"INSERT INTO users(username, password) VALUES (?,?) RETURNING user_id",
		username, hash(password),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GetUser(uid model.Uid) (*model.User, error) {
	str, err := rclient.Get(keyUser(uid)).Result()
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	err = json.Unmarshal([]byte(str), user)
	if err != nil {
		return nil, err
	}

	cPoint, _ := rclient.ZScore(keyCPoints, uid.ToString()).Result()
	user.CPoint = int(cPoint)

	return user, nil
}

func GetUsers(uids *[4]model.Uid) [4]*model.User {
	var users [4]*model.User

	for i, uid := range uids {
		user, err := GetUser(uid)
		if err != nil {
			log.Fatalln(err)
		}
		users[i] = user
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

func hash(password string) string {
	sha := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(sha[:])
}
