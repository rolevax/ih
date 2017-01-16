package srv

import (
	"log"
	"bitbucket.org/rolevax/sakilogy-server/model"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

type Tables struct {
	Create		chan [4]model.Uid
	Ready		chan model.Uid
	conns		*Conns
	sessions	[]*session
}

func NewTables(conns *Conns) *Tables {
	var tables Tables

	tables.Create = make(chan [4]model.Uid)
	tables.Ready = make(chan model.Uid)
	tables.conns = conns
	tables.sessions = make([]*session, 16)[0:0]

	return &tables
}

func (tables *Tables) Loop() {
	for {
		select {
		case uids := <-tables.Create:
			tables.add(uids)
		case uid := <-tables.Ready:
			tables.ready(uid)
		}
	}
}

func (tables *Tables) add(uids [4]model.Uid) {
	s := newSession(tables.conns, uids)
	tables.sessions = append(tables.sessions, s)
	s.notifyLoad()
}

func (tables *Tables) ready(uid model.Uid) {
	for _, s := range tables.sessions {
		if i, ok := s.findUid(uid); ok {
			s.readys[i] = true
			if s.readys[0] && s.readys[1] && s.readys[2] && s.readys[3] {
				s.start()
			}
			return
		}
	}
	log.Println("Tables.ready", uid, "not found")
}

type session struct {
	table	saki.TableSession
	uids	[4]model.Uid
	readys	[4]bool
	conns	*Conns
}

func newSession(conns *Conns, uids [4]model.Uid) *session {
	var s session

	s.table = saki.NewTableSession()
	s.uids = uids
	s.conns = conns;

	return &s
}

func (s *session) notifyLoad() {
	var users [4]*model.User
	for i := range users {
		users[i] = s.conns.dao.GetUser(s.uids[i])
		if users[i] == nil {
			log.Fatal("session.nofityLoad:", s.uids[i], "not in DB")
		}
	}

	msg := struct {
		Type		string
		Users		[4]*model.User
		GirlIds		[4]int
		TempDealer	int
	}{"start", users, [4]int{0,0,0,0}, 0}

	for i, uid := range s.uids {
		msg.TempDealer = (4 - i) % 4
		s.conns.Peer <- &Mail{uid, msg}
		// rotate perspectives
		u0 := msg.Users[0]
		msg.Users[0] = msg.Users[1]
		msg.Users[1] = msg.Users[2]
		msg.Users[2] = msg.Users[3]
		msg.Users[3] = u0

		g0 := msg.GirlIds[0]
		msg.GirlIds[0] = msg.GirlIds[1]
		msg.GirlIds[1] = msg.GirlIds[2]
		msg.GirlIds[2] = msg.GirlIds[3]
		msg.GirlIds[3] = g0
	}
	// TODO check client startable, (set PClient* to PTableThread)
}

func (s *session) findUid(uid model.Uid) (int, bool) {
	for i, u := range s.uids {
		if u == uid {
			return i, true
		}
	}
	return -1, false
}

func (s *session) start() {
	mails := s.table.Start()
	defer saki.DeleteMailVector(mails)

	size := int(mails.Size())
	for i := 0; i < size; i++ {
		mail := Mail{s.uids[mails.Get(i).GetTo()], mails.Get(i).GetMsg()}
		log.Println("session.start get mail", mail.To, mail.Msg)
		s.conns.Peer <- &mail
	}
}




