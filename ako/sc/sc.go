package sc

import "github.com/rolevax/ih/ako/model"

type Auth struct {
	Error  string
	Resume bool
	User   *model.User
}

type LookAround struct {
	Conn  int
	Play  int
	Water []string
	Rooms []*model.Room
}

type UpdateUser struct {
	User  *model.User
	Stats []model.Culti
}

type RoomJoin struct {
	Error string
}

type GetReplayList struct {
	ReplayIds []uint
}

type GetReplay struct {
	ReplayId   uint
	ReplayJson string
}

type Seat struct {
	Room       model.Room
	TempDealer int
}

// rotate perspective
func (msg *Seat) RightPers() *Seat {
	next := &Seat{}

	*next = *msg
	next.TempDealer = (msg.TempDealer + 3) % 4
	next.Room.Users = append(msg.Room.Users[1:4], msg.Room.Users[0])
	next.Room.Gids = append(msg.Room.Gids[1:4], msg.Room.Gids[0])

	return next
}

type TableEvent struct {
	Event string
	Args  map[string]interface{}
	Nonce int
}
