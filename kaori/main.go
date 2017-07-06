package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	prefix := time.Now().Format("15:04:05")
	return fmt.Print(prefix, " ", string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	l := startReadline()
	defer l.Close()

	loop(l)
}

func startReadline() *readline.Instance {
	l, err := readline.New("\033[31m»\033[0m ")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(l.Stderr())
	return l
}

func loop(rl *readline.Instance) {
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		handle(rl, line)
	}
}

func handle(rl *readline.Instance, line string) {
	args := strings.Split(line, " ")
	h, ok := handlers[args[0]]
	if ok {
		h(rl, args[1:])
	} else {
		fmt.Println("what?")
	}
}
