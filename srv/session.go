package srv

import (
	"log"
	"time"
	"bitbucket.org/rolevax/sakilogy-server/saki"
)

type session struct {
	table	saki.TableSession
	uids	[4]uid
	readys	[4]bool
	tables	*tables
}

func newSession(tables *tables, uids [4]uid) *session {
	var s session

	s.table = saki.NewTableSession()
	s.uids = uids
	s.tables = tables;

	return &s
}

func (s *session) notifyLoad() {
	var users [4]*user
	for i := range users {
		users[i] = s.tables.conns.dao.getUser(s.uids[i])
		if users[i] == nil {
			log.Fatal("session.nofityLoad:", s.uids[i], "not in DB")
		}
	}

	msg := struct {
		Type		string
		Users		[4]*user
		GirlIds		[4]int
		TempDealer	int
	}{"start", users, [4]int{0,0,0,0}, 0}

	for i, uid := range s.uids {
		msg.TempDealer = (4 - i) % 4
		s.tables.conns.peer <- &Mail{uid, msg}
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

func (s *session) findUser(uid uid) (int, bool) {
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
		s.tables.conns.peer <- &mail
	}
}

func (s *session) action(act *reqAction) {
	i, _ := s.findUser(act.uid)
	mails := s.table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)

	size := int(mails.Size())
	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		msg := mails.Get(i).GetMsg()
		if msg == "auto" { // special mark
			go func() {
				time.Sleep(300 * time.Millisecond)
				act := reqAction{s.uids[toWhom], "SPIN_OUT", "-1"}
				s.tables.action <- &act
			}()
		} else {
			mail := Mail{s.uids[toWhom], msg}
			s.tables.conns.peer <- &mail
		}
	}
}

func (s *session) gameOver() bool {
	return s.table.GameOver()
}

func (s *session) destroy()  {
	saki.DeleteTableSession(s.table)
}

