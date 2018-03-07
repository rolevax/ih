package nodoka

import (
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
)

type MtHasUser struct {
	Uid model.Uid
}

type MtCtTables struct{}

type MtAction struct {
	Uid model.Uid
	Act *cs.TableAction
}

type MtChoose struct {
	Uid model.Uid
	*cs.TableChoose
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

type MbGetMatchWaits struct{}

type MbMatchJoin struct {
	User model.User
	cs.MatchJoin
}

type MbMatchCancel struct {
	Uid model.Uid
}

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
