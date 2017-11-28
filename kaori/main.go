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

type handle func(rl *readline.Instance, line string)

var currHandle handle

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	l := startReadline()
	defer l.Close()

	currHandle = handleTop
	loop(l)
}

func startReadline() *readline.Instance {
	children := []readline.PrefixCompleterInterface{}
	for key, _ := range handlers {
		children = append(children, readline.PcItem(key))
	}
	comp := readline.NewPrefixCompleter()
	comp.SetChildren(children)

	l, err := readline.NewEx(&readline.Config{
		Prompt:       "\033[31mkaori»\033[0m ",
		AutoComplete: comp,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(l.Stderr())
	l.SetVimMode(true)

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

		currHandle(rl, line)
	}
}
