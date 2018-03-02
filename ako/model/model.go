package model

// girl id, signed-int for compatibility to libsaki
type Gid int

type RuleId int

const (
	RuleFourDoges   RuleId = 0
	RuleClassic1In2 RuleId = 1
	RuleMax         RuleId = 2
)

func (ri RuleId) Valid() bool {
	i := int(ri)
	return 0 <= i && i < int(RuleMax)
}

// level, pt, and rating
type Lpr struct {
	Level  int
	Pt     int
	Rating float64
}

// deprecated
type Girl struct {
	Id Gid
	Lpr
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
