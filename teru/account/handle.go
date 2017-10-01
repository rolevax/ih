package account

import (
	"log"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/msg"
)

func PostCreate(request *restful.Request, response *restful.Response) {
	slow()

	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAccountCreate{}
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

func PostActivate(request *restful.Request, response *restful.Response) {
	slow()

	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAccountActivate{}
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

func GetCPoints(request *restful.Request, response *restful.Response) {
	sc := &msg.ScCpoints{}
	defer response.WriteEntity(sc)

	log.Println("account/getCPoints")
	sc.Entries = mako.GetCPoints()
}

func slow() {
	time.Sleep(3 * time.Second)
}
