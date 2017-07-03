package book

import (
	"log"

	"github.com/mjpancake/ih/ako/model"
)

type BookState struct {
	Waits [4]model.Uid
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

func (bs *BookState) fillByAi() {
	if bs.Wait == 4 {
		// do nothing
	} else if bs.Wait == 2 {
		bs.Waits[2] = bs.Waits[1]
		bs.Waits[1] = model.UidAi1
		bs.Waits[3] = model.UidAi2
		bs.Wait = 4
	} else {
		log.Fatalln("BookState.fillByAi: wrong wait ct", bs.Wait)
	}
}
