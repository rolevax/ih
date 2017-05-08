package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/mjpancake/hisa/model"
)

func GetReplayList(uid model.Uid) []uint {
	var ids []uint

	rows, err := db.Query(
		`select replay_id from user_replay where user_id=? order by replay_id desc`,
		uid)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uint
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalln(err)
		}
		ids = append(ids, id)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return ids
}

func GetReplay(replayId uint) (string, error) {
	var text string

	err := db.QueryRow(
		"select content from replays where replay_id=?", replayId).
		Scan(&text)

	return text, err
}

func updateReplay(tx *sql.Tx, uids [4]model.Uid,
	replay map[string]interface{}) error {
	jsonb, err := json.Marshal(replay)
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into replays(content) values (?)",
		string(jsonb))
	if err != nil {
		return err
	}

	_, err = tx.Exec(`insert into user_replay(user_id,replay_id) values
		(?, last_insert_id()),
		(?, last_insert_id()),
		(?, last_insert_id()),
		(?, last_insert_id())`,
		uids[0], uids[1], uids[2], uids[3])

	return err
}
