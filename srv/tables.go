package srv

import (
	"log"
	"bitbucket.org/rolevax/sakilogy-server/model"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

type Tables struct {
	Create		chan [4]model.Uid
	conns		*Conns
	sessions	[]*session
}

func NewTables(conns *Conns) *Tables {
	var tables Tables

	tables.Create = make(chan [4]model.Uid)
	tables.conns = conns
	tables.sessions = make([]*session, 16)

	return &tables
}

func (tables *Tables) Loop() {
	for {
		select {
		case uids := <-tables.Create:
			tables.add(uids)
			log.Println("table created", uids)
		}
	}
}

func (tables *Tables) add(uids [4]model.Uid) {
	s := newSession(uids)
	tables.sessions = append(tables.sessions, s)
	s.start()
}

type session struct {
	table	saki.TableSession
	uids	[4]model.Uid
}

func newSession(uids [4]model.Uid) *session {
	var s session
	s.table = saki.NewTableSession()
	s.uids = uids
	return &s
}

func (s *session) start() {
	mails := s.table.Start()
	defer saki.DeleteMailVector(mails)
}




