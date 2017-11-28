package ussn

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/sc"
	"github.com/rolevax/ih/hayari"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/nodoka"
)

const uuReqTmot = 1 * time.Second

func (ussn *ussn) handleReject(msg error) {
	jsonb := sc.ToJson(&sc.Auth{
		Error: msg.Error(),
	})
	hayari.Write(ussn.conn, jsonb)
	ussn.handleError(fmt.Errorf("rejected: %v", msg))
}

func (ussn *ussn) handleError(msg error) {
	if ussn.user != nil {
		log.Println(ussn.user.Id, "----", msg)
	} else {
		log.Println(ussn.conn.RemoteAddr(), "----", msg)
	}
	ussn.p.Stop()
}

func (ussn *ussn) handleSc(msg interface{}, resp func(interface{})) {
	jsonb := sc.ToJson(msg)
	err := hayari.Write(ussn.conn, jsonb)
	if err != nil {
		if e, ok := err.(*net.OpError); ok {
			err = e.Err
		}
		ussn.handleError(err)
	} else {
		//log.Println(ussn.user.Id, "<---", string(jsonb))
	}

	resp(err)
}

func (ussn *ussn) handleUpdateInfo() {
	ussn.user = mako.GetUser(ussn.user.Id)
	sc := &sc.UpdateUser{
		User:  ussn.user,
		Stats: mako.GetCultis(ussn.user.Id),
	}
	ussn.handleSc(sc, noResp)
}

func (ussn *ussn) handleRead(breq []byte) {
	msg, err := cs.FromJson(breq)
	if err != nil {
		log.Println(ussn.user.Id, "-X->", string(breq))
		ussn.handleError(err)
	} else {
		ussn.handleCs(msg)
	}
}

func (ussn *ussn) handleCs(i interface{}) {
	switch msg := i.(type) {
	case *cs.LookAround:
		ussn.handleLookAround()
	case *cs.Heartbeat:
		// do nothing
	case *cs.RoomCreate:
		ussn.handleRoomCreate(msg)
	case *cs.RoomJoin:
		ussn.handleRoomJoin(msg)
	case *cs.RoomQuit:
		ussn.handleRoomQuit()
	case *cs.MatchJoin:
		ussn.handleMatchJoin(msg)
	case *cs.GetReplayList:
		ussn.handleGetReplayList()
	case *cs.GetReplay:
		ussn.handleGetReplay(msg.ReplayId)
	case *cs.TableChoose:
		ussn.handleChoose(msg)
	case *cs.TableSeat:
		ussn.handleSeat()
	case *cs.TableAction:
		ussn.handleAction(msg)
	default:
		ussn.handleError(fmt.Errorf("unexpected CsMsg %T\n", msg))
	}
}

func (ussn *ussn) handleRoomCreate(msg *cs.RoomCreate) {
	// TODO check gid pcs-ed
	nodoka.Bmgr.Tell(&nodoka.MbRoomCreate{
		Creator:    *ussn.user,
		RoomCreate: *msg,
	})
}

func (ussn *ussn) handleRoomJoin(msg *cs.RoomJoin) {
	nodoka.Bmgr.Tell(&nodoka.MbRoomJoin{
		User:     *ussn.user,
		RoomJoin: *msg,
	})
}

func (ussn *ussn) handleRoomQuit() {
	nodoka.Bmgr.Tell(&nodoka.MbRoomQuit{Uid: ussn.user.Id})
}

func (ussn *ussn) handleMatchJoin(msg *cs.MatchJoin) {
	nodoka.Bmgr.Tell(&nodoka.MbMatchJoin{
		User:      *ussn.user,
		MatchJoin: *msg,
	})
}

func (ussn *ussn) handleLookAround() {
	// TODO
	// all user have the same result,
	// no need to compute individually
	// Umgr periodically gather data from other Mgr's
	// ussn tell Umgr cpLookAround, Umgr tell ussn pcSc
	res, err := nodoka.Umgr.RequestFuture(&cpWater{}, uuReqTmot).Result()
	if err != nil {
		ussn.handleError(err)
		return
	}
	water := res.(*pcWater)

	playCt, err := (&nodoka.MtCtPlays{}).Req()
	if err != nil {
		ussn.handleError(err)
		return
	}

	/* feature hidden
	rooms, err := (&nodoka.MbGetRooms{}).Req()
	if err != nil {
		ussn.handleError(err)
		return
	}
	*/

	waits, err := (&nodoka.MbGetMatchWaits{}).Req()
	if err != nil {
		ussn.handleError(err)
		return
	}

	msg := &sc.LookAround{
		Conn:  water.ct,
		Play:  playCt,
		Water: water.water,
		//Rooms: rooms,
		MatchWaits: waits,
	}
	ussn.handleSc(msg, noResp)
}

func (ussn *ussn) handleGetReplayList() {
	ids := mako.GetReplayList(ussn.user.Id)
	ussn.handleSc(&sc.GetReplayList{ids}, noResp)
}

func (ussn *ussn) handleGetReplay(replayId uint) {
	text, err := mako.GetReplay(replayId)
	if err != nil {
		ussn.handleError(err)
		return
	}
	time.Sleep(2 * time.Second) // no reason, just wanna sleep
	ussn.handleSc(&sc.GetReplay{replayId, text}, noResp)
}

func (ussn *ussn) handleChoose(msg *cs.TableChoose) {
	nodoka.Tmgr.Tell(&nodoka.MtChoose{Uid: ussn.user.Id, TableChoose: msg})
}

func (ussn *ussn) handleSeat() {
	nodoka.Tmgr.Tell(&nodoka.MtSeat{Uid: ussn.user.Id})
}

func (ussn *ussn) handleAction(msg *cs.TableAction) {
	nodoka.Tmgr.Tell(&nodoka.MtAction{Uid: ussn.user.Id, Act: msg})
}
