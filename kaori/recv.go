package main

import (
	"log"

	"github.com/mjpancake/ih/ako/cs"
	"github.com/mjpancake/ih/ako/model"
	"github.com/mjpancake/ih/ako/sc"
)

func onRecv(b []byte) {
	s := string(b)
	msg := sc.FromJson(b)
	switch msg := msg.(type) {
	case *sc.Auth:
		if msg.Error == "" {
			log.Println("logged in")
		} else {
			log.Println(msg.Error)
		}
	case *sc.UpdateUser:
		log.Println("synced user data")
	case *sc.LookAround:
		recvLookAround(msg)
	case *sc.Seat:
		recvSeat(msg)
	default:
		log.Println(s)
	}
}

func recvLookAround(sc *sc.LookAround) {
	log.Println("conn", sc.Conn, "play", sc.Play)
	for _, room := range sc.Rooms {
		uids := []model.Uid{}
		for _, u := range room.Users {
			uids = append(uids, u.Id)
		}
		log.Println(room.Id, room.AiNum, uids, room.Gids)
	}
}

func recvSeat(sc *sc.Seat) {
	cl.send(&cs.Seat{})
	log.Println("seat")
}
