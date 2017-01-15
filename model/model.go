package model

import (
	"net"
)

type Login struct {
	Username	string
	Password	string
	Conn		net.Conn
}

type Uid int

type User struct {
	Id			Uid
	Username	string
	Nickname	string
}

