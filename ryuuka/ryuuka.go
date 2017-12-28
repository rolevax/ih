package ryuuka

import (
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorlog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
	"github.com/rolevax/ih/ako/ss"
	"github.com/rolevax/ih/toki/api"
)

const (
	timeout = 5 * time.Second
)

var (
	tokiPid  *actor.PID
	chReconn chan struct{}
)

func Init(selfAddr, tokiAddr string) {
	remote.SetLogLevel(actorlog.OffLevel)
	remote.Start(selfAddr)

	chReconn = make(chan struct{})
	go func() {
		for _ = range chReconn {
			reconnectToki(tokiAddr)
		}
	}()
	chReconn <- struct{}{}
}

func SendToToki(msg proto.Message) (proto.Message, error) {
	if tokiPid == nil {
		chReconn <- struct{}{}
		return nil, fmt.Errorf("toki down")
	}

	res, err := tokiPid.RequestFuture(msg, timeout).Result()
	if err != nil {
		log.Println("send to toki:", err)
		tokiPid = nil
		return nil, err
	} else {
		return res.(proto.Message), nil
	}
}

func reconnectToki(tokiAddr string) {
	if tokiPid != nil {
		return
	}

	pid := actor.NewPID(tokiAddr, toki.ActorName)
	_, err := pid.RequestFuture(&ss.TablePing{}, 10*time.Second).Result()
	if err != nil {
		log.Println("connect toki:", err)
	} else {
		log.Println("connected to toki")
		tokiPid = pid
	}
}
