package srv

import (
)

// user id
type uid uint

// girl id, signed-int for compatibility to libsaki
type gid int

type user struct {
	Id			uid
	Username	string
	Level		int
	Pt			int
	Rating		float64
	Ranks		[4]int
}

type girl struct {
	Id			gid
	Level		int
	Pt			int
	Rating		float64
	Ranks		[4]int
}

type reqTypeOnly struct {
	Type		string
}

type reqAuth struct {
	Type		string
	Version		string
	Username	string
	Password	string
}

type reqAction struct {
	Nonce		int
	ActStr		string
	ActArg		string
}

type respTypeOnly struct {
	Type		string
}

type respAuthFail struct {
	Type		string
	Ok			bool
	Reason		string
}

func newRespAuthFail(str string) interface{} {
	return respAuthFail{"auth", false, str}
}

type respAuthOk struct {
	Type	string
	Ok		bool
	User	*user
}

func newRespAuthOk(u *user) interface{} {
	return respAuthOk{"auth", true, u}
}

type bookEntry struct {
	Bookable	bool
	Book		int
	Play		int
}

type respLookAround struct {
	Type		string
	Conn		int
	Books		map[string]bookEntry
}

func newRespLookAround(bookable bool, conn, book, play int) interface{} {
	resp := new(respLookAround)
	resp.Type = "look-around"
	resp.Conn = conn
	resp.Books = make(map[string]bookEntry)
	resp.Books["DS71"] = bookEntry{bookable, book, play}
	resp.Books["CS71"] = bookEntry{false,0,0}
	resp.Books["BS71"] = bookEntry{false,0,0}
	resp.Books["AS71"] = bookEntry{false,0,0}
	return resp
}

type respUpdateUser struct {
	Type		string
	User		*user
}

func newRespUpdateUser(user *user) *respUpdateUser {
	resp := new(respUpdateUser)
	resp.Type = "update-user"
	resp.User = user
	return resp
}



