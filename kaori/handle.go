package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chzyer/readline"
)

type handler func(rl *readline.Instance, args []string)

var handlers = map[string]handler{
	"login":         handleLogin,
	"logout":        handleLogout,
	"getReplayList": handleGetReplayList,
	"getReplay":     handleGetReplay,
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
