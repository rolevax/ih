package model

// user id
type Uid uint

const (
	UidSys Uid = 500
	UidAi1 Uid = 501
	UidAi2 Uid = 502
)

func (uid Uid) IsBot() bool {
	return uint(uid) < 1000
}

func (uid Uid) IsHuman() bool {
	return !uid.IsBot()
}

// girl id, signed-int for compatibility to libsaki
type Gid int

// level, pt, and rating
type Lpr struct {
	Level  int
	Pt     int
	Rating float64
}

type User struct {
	Id       Uid `sql:"user_id,pk"`
	Username string
	CPoint   int
}

// deprecated
type Girl struct {
	Id Gid
	Lpr
}

type CPointEntry struct {
	tableName struct{} `sql:"users"`
	Username  string
	CPoint    int
}

type TaskState int

const (
	TaskStateToDo    TaskState = 0
	TaskStateDoing   TaskState = 1
	TaskStateToCheck TaskState = 2
	TaskStateClosed  TaskState = 3
)

type Task struct {
	Id         int `sql:"task_id,pk"`
	Title      string
	Content    string
	State      TaskState
	AssigneeId Uid
	Assignee   *User
	CPoint     int
}
