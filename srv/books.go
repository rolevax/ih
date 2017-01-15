package srv

import (
	"log"
	"bitbucket.org/rolevax/sakilogy-server/model"
)

type Books struct {
	Book	chan model.Uid
	Unbook	chan model.Uid
	conns	*Conns
	waits	[4]model.Uid
	wait	int
}

func NewBooks(conns *Conns) *Books {
	var books Books

	books.Book = make(chan model.Uid)
	books.Unbook = make(chan model.Uid)
	books.conns = conns
	books.wait = 0

	return &books;
}

func (books *Books) Loop() {
	for {
		select {
		case uid := <-books.Book:
			books.book(uid)
		case uid := <-books.Unbook:
			books.unbook(uid)
		}
	}
}

func (books *Books) book(uid model.Uid) {
	log.Println("book", uid)
	books.waits[books.wait] = uid;
	books.wait++
	if books.wait == 4 {
		books.conns.Start <- books.waits
		books.wait = 0
	}
}

func (books *Books) unbook(uid model.Uid) {
	log.Println("unbook", uid)
	i := 0
	for i < books.wait && books.waits[i] != uid {
		i++
	}

	if i == books.wait {
		log.Println("unbook", uid, "not found")
		return
	}

	// swap to back, and pop back
	e := books.wait - 1;
	books.waits[i], books.waits[e] = books.waits[e], books.waits[i];
	books.wait--
}



