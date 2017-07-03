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

type pcChoose struct {
	*nodoka.MtChoose
}

type pcReady struct {
	*nodoka.MtReady
}

type pcAction struct {
	*nodoka.MtAction
}

type ccChoose struct {
	Uid  model.Uid
	Gidx int
}

type ccReady struct {
	Uid model.Uid
}

type ccAction struct {
	UserIndex int
	Act       *cs.Action
}
