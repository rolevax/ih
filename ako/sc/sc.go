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
	GirlKeys   [4]model.GirlKey
	TempDealer int
}

// rotate perspective
func (msg *TableSeat) RightPers() *TableSeat {
	next := &TableSeat{}

	*next = *msg
	next.TempDealer = (msg.TempDealer + 3) % 4
	next.GirlKeys[0] = msg.GirlKeys[1]
	next.GirlKeys[1] = msg.GirlKeys[2]
	next.GirlKeys[2] = msg.GirlKeys[3]
	next.GirlKeys[3] = msg.GirlKeys[0]

	return next
}

type VarMap map[string]interface{}

type TableEvent struct {
	Event string
	Args  VarMap
}

type TableEnd struct {
	Abortive    bool
	FoodChanges []*model.FoodChange
}
