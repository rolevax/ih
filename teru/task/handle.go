package task

import (
	"fmt"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/account"
	"github.com/rolevax/ih/teru/msg"
)

func GetRoot(req *restful.Request, resp *restful.Response) {
	sc := &msg.ScTaskRoot{}
	defer resp.WriteEntity(sc)

	tasks, err := mako.GetTasks()
	if err != nil {
		sc.Error = err.Error()
		return
	}

	sc.Tasks = tasks
	sc.Waters = mako.GetTaskWaters(30)
}

func GetTask(req *restful.Request, resp *restful.Response) {
	sc := &msg.ScTask{}
	defer resp.WriteEntity(sc)

	taskId, err := getTaskId(req)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	task, err := mako.GetTask(taskId)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	sc.Task = task
}

func PostStart(req *restful.Request, resp *restful.Response) {
	time.Sleep(1 * time.Second)
	sc := &msg.Sc{}
	defer resp.WriteEntity(sc)

	taskId, err := getTaskId(req)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	user, err := getUser(req)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	err = mako.StartTask(user.Id, taskId)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	mako.AddTaskWater(fmt.Sprintf(
		"%v %v开始执行任务[%v]",
		time.Now().Format("2006-01-02 15:04"),
		user.Username,
		taskId,
	))
}

func PostPr(req *restful.Request, resp *restful.Response) {
	time.Sleep(1 * time.Second)
	sc := &msg.Sc{}
	defer resp.WriteEntity(sc)

	taskId, err := getTaskId(req)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	user, err := getUser(req)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	err = mako.NotifyCheckTask(user.Id, taskId)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	mako.AddTaskWater(fmt.Sprintf(
		"%v %v通知验收任务[%v]",
		time.Now().Format("2006-01-02 15:04"),
		user.Username,
		taskId,
	))
}

func getTaskId(req *restful.Request) (int, error) {
	taskIdStr := req.PathParameter("task-id")
	taskId, err := strconv.Atoi(taskIdStr)
	return taskId, err
}

func getUser(req *restful.Request) (*model.User, error) {
	uid, err := account.GetUid(req)
	if err != nil {
		return nil, err
	}

	return mako.GetUser(uid)
}
