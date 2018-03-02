package model

import (
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/rolevax/ih/hitomi"
)

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

type User struct {
	Id        Uid `sql:"user_id,pk"`
	Username  string
	CPoint    int
	GotFoodAt *time.Time
	Food      int
}

func (user *User) AfterQuery(db orm.DB) error {
	user.Username = hitomi.Filter(user.Username)
	return nil
}

type CPointEntry struct {
	tableName struct{} `sql:"users"`
	Username  string
	CPoint    int
}

func (cpe *CPointEntry) AfterQuery(db orm.DB) error {
	cpe.Username = hitomi.Filter(cpe.Username)
	return nil
}

type FoodChange struct {
	Delta  int
	Reason string
}
