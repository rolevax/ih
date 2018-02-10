package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/account"
	"github.com/rolevax/ih/teru/admin"
	"github.com/rolevax/ih/teru/my"
	"github.com/rolevax/ih/teru/task"
)

const (
	Port        = ":8080"
	CertPath    = "/srv/cert.pem"
	PrivKeyPath = "/srv/key.pem"
	DictPath    = "/srv/dict.txt"
)

type logWriter struct{}

func (w logWriter) Write(bytes []byte) (int, error) {
	prefix := time.Now().Format("01/02 15:04:05")
	return fmt.Print(prefix, " ", string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	if flag.Parsed() {
		log.Fatalln("unexpected flag parse before main()")
	}

	redis := flag.String("redis", "localhost:6379", "redis server addr")
	db := flag.String("db", "localhost:5432", "pg db server addr")
	flag.Parse()
	mako.InitRedis(*redis)
	mako.InitDb(*db)
	hitomi.Init(DictPath)

	addWebService()
	supportCors()
	log.Fatal(http.ListenAndServeTLS(Port, CertPath, PrivKeyPath, nil))
}

func addWebService() {
	restful.Filter(globalLogging)
	addWebServiceAccount()
	addWebServiceAdmin()
	addWebServiceMy()
	addWebServiceTask()
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
	ws.Route(ws.POST("/upsert-task").To(admin.PostUpsertTask))
	ws.Route(ws.POST("/check-task").To(admin.PostCheckTask))

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

func addWebServiceTask() {
	ws := &restful.WebService{}
	ws.
		Path("/task").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(
		ws.GET("/").
			To(task.GetRoot),
	)
	ws.Route(
		ws.GET("/{task-id}").
			To(task.GetTask),
	)
	ws.Route(
		ws.POST("/start/{task-id}").
			Filter(account.FilterAuth).
			To(task.PostStart),
	)
	ws.Route(
		ws.POST("/pr/{task-id}").
			Filter(account.FilterAuth).
			To(task.PostPr),
	)

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
