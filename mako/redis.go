package mako

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/redis.v5"
)

const (
	DefaultAdminToken = "111111"
)

var rclient *redis.Client

func init() {
	rclient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // hostname in docker net
		Password: "",
		DB:       0,
	})

	for i := 0; i < 5; i++ {
		_, err := rclient.Ping().Result()
		if err == nil {
			log.Println("mako.init redis: Ok")
			if CheckAdminToken(DefaultAdminToken) {
				log.Println("mako.init WARNING: using default admin token")
			}
			return
		}

		log.Println("mako.init redis", err)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("mako.init redis: tried too many times")
}

func AddAcceptingVersion(ver string) {
	err := rclient.SAdd("mako.vers", ver).Err()
	if err != nil {
		log.Fatal("redis", err)
	}
}

func AcceptVersion(ver string) bool {
	res, err := rclient.SIsMember("mako.vers", ver).Result()
	if err != nil {
		log.Fatal("mako.AcptVer", err)
	}
	return res
}

func CheckAdminToken(token string) bool {
	return token == GetAdminToken()
}

func GetAdminToken() string {
	res, err := rclient.Get("mako.admin.token").Result()
	if err != nil {
		if err == redis.Nil {
			res = DefaultAdminToken
		} else {
			log.Fatal("mako.GetAdminToken", err)
		}
	}
	return res
}

func checkAnswer(answer string) ([]int, error) {
	corr, err := rclient.Get("mako.answer").Result()
	if err != nil {
		// should prepare 'mako.answer' manualy in redis
		log.Fatal("mako.Answer", err)
	}

	as := []byte(answer)
	cas := []byte(corr)
	if len(as) != len(cas) {
		str := "wrong answer len %d, want %d"
		return nil, fmt.Errorf(str, len(as), len(cas))
	}

	res := []int{}
	for i, _ := range as {
		if as[i] != cas[i] {
			res = append(res, i)
		}
	}

	return res, nil
}
