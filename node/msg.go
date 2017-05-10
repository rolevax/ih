package node

import "github.com/mjpancake/hisa/model"

type MtHasUser struct {
	Uid model.Uid
}

type MtCtPlays struct{}

type MtAction struct {
	Uid model.Uid
	Act *model.CsAction
}

type MtChoose struct {
	Uid  model.Uid
	Gidx int
}

type MtReady struct {
	Uid model.Uid
}

type MbBook struct {
	Uid      model.Uid
	BookType model.BookType
}

type MbUnbook struct {
	Uid model.Uid
}

type MbCtBooks struct{}

type MuSc struct {
	To  model.Uid
	Msg interface{}
}

type MuKick struct {
	Uid    model.Uid
	Reason string
}

type MuUpdateInfo struct {
	Uid model.Uid
}
