package mako

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/hitomi"
)

func Login(username, password string) (*model.User, error) {
	auth, err := GetAuth(username)
	if err != nil {
		return nil, err
	}

	if hash(password) != auth.Password {
		return nil, errors.New("密码错误")
	}

	return GetUser(auth.Uid)
}

func SignUp(username, password string) error {
	if !hitomi.CheckName(username) {
		return errors.New("用户名不可用")
	}

	uid, err := rclient.Incr(keyGenUid).Result()
	if err != nil {
		return err
	}

	auth := &model.Auth{
		Uid:      model.Uid(uid),
		Password: hash(password),
	}
	bytes, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	ok, err := rclient.SetNX(keyAuth(username), bytes, 0).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("用户名已存在")
	}

	return nil
}

func GetAuth(username string) (*model.Auth, error) {
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

	return auth, nil
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

func GetCPoints() ([]model.CPointEntry, error) {
	zs, err := rclient.ZRangeWithScores(keyCPoints, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var res []model.CPointEntry

	for _, z := range zs {
		i, err := strconv.Atoi(z.Member.(string))
		if err != nil {
			return nil, err
		}

		user, err := GetUser(model.Uid(i))
		if err != nil {
			return nil, err
		}

		res = append(res, model.CPointEntry{
			Username: user.Username,
			CPoint:   int(z.Score),
		})
	}

	return res, nil
}

func UpdateCPoint(username string, delta int) error {
	auth, err := GetAuth(username)
	if err != nil {
		return err
	}

	return rclient.ZIncrBy(
		keyCPoints,
		float64(delta),
		auth.Uid.ToString(),
	).Err()
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
