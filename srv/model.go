package srv

import (
	"net"
)

type login struct {
	Version		string
	Username	string
	Password	string
	conn		net.Conn
}

type uid uint

type user struct {
	Id			uid
	Username	string
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
	uid			uid
	Nonce		int
	ActStr		string
	ActArg		string
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

type respLookAround struct {
	Type		string
	Bookable	bool
	Conn		int
	Book		int
	Play		int
}

func newRespLookAround(bookable bool, conn, book, play int) interface{} {
	return respLookAround{"look-around", bookable, conn, book, play}
}



