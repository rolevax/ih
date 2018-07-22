package mako

import (
	"github.com/rolevax/ih/ako/model"
)

func UpdateUserGirl(uids [4]model.Uid, args *model.EndTableStat) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = updateReplay(tx, uids, args.Replay)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
