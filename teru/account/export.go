package account

import (
	"strconv"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/ako/model"
)

func GetUid(req *restful.Request) (model.Uid, error) {
	uidStr := req.Request.Header.Get("X-User-Id")
	uidUint, err := strconv.ParseUint(uidStr, 0, 0)
	if err != nil {
		return model.Uid(0), err
	}

	return model.Uid(uidUint), nil
}
