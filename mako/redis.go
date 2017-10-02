package mako

import (
	"fmt"
	"log"

	"gopkg.in/redis.v5"
)

var rclient *redis.Client

func init() {
	rclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})

	_, err := rclient.Ping().Result()
	if err != nil {
		log.Fatalln("redis", err)
	}
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
	res, err := rclient.Get("mako.admin.token").Result()
	if err != nil {
		log.Fatal("mako.CheckAdminToken", err)
	}
	return token == res
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
