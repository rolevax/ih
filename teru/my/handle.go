package my

import (
	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/teru/msg"
)

func GetNull(request *restful.Request, response *restful.Response) {
	sc := &msg.Sc{}
	defer response.WriteEntity(sc)
}
