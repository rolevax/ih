package tssn

import (
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/nodoka"
)

type cpReg struct {
	add  bool
	tssn *tssn
}

type pcChoose struct {
	*nodoka.MtChoose
}

type pcSeat struct {
	*nodoka.MtSeat
}

type pcAction struct {
	*nodoka.MtAction
}

type ccChoose struct {
	Uid model.Uid
}

type ccSeat struct {
	Uid model.Uid
}

type ccAction struct {
	UserIndex int
	Act       *cs.TableAction
}
