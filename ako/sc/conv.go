package sc

import (
	"log"

	"github.com/rolevax/ih/ako"
)

var decoder = ako.NewDecoder([]interface{}{
	Auth{},
	UpdateUser{},
	LookAround{},
	RoomJoin{},
	TableInit{},
	TableSeat{},
	TableEvent{},
	GetReplayList{},
	GetReplay{},
})

func FromJson(breq []byte) interface{} {
	sc, err := decoder.FromJson(breq)
	if err != nil {
		log.Fatal("sc.FromJson: ", err)
	}
	return sc
}

func ToJson(sc interface{}) []byte {
	return ako.ToJson(sc)
}
