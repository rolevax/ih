package ussn

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mjpancake/ih/ako/cs"
	"github.com/mjpancake/ih/ako/sc"
	"github.com/mjpancake/ih/hayari"
	"github.com/mjpancake/ih/mako"
	"github.com/mjpancake/ih/nodoka"
)

const uuReqTmot = 1 * time.Second

func (ussn *ussn) handleReject(msg error) {
	// good bye, no need to check error anymore
	jsonb, _ := json.Marshal(sc.NewAuthFail(msg.Error()))
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
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}

	err = hayari.Write(ussn.conn, jsonb)
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
	sc := sc.NewUpdateUser(ussn.user, mako.GetStats(ussn.user.Id))
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
	case *cs.HeartBeat:
		// do nothing
	case *cs.Book:
		nodoka.Bmgr.Tell(&nodoka.MbBook{
			Uid:      ussn.user.Id,
			BookType: msg.BookType,
		})
	case *cs.Unbook:
		nodoka.Bmgr.Tell(&nodoka.MbUnbook{Uid: ussn.user.Id})
	case *cs.GetReplayList:
		ussn.handleGetReplayList()
	case *cs.GetReplay:
		ussn.handleGetReplay(msg.ReplayId)
	case *cs.Choose:
		ussn.handleChoose(msg.GirlIndex)
	case *cs.Ready:
		ussn.handleReady()
	case *cs.Action:
		ussn.handleAction(msg)
	default:
		ussn.handleError(fmt.Errorf("unexpected CsMsg %T\n", msg))
	}
}

func (ussn *ussn) handleLookAround() {
	playing, err := (&nodoka.MtHasUser{Uid: ussn.user.Id}).Req()
	if err != nil {
		ussn.handleError(err)
		return
	}

	if playing {
		ussn.handleSc(&sc.TypeOnly{"resume"}, noResp)
	} else {
		res, err := nodoka.Umgr.RequestFuture(&cpWater{}, uuReqTmot).Result()
		if err != nil {
			ussn.handleError(err)
			return
		}
		water := res.(*pcWater)

		tables, err := (&nodoka.MtCtPlays{}).Req()
		if err != nil {
			ussn.handleError(err)
			return
		}
		waits, err := (&nodoka.MbCtBooks{}).Req()
		if err != nil {
			ussn.handleError(err)
			return
		}

		user := ussn.user
		dcbaBookable := [4]bool{
			user.Level < 13 || user.Rating < 1800.0,
			user.Level >= 9,
			user.Level >= 13 && user.Rating >= 1800.0,
			user.Level >= 16 && user.Rating >= 2000.0,
		}

		msg := sc.NewLookAround(water.ct, water.water,
			&dcbaBookable, &waits, &tables)
		ussn.handleSc(msg, noResp)
	}
}

func (ussn *ussn) handleGetReplayList() {
	ids := mako.GetReplayList(ussn.user.Id)
	ussn.handleSc(sc.NewGetReplayList(ids), noResp)
}

func (ussn *ussn) handleGetReplay(replayId uint) {
	text, err := mako.GetReplay(replayId)
	if err != nil {
		ussn.handleError(err)
		return
	}
	time.Sleep(2 * time.Second) // no reason, just wanna sleep
	ussn.handleSc(sc.NewGetReplay(replayId, text), noResp)
}

func (ussn *ussn) handleChoose(gidx int) {
	nodoka.Tmgr.Tell(&nodoka.MtChoose{Uid: ussn.user.Id, Gidx: gidx})
}

func (ussn *ussn) handleReady() {
	nodoka.Tmgr.Tell(&nodoka.MtReady{Uid: ussn.user.Id})
}

func (ussn *ussn) handleAction(msg *cs.Action) {
	nodoka.Tmgr.Tell(&nodoka.MtAction{Uid: ussn.user.Id, Act: msg})
}
