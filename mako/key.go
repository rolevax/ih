package mako

import (
	"fmt"

	"github.com/rolevax/ih/ako/model"
)

const (
	keyCPoints = "mako.c.points"
	keyGenUid  = "mako.gen.uid"
)

func keyAuth(username string) string {
	return fmt.Sprintf("mako.auth:%v", username)
}

func keyUser(uid model.Uid) string {
	return fmt.Sprintf("mako.user:%v", uid)
}
