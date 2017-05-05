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

type CsAction struct {
	Type   string
	Nonce  int
	ActStr string
	ActArg string
}

type CsBook struct {
	Type     string
	BookType BookType
}

type CsChoose struct {
	GirlIndex int
}
