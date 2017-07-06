package cs

import "github.com/mjpancake/ih/ako/model"

type TypeOnly struct {
	Type string
}

type Auth struct {
	Version  string
	Username string
	Password string
}

type LookAround struct{}

type HeartBeat struct{}

type Choose struct {
	GirlIndex int
}

type Ready struct{}

type Action struct {
	Nonce  int
	ActStr string
	ActArg string
}

type Book struct {
	BookType model.BookType
}

type Unbook struct{}

type GetReplay struct {
	ReplayId uint
}

type GetReplayList struct{}
