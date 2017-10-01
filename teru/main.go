package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/rolevax/ih/teru/account"
	"github.com/rolevax/ih/teru/admin"
)

func main() {
	addWebService()
	supportCors()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addWebService() {
	addWebServiceAccount()
	addWebServiceAdmin()
}

func addWebServiceAccount() {
	ws := &restful.WebService{}
	ws.
		Path("/account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/create").To(account.PostCreate))
	ws.Route(ws.POST("/activate").To(account.PostActivate))

	ws.Route(ws.GET("/c-points").To(account.GetCPoints))

	restful.Add(ws)
}

func addWebServiceAdmin() {
	ws := &restful.WebService{}
	ws.
		Path("/admin").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/c-point").To(admin.PostCPoint))

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
