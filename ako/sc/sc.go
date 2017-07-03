package sc

import "github.com/mjpancake/ih/ako/model"

type TypeOnly struct {
	Type string
}

type AuthFail struct {
	Type   string
	Ok     bool
	Reason string
}

func NewAuthFail(str string) *AuthFail {
	return &AuthFail{"auth", false, str}
}

type AuthOk struct {
	Type string
	Ok   bool
	User *model.User
}

func NewAuthOk(u *model.User) *AuthOk {
	return &AuthOk{"auth", true, u}
}

type LookAround struct {
	Type  string
	Conn  int
	Water []string
	Books [model.BookTypeKinds]model.BookEntry
}

func NewLookAround(conn int, water []string, dcbaBookable *[4]bool,
	waits *[model.BookTypeKinds]int,
	tables *[model.BookTypeKinds]int) *LookAround {
	msg := &LookAround{
		Type:  "look-around",
		Conn:  conn,
		Water: water,
	}

	for i := 0; i < model.BookTypeKinds; i++ {
		bt := model.BookType(i)
		msg.Books[i].Bookable = dcbaBookable[int(bt.Abcd())]
		msg.Books[i].Book = waits[i]
		msg.Books[i].Play = tables[i] * bt.NeedUser()
	}

	return msg
}

type UpdateUser struct {
	Type  string
	User  *model.User
	Stats []model.StatRow
}

func NewUpdateUser(user *model.User, stats []model.StatRow) *UpdateUser {
	resp := &UpdateUser{
		Type:  "update-user",
		User:  user,
		Stats: stats,
	}
	return resp
}

type GetReplayList struct {
	Type      string
	ReplayIds []uint
}

func NewGetReplayList(replayIds []uint) *GetReplayList {
	return &GetReplayList{
		Type:      "get-replay-list",
		ReplayIds: replayIds,
	}
}

type GetReplay struct {
	Type       string
	ReplayId   uint
	ReplayJson string
}

func NewGetReplay(replayId uint, replayJson string) *GetReplay {
	return &GetReplay{
		Type:       "get-replay",
		ReplayId:   replayId,
		ReplayJson: replayJson,
	}
}

type Start struct {
	Type       string
	Users      [4]*model.User
	TempDealer int
	Choices    [12]model.Gid
}

func NewStart(users [4]*model.User, td int, cs [12]model.Gid) *Start {
	return &Start{
		Type:       "start",
		Users:      users,
		TempDealer: td,
		Choices:    cs,
	}
}

func (msg *Start) RightPers() {
	msg.TempDealer = (msg.TempDealer + 3) % 4

	// rotate perspectives
	u0 := msg.Users[0]
	msg.Users[0] = msg.Users[1]
	msg.Users[1] = msg.Users[2]
	msg.Users[2] = msg.Users[3]
	msg.Users[3] = u0

	cs := &msg.Choices
	cpu := len(cs) / 4 // choice per user
	for i := 0; i < cpu; i++ {
		tmp := cs[i]
		for w := 0; w < 3; w++ {
			cs[w*cpu+i] = cs[(w+1)*cpu+i]
		}
		cs[3*cpu+i] = tmp
	}
}

type Chosen struct {
	Type    string
	GirlIds [4]model.Gid
}

func NewChosen(gids [4]model.Gid) *Chosen {
	return &Chosen{
		Type:    "chosen",
		GirlIds: gids,
	}
}

func (msg *Chosen) RightPers() {
	gs := &msg.GirlIds
	gs[0], gs[1], gs[2], gs[3] = gs[1], gs[2], gs[3], gs[0]
}

type TableEvent struct {
	Type  string
	Event string
	Args  map[string]interface{}
	Nonce int
}
