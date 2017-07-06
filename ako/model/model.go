package model

import (
	"log"
	"strconv"
)

// user id
type Uid uint

const (
	UidSys Uid = 500
	UidAi1 Uid = 501
	UidAi2 Uid = 502
)

func (uid Uid) IsBot() bool {
	return uint(uid) < 1000
}

func (uid Uid) IsHuman() bool {
	return !uid.IsBot()
}

// girl id, signed-int for compatibility to libsaki
type Gid int

// level, pt, and rating
type Lpr struct {
	Level  int
	Pt     int
	Rating float64
}

type User struct {
	Id       Uid `sql:"user_id,pk"`
	Username string
	Lpr
}

type Girl struct {
	Id Gid
	Lpr
}

type BookType int

type Abcd int
type Rule int

const (
	BookTypeKinds = 8
	BookD         = Abcd(0)
	BookC         = Abcd(1)
	BookB         = Abcd(2)
	BookA         = Abcd(3)
	Rule4p        = Rule(0)
	Rule2p        = Rule(1)
)

func (b BookType) Index() int {
	return int(b)
}

func (b BookType) Valid() bool {
	i := int(b)
	return 0 <= i && i < BookTypeKinds
}

func (b BookType) Abcd() Abcd {
	return Abcd(int(b) % 4)
}

func (b BookType) Rule() Rule {
	return Rule(int(b) / 4)
}

func (b BookType) NeedUser() int {
	switch b.Rule() {
	case Rule4p:
		return 4
	case Rule2p:
		return 2
	default:
		log.Fatalln("BookType.NeedUser")
		return -1
	}
}

func (b BookType) String() string {
	dcba := [4]string{"D", "C", "B", "A"}
	return dcba[int(b.Abcd())] + strconv.Itoa(b.NeedUser()) + "P"
}

type BookEntry struct {
	Bookable bool
	Book     int
	Play     int
}
