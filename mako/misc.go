package mako

import (
	"log"

	"github.com/go-redis/redis"
)

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

func AddTaskWater(water string) {
	err := rclient.LPush("mako.task.water", water).Err()
	if err != nil {
		log.Fatalln("mako.AddTaskWater", err)
	}
}

func GetTaskWaters(ct int) []string {
	res, err := rclient.LRange("mako.task.water", 0, int64(ct)).Result()
	if err != nil {
		log.Fatalln("mako.GetTaskWaters", err)
	}
	return res
}
