package main

import (
	"log"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/rolevax/ih/mako"
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

	err = mako.SignUp(cs.Username, cs.Password)
	if err != nil {
		sc.Error = err.Error()
	}
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

	log.Println("account/activate", cs.Username, "answers", cs.Answers)
	err = mako.Activate(cs.Username, cs.Password, cs.Answers)
	if err != nil {
		sc.Error = err.Error()
	}
}

func slow() {
	time.Sleep(3 * time.Second)
}
