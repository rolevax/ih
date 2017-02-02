package srv

import (
	"log"
)

type tables struct {
	create		chan [4]uid
	ready		chan uid
	action		chan *reqAction
	endSession	chan *session
	conns		*conns
	sessions	[]*session
}

func newTables(conns *conns) *tables {
	tables := new(tables)

	tables.create = make(chan [4]uid)
	tables.ready = make(chan uid)
	tables.action = make(chan *reqAction)
	tables.endSession = make(chan *session)
	tables.conns = conns
	tables.sessions = make([]*session, 16)[0:0]

	return tables
}

func (tables *tables) Loop() {
	for {
		select {
		case uids := <-tables.create:
			tables.add(uids)
		case uid := <-tables.ready:
			tables.markReady(uid)
		case act := <-tables.action:
			tables.route(act)
		case s := <-tables.endSession:
			tables.sub(s)
		}
	}
}

func (tables *tables) Create() chan<- [4]uid {
	return tables.create
}

func (tables *tables) Ready() chan<- uid {
	return tables.ready
}

func (tables *tables) Action() chan<- *reqAction {
	return tables.action
}

func (tables *tables) EndSession() chan<- *session {
	return tables.endSession
}

func (tables *tables) HasUser(uid uid) bool {
	for _, s := range tables.sessions {
		if _, ok := s.FindUser(uid); ok {
			return true
		}
	}
	return false
}

func (tables *tables) add(uids [4]uid) {
	s := newSession(tables, uids)
	tables.sessions = append(tables.sessions, s)
	go s.Loop()
}

func (tables *tables) sub(s *session) {
	i := 0
	for i < len(tables.sessions) && tables.sessions[i] != s {
		i++
	}
	if i == len(tables.sessions) {
		log.Fatalln("tables.sub: session not found")
	}
	// overwrite by back and pop back
	last := len(tables.sessions) - 1;
	tables.sessions[i] = tables.sessions[last]
	tables.sessions[last] = nil
	tables.sessions = tables.sessions[:last]
}

func (tables *tables) markReady(uid uid) {
	for _, s := range tables.sessions {
		if i, ok := s.FindUser(uid); ok {
			s.Ready() <- i
			return
		}
	}
	log.Println("tables.ready", uid, "not found")
}

func (tables *tables) route(act *reqAction) {
	for _, s := range tables.sessions {
		if _, ok := s.FindUser(act.uid); ok {
			s.Action() <- act
			return
		}
	}
	log.Println("tables.route", act.uid, "not found")
}

