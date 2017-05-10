package model

type CsTypeOnly struct {
	Type string
}

type CsAuth struct {
	Type     string
	Version  string
	Username string
	Password string
}

type CsLookAround struct{}

type CsHeartBeat struct{}

type CsChoose struct {
	GirlIndex int
}

type CsReady struct{}

type CsAction struct {
	Nonce  int
	ActStr string
	ActArg string
}

type CsBook struct {
	BookType BookType
}

type CsUnbook struct{}

type CsGetReplay struct {
	ReplayId uint
}

type CsGetReplayList struct{}
