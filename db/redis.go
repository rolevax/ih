package db

import (
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
	err := rclient.SAdd("accepting.versions", ver).Err()
	if err != nil {
		log.Fatalln("redis", err)
	}
}

func AcceptVersion(ver string) bool {
	res, err := rclient.SIsMember("accepting.versions", ver).Result()
	if err != nil {
		log.Fatalln("db.AcptVer", err)
	}
	return res
}
