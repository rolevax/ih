package ussn

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/model"
)

func botSc(to model.Uid, msg interface{}, sender *actor.PID) {
	switch msg.(type) {
	default:
		log.Printf("ussn.botSc unhandled %T\n", msg)
	}
}
