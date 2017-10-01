package admin

import (
	"log"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/teru/msg"
)

func PostCPoint(request *restful.Request, response *restful.Response) {
	slow()

	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAdminCPoint{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	/*
		if !mako.CheckAdminToken(cs.Token) {
			sc.Error = "wrong token"
			return
		}
	*/

	log.Println("admin/c-point", cs.Username, cs.CPointDelta)

	/* FUCK temp
	err = mako.UpdateCPoint(cs.Username, cs.CPointDelta)
	if err != nil {
		sc.Error = err.Error()
	}
	*/
}

func slow() {
	time.Sleep(3 * time.Second)
}
