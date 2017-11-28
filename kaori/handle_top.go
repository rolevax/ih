package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type handlerTop func(rl *readline.Instance, args []string)

var handlers = map[string]handlerTop{
	"login":           handleLogin,
	"logout":          handleLogout,
	"look-around":     handleLookAround,
	"get-replay-list": handleGetReplayList,
	"get-replay":      handleGetReplay,
	"room-create":     handleRoomCreate,
	"room-join":       handleRoomJoin,
	"room-quit":       handleRoomQuit,
	"match-join":      handleMatchJoin,
}

func init() {
	// prevent init loop
	handlers["help"] = handleHelp
}

func handleTop(rl *readline.Instance, line string) {
	args := strings.Split(line, " ")
	h, ok := handlers[args[0]]
	if ok {
		h(rl, args[1:])
	} else {
		fmt.Println("what?")
	}
}

func handleHelp(rl *readline.Instance, args []string) {
	fmt.Println("avaiable commands:")
	cmds := []string{}
	for cmd, _ := range handlers {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	for _, cmd := range cmds {
		fmt.Println("  " + cmd)
	}
	fmt.Println("vim mode and tab-completion is supported")
}

func handleLogin(rl *readline.Instance, args []string) {
	if !(len(args) == 1 || len(args) == 2) {
		fmt.Println("usage: login <username> [password]")
		return
	}

	username := args[0]
	password := ""

	if len(args) == 2 {
		password = args[1]
	} else {
		pw, err := rl.ReadPassword("password:")
		if err != nil {
			log.Fatalln(err)
		}
		password = string(pw)
	}

	err := login(rl, username, password)
	if err != nil {
		log.Fatalln(err)
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

func handleMatchJoin(rl *readline.Instance, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: match-join <rule-id>")
		return
	}

	ruleId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	matchJoin(ruleId)
}
