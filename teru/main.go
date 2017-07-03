package main

import (
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

type CsAccountCreate struct {
	Username string
	Password string
}

type CsAccountActivate struct {
	Username string
	Password string
	Answers  string
}

type Sc struct {
	Error string // no news is good news
}

type ScAccountActivate struct {
	Sc
	Result string
}

func main() {
	ws := &restful.WebService{}
	ws.
		Path("/account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/create").To(create))
	ws.Route(ws.POST("/activate").To(activate))

	restful.Add(ws)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer,
	}
	restful.Filter(cors.Filter)
	restful.Filter(restful.OPTIONSFilter())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

	log.Println("create: username", cs.Username)
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
	log.Println("activate: username", cs.Username, "answers", cs.Answers)
	// TODO set sc.Result
}

func slow() {
	time.Sleep(3 * time.Second)
}
