package mako

import (
	"fmt"

	"github.com/rolevax/ih/ako/model"
)

func UpsertTask(task *model.Task) error {
	_, err := db.Model(task).
		OnConflict("(task_id) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("content = EXCLUDED.content").
		Set("c_point = EXCLUDED.c_point").
		Insert()

	return err
}

func GetTasks() ([]model.Task, error) {
	res := []model.Task{}
	err := db.Model(&res).
		Column("task.*", "Assignee").
		Where("task.state <= 2").
		Select()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetTask(taskId int) (*model.Task, error) {
	res := &model.Task{Id: taskId}
	err := db.Select(res)
	return res, err
}

func StartTask(uid model.Uid, taskId int) error {
	ct, err := db.Model(&model.Task{}).
		Where("assignee_id=?", uid).
		Count()
	if err != nil {
		return err
	}
	if ct != 0 {
		return fmt.Errorf("当前任务验收通过后才能接受新任务")
	}

	_, err = db.Model(&model.Task{}).
		Set("state=?", model.TaskStateDoing).
		Set("assignee_id=?", uid).
		Where("task_id=? AND state=?", taskId, model.TaskStateToDo).
		Update()
	return err
}

func NotifyCheckTask(uid model.Uid, taskId int) error {
	_, err := db.Model(&model.Task{}).
		Set("state=?", model.TaskStateToCheck).
		Where("task_id=? AND state=?", taskId, model.TaskStateDoing).
		Update()
	return err
}

func AcceptWorkOnTask(taskId int) error {
	task := &model.Task{Id: taskId}
	err := db.Select(task)
	if err != nil {
		return err
	}

	res, err := db.Model(&model.Task{}).
		Set("state=?", model.TaskStateClosed).
		Set("assignee_id=NULL").
		Where("task_id=? AND state=?", taskId, model.TaskStateToCheck).
		Update()
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("task: %d rows affected", res.RowsAffected())
	}

	res, err = db.Model(&model.User{}).
		Set("c_point=c_point+?", task.CPoint).
		Where("user_id=?", task.AssigneeId).
		Update()
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("user: %d rows affected", res.RowsAffected())
	}

	return nil
}

func ExpectAssigneeOnTask(taskId int) error {
	res, err := db.Model(&model.Task{}).
		Set("state=?", model.TaskStateDoing).
		Where("task_id=? AND state=?", taskId, model.TaskStateToCheck).
		Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return fmt.Errorf("%d rows affected", res.RowsAffected())
	}

	return nil
}

func FireAssigneeFromTask(taskId int, isDoing bool) error {
	var currState model.TaskState
	if isDoing {
		currState = model.TaskStateDoing
	} else {
		currState = model.TaskStateToCheck
	}

	res, err := db.Model(&model.Task{}).
		Set("state=?", model.TaskStateToDo).
		Set("assignee_id=NULL").
		Where("task_id=? AND state=?", taskId, currState).
		Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return fmt.Errorf("%d rows affected", res.RowsAffected())
	}

	return nil
}
