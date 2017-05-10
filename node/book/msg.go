package book

import "github.com/mjpancake/hisa/model"

type BookState struct {
	Waits [model.BookTypeKinds]model.Uid
	Wait  int
}

func (bs *BookState) removeIfAny(uid model.Uid) {
	i := 0
	for i < bs.Wait && bs.Waits[i] != uid {
		i++
	}
	if i == bs.Wait {
		return
	}
	// swap to back, and pop back
	e := bs.Wait - 1
	bs.Waits[i], bs.Waits[e] = bs.Waits[e], bs.Waits[i]
	bs.Wait--
}
