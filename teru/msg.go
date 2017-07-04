package main

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

type ScAccountActivate struct {
	Sc
	Result string
}
