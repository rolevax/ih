package srv

import (
	"log"
	"gopkg.in/redis.v5"
)

// redis access object
type rao struct {
	client		*redis.Client
}

func newRao() *rao {
	rao := new(rao)

	rao.client = redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0})

	_, err := rao.client.Ping().Result()
	if err != nil {
		log.Fatalln("redis", err)
	}

	err = rao.client.SAdd("accepting.versions", Version).Err()
	if err != nil {
		log.Fatalln("redis", err)
	}

	return rao
}

func (rao *rao) AcceptVersion(ver string) bool {
	res, err := rao.client.SIsMember("accepting.versions", ver).Result()
	if err != nil {
		log.Fatalln("rao.AcptVer", err)
	}
	return res
}

