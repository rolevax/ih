package main

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

type Sc struct {
	Error string // no news is good news
}

type ScCpoints struct {
	Sc
	Entries []model.CpointEntry
}
