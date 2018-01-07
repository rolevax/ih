package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorlog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/rolevax/ih/ako/ss"
	"github.com/rolevax/ih/toki/api"
)

func Receive(ctx actor.Context) {
	log.Println(ctx.Message())

	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *ss.TablePing:
		handlePing(msg, ctx.Respond)
	case *ss.TableStart:
		handleStart(msg, ctx.Respond)
	case *ss.TableAction:
		handleAction(msg, ctx.Respond)
	case *ss.TableSweepOne:
		handleSweepOne(msg, ctx.Respond)
	case *ss.TableSweepAll:
		handleSweepAll(msg, ctx.Respond)
	case *ss.TableDeleteIfAny:
		handleDeleteIfAny(msg, ctx.Respond)
	default:
		log.Fatalln("unknown message type %T\n", msg)
	}
}

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

	addr := flag.String("addr", "localhost:8900", "addr to listen")
	flag.Parse()

	remote.SetLogLevel(actorlog.OffLevel)
	remote.Start(*addr)
	props := actor.FromFunc(Receive)
	actor.SpawnNamed(props, toki.ActorName)
	log.Println("started", toki.ActorName, "at", *addr)

	<-make(chan struct{}) // block
}
