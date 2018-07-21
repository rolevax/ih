package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/hitomi"
	"github.com/rolevax/ih/mako"
)

const (
	DictPath = "/srv/dict.txt"
)

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	prefix := time.Now().Format("01/02 15:04:05")
	return fmt.Print(prefix, " ", string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	if flag.Parsed() {
		log.Fatalln("unexpected flag parse before main()")
	}

	redis := flag.String("redis", "localhost:6379", "redis server addr")
	db := flag.String("db", "localhost:5432", "pg db server addr")
	dictPath := flag.String("dict", DictPath, "path to sensitive word dictionary")

	flag.Parse()
	mako.InitRedis(*redis)
	mako.InitDb(*db)
	hitomi.Init(*dictPath)

	migrate()
	log.Println("End of migration")
}

type MigUser struct {
	tableName struct{}  `sql:"users"`
	Id        model.Uid `sql:"user_id,pk"`
	Username  string
	Password  string
	CPoint    int
}

type RAuth struct {
	Uid      model.Uid
	Password string
}

func migrate() {
	db := mako.GetDb()
	r := mako.GetRClient()

	_ = r
	var users []MigUser

	err := db.Model(&users).Select()
	if err != nil {
		log.Fatalln("get users failed", err)
	}

	maxId := 0
	for _, user := range users {
		if int(user.Id) > maxId {
			maxId = int(user.Id)
		}

		auth := &RAuth{
			Uid:      user.Id,
			Password: user.Password,
		}
		bytes, err := json.Marshal(auth)
		if err != nil {
			log.Fatalln(err)
		}
		err = r.Set("mako.auth:"+user.Username, bytes, 0).Err()
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err = json.Marshal(user)
		if err != nil {
			log.Fatalln(err)
		}
		err = r.Set(fmt.Sprintf("mako.user:%v", user.Id), bytes, 0).Err()
		if err != nil {
			log.Fatalln(err)
		}

		if user.CPoint != 0 {
			err = r.ZAdd(
				"mako.c.points",
				redis.Z{float64(user.CPoint), int(user.Id)},
			).Err()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	log.Println("max id", maxId)
	err = r.Set("mako.gen.uid", maxId, 0).Err()
	if err != nil {
		log.Fatalln(err)
	}

	tasks, err := mako.GetTasks()
	if err != nil {
		log.Fatalln(err)
	}

	for _, task := range tasks {
		bytes, err := json.Marshal(task)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(bytes))
		err = r.Set(fmt.Sprintf("mako.task:%v", task.Id), bytes, 0).Err()
		if err != nil {
			log.Fatalln(err)
		}
		if task.State != model.TaskStateClosed {
			err = r.SAdd("mako.open.tasks", task.Id).Err()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
