package msg

import "github.com/rolevax/ih/ako/model"

type CsAccountCreate struct {
	Username string
	Password string
}

type CsAccountActivate struct {
	Username string
	Password string
	Answers  string
}

type CsAdminCPoint struct {
	Token       string
	Username    string
	CPointDelta int
}

type Sc struct {
	Error string // no news is good news
}

type ScCpoints struct {
	Sc
	Entries []model.CPointEntry
}
