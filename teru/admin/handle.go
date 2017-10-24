package admin

import (
	"fmt"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/msg"
)

func PostCPoint(request *restful.Request, response *restful.Response) {
	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAdminCPoint{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	if !mako.CheckAdminToken(cs.Token) {
		sc.Error = "wrong token"
		return
	}

	err = mako.UpdateCPoint(cs.Username, cs.CPointDelta)
	if err != nil {
		sc.Error = err.Error()
	}
}

func PostUpsertTask(request *restful.Request, response *restful.Response) {
	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAdminUpsertTask{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	if !mako.CheckAdminToken(cs.Token) {
		sc.Error = "wrong token"
		return
	}

	err = mako.UpsertTask(&cs.Task)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	mako.AddTaskWater(fmt.Sprintf(
		"%v 任务更新 [%v] %v",
		time.Now().Format("2006-01-02 15:04"),
		cs.Task.Id,
		cs.Task.Title,
	))
}

func PostCheckTask(req *restful.Request, resp *restful.Response) {
	sc := &msg.Sc{}
	defer resp.WriteEntity(sc)

	cs := &msg.CsAdminCheckTask{}
	err := req.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	if !mako.CheckAdminToken(cs.Token) {
		sc.Error = "wrong token"
		return
	}

	waterStr := ""

	switch cs.Op {
	case "accept":
		if err := mako.AcceptWorkOnTask(cs.TaskId); err != nil {
			sc.Error = err.Error()
			return
		}
		waterStr = "通过，已更新贡献度"
	case "expect":
		if err := mako.ExpectAssigneeOnTask(cs.TaskId); err != nil {
			sc.Error = err.Error()
			return
		}
		waterStr = "未通过，求改进"
	case "fire":
		if err := mako.FireAssigneeFromTask(cs.TaskId, false); err != nil {
			sc.Error = err.Error()
			return
		}
		waterStr = "未通过，任务重新发布"
	case "fire-doing":
		if err := mako.FireAssigneeFromTask(cs.TaskId, true); err != nil {
			sc.Error = err.Error()
			return
		}
	default:
		sc.Error = fmt.Sprintf("unexpected op %s", cs.Op)
		return
	}

	if cs.Op == "fire-doing" {
		mako.AddTaskWater(fmt.Sprintf(
			"%v 任务[%v]被强制撤回",
			time.Now().Format("2006-01-02 15:04"),
			cs.TaskId,
		))
	} else {
		mako.AddTaskWater(fmt.Sprintf(
			"%v 任务[%v]验收结果: %v",
			time.Now().Format("2006-01-02 15:04"),
			cs.TaskId,
			waterStr,
		))
	}
}
