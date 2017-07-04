package main

import (
	"log"
	"time"

	"github.com/emicklei/go-restful"
)

func create(request *restful.Request, response *restful.Response) {
	slow()

	sc := &Sc{}
	defer response.WriteEntity(sc)

	cs := &CsAccountCreate{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	log.Println("account/create", cs.Username)
	// TODO check and add to db
}

func activate(request *restful.Request, response *restful.Response) {
	slow()

	sc := &Sc{}
	defer response.WriteEntity(sc)

	cs := &CsAccountActivate{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	// TODO check password and set error
	log.Println("account/activate", cs.Username, "answers", cs.Answers)
	// TODO set sc.Result
}

func slow() {
	time.Sleep(3 * time.Second)
}
