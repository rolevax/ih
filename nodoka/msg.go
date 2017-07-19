package nodoka

import (
	"github.com/mjpancake/ih/ako/cs"
	"github.com/mjpancake/ih/ako/model"
)

type MtHasUser struct {
	Uid model.Uid
}

type MtCtPlays struct{}

type MtAction struct {
	Uid model.Uid
	Act *cs.Action
}

type MtSeat struct {
	Uid model.Uid
}

type MbRoomCreate struct {
	Creator model.User
	cs.RoomCreate
}

type MbRoomJoin struct {
	User model.User
	cs.RoomJoin
}

type MbRoomQuit struct {
	Uid model.Uid
}

type MbGetRooms struct{}

type MuSc struct {
	To  model.Uid
	Msg interface{}
}

type MuKick struct {
	Uid    model.Uid
	Reason string
}

type MuUpdateInfo struct {
	Uid model.Uid
}
