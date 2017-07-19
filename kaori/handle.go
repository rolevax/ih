package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chzyer/readline"
)

type handler func(rl *readline.Instance, args []string)

var handlers = map[string]handler{
	"login":           handleLogin,
	"logout":          handleLogout,
	"look-around":     handleLookAround,
	"get-replay-list": handleGetReplayList,
	"get-replay":      handleGetReplay,
	"room-create":     handleRoomCreate,
	"room-join":       handleRoomJoin,
	"room-quit":       handleRoomQuit,
}

func handleLogin(rl *readline.Instance, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: login <username>")
		return
	}

	username := args[0]

	pw, err := rl.ReadPassword("password:")
	password := string(pw)

	err = login(username, password)
	if err != nil {
		log.Fatal(err)
	}
}

func handleLogout(rl *readline.Instance, args []string) {
	logout()
}

func handleLookAround(rl *readline.Instance, args []string) {
	lookAround()
}

func handleGetReplayList(rl *readline.Instance, args []string) {
	getReplayList()
}

func handleGetReplay(rl *readline.Instance, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: getReplay <replay-id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	getReplay(uint(id))
}

func handleRoomCreate(rl *readline.Instance, args []string) {
	roomCreate()
}

func handleRoomJoin(rl *readline.Instance, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: room-join <rid>")
		return
	}

	rid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	roomJoin(rid)
}

func handleRoomQuit(rl *readline.Instance, args []string) {
	roomQuit()
}
