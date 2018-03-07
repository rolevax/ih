package sc

import "github.com/rolevax/ih/ako/model"

type Auth struct {
	Error  string
	Resume bool
	User   *model.User
}

type LookAround struct {
	Conn       int
	Table      int
	Water      []string
	Rooms      []*model.Room
	MatchWaits []int
}

type UpdateUser struct {
	User *model.User
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

type TableInit struct {
	MatchResult model.MatchResult
	Choices     [3]model.Gid
	FoodCosts   [3]int
}

func (msg *TableInit) RightPers() *TableInit {
	next := &TableInit{}

	*next = *msg
	next.MatchResult = *msg.MatchResult.RightPers()

	// choices are assigned, not rotated

	return next
}

type TableSeat struct {
	Gids       [4]model.Gid
	TempDealer int
}

// rotate perspective
func (msg *TableSeat) RightPers() *TableSeat {
	next := &TableSeat{}

	*next = *msg
	next.TempDealer = (msg.TempDealer + 3) % 4
	next.Gids[0] = msg.Gids[1]
	next.Gids[1] = msg.Gids[2]
	next.Gids[2] = msg.Gids[3]
	next.Gids[3] = msg.Gids[0]

	return next
}

type TableEvent struct {
	Event string
	Args  map[string]interface{}
}

type TableEnd struct {
	Abortive    bool
	FoodChanges []*model.FoodChange
}
