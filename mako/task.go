package mako

import (
	"errors"

	"github.com/rolevax/ih/ako/model"
)

func UpsertTask(task *model.Task) error {
	return errors.New("任务系统已关闭")
}

func GetTasks() ([]model.Task, error) {
	return nil, errors.New("任务系统已关闭")
}

func GetTask(taskId int) (*model.Task, error) {
	return nil, errors.New("任务系统已关闭")
}

func StartTask(uid model.Uid, taskId int) error {
	return errors.New("任务系统已关闭")
}

func NotifyCheckTask(uid model.Uid, taskId int) error {
	return errors.New("任务系统已关闭")
}

func AcceptWorkOnTask(taskId int) error {
	return errors.New("任务系统已关闭")
}

func ExpectAssigneeOnTask(taskId int) error {
	return errors.New("任务系统已关闭")
}

func FireAssigneeFromTask(taskId int, isDoing bool) error {
	return errors.New("任务系统已关闭")
}
