package nodoka

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
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

func (msg *MtCtPlays) Req() (int, error) {
	res, err := Req(Tmgr, msg)
	if err != nil {
		return 0, err
	} else {
		return res.(int), nil
	}
}

func (msg *MbGetRooms) Req() ([]*model.Room, error) {
	res, err := Req(Bmgr, msg)
	if err != nil {
		return nil, err
	} else {
		return res.([]*model.Room), nil
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
