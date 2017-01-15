package model

import (
	"net"
)

type UserAuth struct {
	Username	string
	Password	string
	Conn		net.Conn
}

type User struct {
	Id			int
	Username	string
	Nickname	string
	Conn		net.Conn	`json:"-"`
}

