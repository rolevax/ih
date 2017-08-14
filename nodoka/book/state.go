package book

import (
	"fmt"
	"log"

	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/nodoka"
)

var (
	idGen          = 0
	rooms          = map[model.Rid]*model.Room{}
	ridOfUid       = map[model.Uid]model.Rid{}
	errRoomTan90   = fmt.Errorf("room tan90")
	errRoomBan     = fmt.Errorf("ban wanton")
	errRoomDupUser = fmt.Errorf("already in same room")
	errRoomDupGirl = fmt.Errorf("room dup girl") // TODO use it
)

func roomCreate(msg *nodoka.MbRoomCreate) (*model.Room, error) {
	if rid, ok := ridOfUid[msg.Creator.Id]; ok {
		return nil, fmt.Errorf("double room create, rid %d", rid)
	}

	rid := genId()
	if _, ok := rooms[rid]; ok {
		log.Fatalf("rid conflict on %d", rid)
	}

	room := &model.Room{
		Id:    rid,
		AiNum: msg.AiNum,
		Users: []model.User{msg.Creator},
		Gids:  []model.Gid{msg.GirlId},
		Bans:  msg.Bans,
	}
	room.FillAi(msg.AiGids)

	rooms[rid] = room
	ridOfUid[msg.Creator.Id] = rid

	return roomPop(room)
}

func roomJoin(msg *nodoka.MbRoomJoin) (*model.Room, error) {
	rid := msg.RoomId

	room, ok := rooms[rid]
	if !ok {
		return nil, errRoomTan90
	}

	for _, gid := range room.Bans {
		if gid == msg.GirlId {
			return nil, errRoomBan
		}
	}

	for _, user := range room.Users {
		if user.Id == msg.User.Id {
			return nil, errRoomDupUser
		}
	}

	room.Users = append(room.Users, msg.User)
	room.Gids = append(room.Gids, msg.GirlId)

	return roomPop(room)
}

func roomPop(room *model.Room) (*model.Room, error) {
	if room.Four() {
		for _, u := range room.Users {
			delete(ridOfUid, u.Id)
		}
		delete(rooms, room.Id)
		return room, nil
	} else {
		return nil, nil
	}
}

func roomQuit(uid model.Uid) {
	rid, ok := ridOfUid[uid]
	if !ok {
		return
	}

	delete(ridOfUid, uid)

	room, ok := rooms[rid]
	if !ok {
		log.Fatalf("roomQuit: room %d tan90", rid)
	}

	for i, u := range room.Users {
		if u.Id == uid {
			room.Users = append(room.Users[:i], room.Users[i+1:]...)
			room.Gids = append(room.Gids[:i], room.Gids[i+1:]...)
			break
		}
	}

	if len(room.Users) == 0 {
		delete(rooms, rid)
	}
}

func genId() model.Rid {
	rid := model.Rid(idGen)
	idGen++
	idGen %= 1000 // assume 0~999 is enough
	return rid
}

/*
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
*/
