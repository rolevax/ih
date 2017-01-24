package srv

import (
	"log"
	"time"
	"bitbucket.org/rolevax/sakilogy-server/model"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

type Action struct {
	Uid			model.Uid
	ActStr		string
	ActArg		string
}

type Tables struct {
	Create		chan [4]model.Uid
	Ready		chan model.Uid
	Action		chan *Action
	conns		*Conns
	sessions	[]*session
}

func NewTables(conns *Conns) *Tables {
	var tables Tables

	tables.Create = make(chan [4]model.Uid)
	tables.Ready = make(chan model.Uid)
	tables.Action = make(chan *Action)
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
		case act := <-tables.Action:
			tables.action(act)
		}
	}
}

func (tables *Tables) SessionCount() int {
	return len(tables.sessions)
}

func (tables *Tables) add(uids [4]model.Uid) {
	s := newSession(tables, uids)
	tables.sessions = append(tables.sessions, s)
	s.notifyLoad()
}

func (tables *Tables) ready(uid model.Uid) {
	for _, s := range tables.sessions {
		if i, ok := s.findUser(uid); ok {
			s.readys[i] = true
			if s.readys[0] && s.readys[1] && s.readys[2] && s.readys[3] {
				s.start()
			}
			return
		}
	}
	log.Println("Tables.ready", uid, "not found")
}

func (tables *Tables) action(act *Action) {
	for i, s := range tables.sessions {
		if _, ok := s.findUser(act.Uid); ok {
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
	log.Println("Tables.action", act.Uid, "not found")
}



type session struct {
	table	saki.TableSession
	uids	[4]model.Uid
	readys	[4]bool
	tables	*Tables
}

func newSession(tables *Tables, uids [4]model.Uid) *session {
	var s session

	s.table = saki.NewTableSession()
	s.uids = uids
	s.tables = tables;

	return &s
}

func (s *session) notifyLoad() {
	var users [4]*model.User
	for i := range users {
		users[i] = s.tables.conns.dao.GetUser(s.uids[i])
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
		s.tables.conns.Peer <- &Mail{uid, msg}
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
}

func (s *session) findUser(uid model.Uid) (int, bool) {
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
		s.tables.conns.Peer <- &mail
	}
}

func (s *session) action(act *Action) {
	i, _ := s.findUser(act.Uid)
	mails := s.table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)

	size := int(mails.Size())
	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		msg := mails.Get(i).GetMsg()
		if msg == "auto" { // special mark
			go func() {
				time.Sleep(300 * time.Millisecond)
				act := Action{s.uids[toWhom], "SPIN_OUT", "-1"}
				s.tables.Action <- &act
			}()
		} else {
			mail := Mail{s.uids[toWhom], msg}
			s.tables.conns.Peer <- &mail
		}
	}
}

func (s *session) gameOver() bool {
	return s.table.GameOver()
}

func (s *session) destroy()  {
	saki.DeleteTableSession(s.table)
}



