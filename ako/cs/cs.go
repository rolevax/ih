package cs

import "github.com/rolevax/ih/ako/model"

type TypeOnly struct {
	Type string
}

type Auth struct {
	Version  string
	Username string
	Password string
}

type LookAround struct{}

type Heartbeat struct{}

type RoomCreate struct {
	GirlId model.Gid
	AiNum  model.AiNum
	Bans   []model.Gid
	AiGids []model.Gid
}

type RoomJoin struct {
	GirlId model.Gid
	RoomId model.Rid
}

type RoomQuit struct{}

type Seat struct{}

type Action struct {
	Nonce   int
	ActStr  string
	ActArg  int
	ActTile string
}

type GetReplay struct {
	ReplayId uint
}

type GetReplayList struct{}
