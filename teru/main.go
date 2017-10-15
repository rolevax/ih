package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/teru/account"
	"github.com/rolevax/ih/teru/admin"
	"github.com/rolevax/ih/teru/my"
)

const (
	Port        = ":8080"
	CertPath    = "/srv/cert.pem"
	PrivKeyPath = "/srv/key.pem"
)

func main() {
	addWebService()
	supportCors()
	log.Fatal(http.ListenAndServeTLS(Port, CertPath, PrivKeyPath, nil))
}

func addWebService() {
	restful.Filter(globalLogging)
	addWebServiceAccount()
	addWebServiceAdmin()
	addWebServiceMy()
}

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("%s %s %v\n", req.Request.Method, req.Request.URL, req.Request.RemoteAddr)
	chain.ProcessFilter(req, resp)
}

func addWebServiceAccount() {
	ws := &restful.WebService{}
	ws.
		Path("/account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/auth").To(account.PostAuth))
	ws.Route(ws.POST("/create").To(account.PostCreate))

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

func addWebServiceMy() {
	ws := &restful.WebService{}
	ws.
		Path("/my").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/null").Filter(account.FilterAuth).To(my.GetNull))

	restful.Add(ws)
}

func supportCors() {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer,
	}
	restful.Filter(cors.Filter)
	restful.Filter(restful.OPTIONSFilter())
}
