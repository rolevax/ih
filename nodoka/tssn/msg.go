package tssn

import (
	"github.com/mjpancake/ih/ako/cs"
	"github.com/mjpancake/ih/ako/model"
	"github.com/mjpancake/ih/nodoka"
)

type cpReg struct {
	add  bool
	tssn *tssn
}

type pcSeat struct {
	*nodoka.MtSeat
}

type pcAction struct {
	*nodoka.MtAction
}

type ccSeat struct {
	Uid model.Uid
}

type ccAction struct {
	UserIndex int
	Act       *cs.Action
}
