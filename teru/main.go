package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {
	addWebService()
	supportCors()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addWebService() {
	addWebServiceAccount()
	addWebServiceQuery()
}

func addWebServiceAccount() {
	ws := &restful.WebService{}
	ws.
		Path("/account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/create").To(create))
	ws.Route(ws.POST("/activate").To(activate))

	restful.Add(ws)
}

func addWebServiceQuery() {
	ws := &restful.WebService{}
	ws.
		Path("/query").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/cpoints").To(getCpoints))

	restful.Add(ws)
}

func supportCors() {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer,
	}
	restful.Filter(cors.Filter)
	restful.Filter(restful.OPTIONSFilter())
}
