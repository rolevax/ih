package model

type ScTypeOnly struct {
	Type string
}

type ScAuthFail struct {
	Type   string
	Ok     bool
	Reason string
}

func NewScAuthFail(str string) *ScAuthFail {
	return &ScAuthFail{"auth", false, str}
}

type ScAuthOk struct {
	Type string
	Ok   bool
	User *User
}

func NewScAuthOk(u *User) *ScAuthOk {
	return &ScAuthOk{"auth", true, u}
}

type ScLookAround struct {
	Type  string
	Conn  int
	Water []string
	Books [BookTypeKinds]BookEntry
}

func NewScLookAround(conn int, water []string, dcbaBookable *[4]bool,
	waits *[BookTypeKinds]int, tables *[BookTypeKinds]int) *ScLookAround {
	msg := &ScLookAround{
		Type:  "look-around",
		Conn:  conn,
		Water: water,
	}

	for i := 0; i < BookTypeKinds; i++ {
		bt := BookType(i)
		msg.Books[i].Bookable = dcbaBookable[int(bt.Abcd())]
		msg.Books[i].Book = waits[i]
		msg.Books[i].Play = tables[i] * bt.NeedUser()
	}

	return msg
}

type ScUpdateUser struct {
	Type  string
	User  *User
	Stats []StatRow
}

func NewScUpdateUser(user *User, stats []StatRow) *ScUpdateUser {
	resp := &ScUpdateUser{
		Type:  "update-user",
		User:  user,
		Stats: stats,
	}
	return resp
}

type ScGetReplayList struct {
	Type      string
	ReplayIds []uint
}

func NewScGetReplayList(replayIds []uint) *ScGetReplayList {
	return &ScGetReplayList{
		Type:      "get-replay-list",
		ReplayIds: replayIds,
	}
}

type ScGetReplay struct {
	Type       string
	ReplayId   uint
	ReplayJson string
}

func NewScGetReplay(replayId uint, replayJson string) *ScGetReplay {
	return &ScGetReplay{
		Type:       "get-replay",
		ReplayId:   replayId,
		ReplayJson: replayJson,
	}
}

type ScStart struct {
	Type       string
	Users      [4]*User
	TempDealer int
	Choices    [12]Gid
}

func NewScStart(users [4]*User, td int, cs [12]Gid) *ScStart {
	return &ScStart{
		Type:       "start",
		Users:      users,
		TempDealer: td,
		Choices:    cs,
	}
}

func (msg *ScStart) RightPers() {
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

type ScChosen struct {
	Type    string
	GirlIds [4]Gid
}

func NewScChosen(gids [4]Gid) *ScChosen {
	return &ScChosen{
		Type:    "chosen",
		GirlIds: gids,
	}
}

func (msg *ScChosen) RightPers() {
	gs := &msg.GirlIds
	gs[0], gs[1], gs[2], gs[3] = gs[1], gs[2], gs[3], gs[0]
}

type ScTableEvent struct {
	Type  string
	Event string
	Args  map[string]interface{}
	Nonce int
}
