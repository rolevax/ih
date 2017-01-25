package srv

import (
	"log"
)

type tables struct {
	create		chan [4]uid
	ready		chan uid
	action		chan *reqAction
	conns		*conns
	sessions	[]*session
}

func newTables(conns *conns) *tables {
	tables := new(tables)

	tables.create = make(chan [4]uid)
	tables.ready = make(chan uid)
	tables.action = make(chan *reqAction)
	tables.conns = conns
	tables.sessions = make([]*session, 16)[0:0]

	return tables
}

func (tables *tables) loop() {
	for {
		select {
		case uids := <-tables.create:
			tables.add(uids)
		case uid := <-tables.ready:
			tables.markReady(uid)
		case act := <-tables.action:
			tables.route(act)
		}
	}
}

func (tables *tables) add(uids [4]uid) {
	s := newSession(tables, uids)
	tables.sessions = append(tables.sessions, s)
	s.notifyLoad()
}

func (tables *tables) markReady(uid uid) {
	for _, s := range tables.sessions {
		if i, ok := s.findUser(uid); ok {
			s.readys[i] = true
			if s.readys[0] && s.readys[1] && s.readys[2] && s.readys[3] {
				s.start()
			}
			return
		}
	}
	log.Println("tables.ready", uid, "not found")
}

func (tables *tables) route(act *reqAction) {
	for i, s := range tables.sessions {
		if _, ok := s.findUser(act.uid); ok {
			s.action(act)
			if s.gameOver() {
				s.destroy()
				// overwrite by back and pop back
				last := len(tables.sessions) - 1;
				tables.sessions[i] = tables.sessions[last]
				tables.sessions[last] = nil
				tables.sessions = tables.sessions[:last]
			}
			return
		}
	}
	log.Println("tables.route", act.uid, "not found")
}

