package mako

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

const (
	DefaultAdminToken = "111111"
)

var rclient *redis.Client

func GetRClient() *redis.Client {
	return rclient
}

func InitRedis(addr string) {
	rclient = redis.NewClient(&redis.Options{
		Addr:     addr,
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
