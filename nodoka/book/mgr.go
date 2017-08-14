package book

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
	"github.com/rolevax/ih/nodoka"
	"github.com/rolevax/ih/nodoka/tssn"
)

func Init() {
	props := actor.FromFunc(Receive)
	pid, err := actor.SpawnNamed(props, "Bmgr")
	if err != nil {
		log.Fatalln(err)
	}
	nodoka.Bmgr = pid
}

func Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *nodoka.MbRoomCreate:
		handleRoomCreate(msg)
	case *nodoka.MbRoomJoin:
		handleRoomJoin(msg)
	case *nodoka.MbRoomQuit:
		handleRoomQuit(msg)
	case *nodoka.MbGetRooms:
		resp := handleGetRooms()
		ctx.Respond(resp)
	default:
		log.Fatalf("Bmgr.Recv unexpected %T\n", msg)
	}
}

func handleRoomCreate(msg *nodoka.MbRoomCreate) {
	uid := msg.Creator.Id

	playing, err := (&nodoka.MtHasUser{Uid: uid}).Req()
	if err != nil {
		log.Println("Bmgr.handleRoomCreate:", err)
		nodoka.Umgr.Tell(&nodoka.MuKick{uid, err.Error()})
		return
	}
	if playing {
		nodoka.Umgr.Tell(&nodoka.MuKick{uid, "create but playing"})
		return
	}

	room, err := roomCreate(msg)
	if err != nil {
		nodoka.Umgr.Tell(&nodoka.MuKick{uid, err.Error()})
	}
	if room != nil {
		tssn.Start(room)
	}
}

func handleRoomJoin(msg *nodoka.MbRoomJoin) {
	room, err := roomJoin(msg)
	if err != nil {
		if err == errRoomTan90 {
			nodoka.Umgr.Tell(&nodoka.MuSc{
				To:  msg.User.Id,
				Msg: &sc.RoomJoin{Error: "来晚一步，房间已开"},
			})
		} else {
			nodoka.Umgr.Tell(&nodoka.MuKick{msg.User.Id, err.Error()})
		}
		return
	}

	nodoka.Umgr.Tell(&nodoka.MuSc{
		To:  msg.User.Id,
		Msg: &sc.RoomJoin{Error: ""},
	})

	if room != nil {
		tssn.Start(room)
	}
}

func handleRoomQuit(msg *nodoka.MbRoomQuit) {
	roomQuit(msg.Uid)
}

func handleGetRooms() []*model.Room {
	res := []*model.Room{}

	for _, r := range rooms {
		v := *r // copy
		res = append(res, &v)
	}

	return res
}
