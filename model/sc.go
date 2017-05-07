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
	Books [4]BookEntry
}

func NewScLookAround(conn int) *ScLookAround {
	return &ScLookAround{Type: "look-around", Conn: conn}
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
