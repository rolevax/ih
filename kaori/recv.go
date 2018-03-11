package main

import (
	"fmt"
	"log"

	"github.com/chzyer/readline"
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
)

func onRecv(b []byte, rl *readline.Instance) {
	s := string(b)
	msg := sc.FromJson(b)
	switch msg := msg.(type) {
	case *sc.Auth:
		recvAuth(msg, rl)
	case *sc.UpdateUser:
		log.Println("synced user data")
	case *sc.LookAround:
		recvLookAround(msg)
	case *sc.TableInit:
		recvTableInit(msg, rl)
	case *sc.TableSeat:
		recvTableSeat(msg)
	default:
		log.Println(s)
	}
}

func recvAuth(sc *sc.Auth, rl *readline.Instance) {
	if sc.Error == "" {
		rl.SetPrompt("\033[31m" + sc.User.Username + "»\033[0m ")
		log.Println("logged in")
	} else {
		log.Println(sc.Error)
	}
}

func recvLookAround(sc *sc.LookAround) {
	log.Println("connCt", sc.Conn, "tableCt", sc.Table)
	for _, room := range sc.Rooms {
		uids := []model.Uid{}
		for _, u := range room.Users {
			uids = append(uids, u.Id)
		}
		log.Println(room.Id, room.AiNum, uids, room.Gids)
	}
}

func recvTableInit(sc *sc.TableInit, rl *readline.Instance) {
	log.Println("table-init, rule is", sc.MatchResult.RuleId)
	for i, u := range sc.MatchResult.Users {
		log.Println("  ", i, ": ", u.Username)
	}
	log.Println("choices:", sc.Choices)
	fmt.Println("Please choose a girl by index")

	currHandle = handleTable
}

func recvTableSeat(sc *sc.TableSeat) {
	cl.send(&cs.TableSeat{})
	log.Println("seat")
}
