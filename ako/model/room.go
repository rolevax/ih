package model

import (
	"log"
)

type AiNum int
type Rid int

const (
	Ai0 = AiNum(0)
	Ai2 = AiNum(1)
	Ai3 = AiNum(2)
)

func (a AiNum) Valid() bool {
	i := int(a)
	return 0 <= i && i <= 2
}

func (a AiNum) String() string {
	switch a {
	case Ai0:
		return "Ai0"
	case Ai2:
		return "Ai2"
	case Ai3:
		return "Ai3"
	default:
		log.Fatalln("AiNum.String", a)
		return ""
	}
}

func (a AiNum) NeedUser() int {
	switch a {
	case Ai0:
		return 4
	case Ai2:
		return 2
	case Ai3:
		return 1
	default:
		log.Fatalln("AiNum.NeedUser")
		return -1
	}
}

func (a AiNum) NeedAi() int {
	return 4 - a.NeedUser()
}

type Room struct {
	Id    Rid
	AiNum AiNum
	Users []User
	Gids  []Gid
	Bans  []Gid
}

func (r *Room) Four() bool {
	return len(r.Users) == 4 && len(r.Gids) == 4
}

func (r *Room) FillAi(gids []Gid) {
	uids := []Uid{501, 502, 503}
	for i, gid := range gids {
		r.fillAi(uids[i], gid)
	}
}

func (r *Room) fillAi(uid Uid, gid Gid) {
	r.Users = append(r.Users, User{
		Id:       uid,
		Username: "ⓝ喵打",
	})
	r.Gids = append(r.Gids, gid)
}
