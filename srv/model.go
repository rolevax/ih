package srv

import (
	"net"
)

type login struct {
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
	Username	string
	Password	string
}

type reqAction struct {
	uid		uid
	ActStr	string
	ActArg	string
}

type respAuthFail struct {
	Type	string
	Ok		bool
	Reason	string
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
	Type	string
	Conn	int
	Idle	int
	Book	int
	Play	int
}

func newRespLookAround(connCt, idleCt, bookCt, playCt int) interface{} {
	return respLookAround{"look-around", connCt, idleCt, bookCt, playCt}
}



