package main

import (
	"fmt"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/rolevax/ih/ako/cs"
)

func handleTable(rl *readline.Instance, line string) {
	handleTableChoose(line)
}

func handleTableChoose(line string) error {
	gidx, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	if !(0 <= gidx && gidx < 3) {
		return fmt.Errorf("gidx out of range")
	}

	cl.send(&cs.TableChoose{Gidx: gidx})
	return nil
}
