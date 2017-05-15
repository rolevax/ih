package ussn

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/netio"
	"github.com/mjpancake/hisa/node"
)

const uuReqTmot = 1 * time.Second

func (ussn *ussn) handleReject(msg error) {
	// good bye, no need to check error anymore
	jsonb, _ := json.Marshal(model.NewScAuthFail(msg.Error()))
	netio.Write(ussn.conn, jsonb)
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

	err = netio.Write(ussn.conn, jsonb)
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
	ussn.user = db.GetUser(ussn.user.Id)
	sc := model.NewScUpdateUser(ussn.user, db.GetStats(ussn.user.Id))
	ussn.handleSc(sc, noResp)
}

func (ussn *ussn) handleRead(breq []byte) {
	msg, err := model.FromJson(breq)
	if err != nil {
		log.Println(ussn.user.Id, "-X->", string(breq))
		ussn.handleError(err)
	} else {
		ussn.handleCs(msg)
	}
}

func (ussn *ussn) handleCs(i interface{}) {
	switch msg := i.(type) {
	case *model.CsLookAround:
		ussn.handleLookAround()
	case *model.CsHeartBeat:
		// do nothing
	case *model.CsBook:
		node.Bmgr.Tell(&node.MbBook{
			Uid:      ussn.user.Id,
			BookType: msg.BookType,
		})
	case *model.CsUnbook:
		node.Bmgr.Tell(&node.MbUnbook{Uid: ussn.user.Id})
	case *model.CsGetReplayList:
		ussn.handleGetReplayList()
	case *model.CsGetReplay:
		ussn.handleGetReplay(msg.ReplayId)
	case *model.CsChoose:
		ussn.handleChoose(msg.GirlIndex)
	case *model.CsReady:
		ussn.handleReady()
	case *model.CsAction:
		ussn.handleAction(msg)
	default:
		ussn.handleError(fmt.Errorf("unexpected CsMsg %T\n", msg))
	}
}

func (ussn *ussn) handleLookAround() {
	playing, err := (&node.MtHasUser{Uid: ussn.user.Id}).Req()
	if err != nil {
		ussn.handleError(err)
		return
	}

	if playing {
		ussn.handleSc(&model.ScTypeOnly{"resume"}, noResp)
	} else {
		res, err := node.Umgr.RequestFuture(&cpCtUser{}, uuReqTmot).Result()
		if err != nil {
			ussn.handleError(err)
			return
		}
		connCt := res.(int)
		msg := model.NewScLookAround(connCt)

		pss, err := (&node.MtCtPlays{}).Req()
		if err != nil {
			ussn.handleError(err)
			return
		}
		bss, err := (&node.MbCtBooks{}).Req()
		if err != nil {
			ussn.handleError(err)
			return
		}

		user := ussn.user
		cBookable := user.Level >= 9
		bBookable := user.Level >= 13 && user.Rating >= 1800.0
		dBookable := !bBookable
		aBookable := user.Level >= 16 && user.Rating >= 2000.0
		msg.Books[0] = model.BookEntry{dBookable, bss[0], 4 * pss[0]}
		msg.Books[1] = model.BookEntry{cBookable, bss[1], 4 * pss[1]}
		msg.Books[2] = model.BookEntry{bBookable, bss[2], 4 * pss[2]}
		msg.Books[3] = model.BookEntry{aBookable, bss[3], 4 * pss[3]}
		ussn.handleSc(msg, noResp)
	}
}

func (ussn *ussn) handleGetReplayList() {
	ids := db.GetReplayList(ussn.user.Id)
	ussn.handleSc(model.NewScGetReplayList(ids), noResp)
}

func (ussn *ussn) handleGetReplay(replayId uint) {
	text, err := db.GetReplay(replayId)
	if err != nil {
		ussn.handleError(err)
		return
	}
	time.Sleep(2 * time.Second) // no reason, just wanna sleep
	ussn.handleSc(model.NewScGetReplay(replayId, text), noResp)
}

func (ussn *ussn) handleChoose(gidx int) {
	node.Tmgr.Tell(&node.MtChoose{Uid: ussn.user.Id, Gidx: gidx})
}

func (ussn *ussn) handleReady() {
	node.Tmgr.Tell(&node.MtReady{Uid: ussn.user.Id})
}

func (ussn *ussn) handleAction(msg *model.CsAction) {
	node.Tmgr.Tell(&node.MtAction{Uid: ussn.user.Id, Act: msg})
}
