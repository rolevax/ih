package tssn

import (
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/node"
)

type cpReg struct {
	add  bool
	tssn *tssn
}

type pcChoose struct {
	*node.MtChoose
}

type pcReady struct {
	*node.MtReady
}

type pcAction struct {
	*node.MtAction
}

type ccChoose struct {
	Uid  model.Uid
	Gidx int
}

type ccReady struct {
	Uid model.Uid
}
