package node

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/model"
)

const reqTimeOut = 500 * time.Millisecond

func Req(mgr *actor.PID, msg interface{}) (interface{}, error) {
	f := mgr.RequestFuture(msg, reqTimeOut)
	return f.Result()
}

func (msg *MtHasUser) Req() (bool, error) {
	res, err := Req(Tmgr, msg)
	if err != nil {
		return false, err
	} else {
		return res.(bool), nil
	}
}

func (msg *MtCtPlays) Req() ([model.BookTypeKinds]int, error) {
	res, err := Req(Tmgr, msg)
	if err != nil {
		return [model.BookTypeKinds]int{}, err
	} else {
		return res.([model.BookTypeKinds]int), nil
	}
}

func (msg *MbCtBooks) Req() ([model.BookTypeKinds]int, error) {
	res, err := Req(Bmgr, msg)
	if err != nil {
		return [model.BookTypeKinds]int{}, err
	} else {
		return res.([model.BookTypeKinds]int), nil
	}
}

func (msg *MuSc) Req() error {
	res, err := Req(Umgr, msg)
	if err != nil {
		return err
	} else if res != nil {
		return res.(error)
	} else {
		return nil
	}
}
