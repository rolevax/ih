package mako

import (
	"encoding/json"
	"log"

	"github.com/go-pg/pg"
	"github.com/rolevax/ih/ako/model"
)

func GetReplayList(uid model.Uid) []uint {
	var ids []uint

	_, err := db.Query(
		pg.Scan(pg.Array(&ids)),
		`SELECT replay_ids FROM replay_of_user WHERE user_id=?`,
		uid,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return ids
}

func GetReplay(replayId uint) (string, error) {
	var text string

	_, err := db.QueryOne(
		&text,
		"SELECT content FROM replays WHERE replay_id=?",
		replayId,
	)

	return text, err
}

func updateReplay(tx *pg.Tx, uids [4]model.Uid,
	replay map[string]interface{}) error {
	jsonb, err := json.Marshal(replay)
	if err != nil {
		return err
	}

	var rid uint
	_, err = tx.QueryOne(
		&rid,
		"INSERT INTO replays(content) VALUES (?) RETURNING replay_id",
		string(jsonb),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO replay_of_user(user_id, replay_ids)
		VALUES (?, ?), (?, ?), (?, ?), (?, ?)
		ON CONFLICT (user_id) DO UPDATE
		SET replay_ids = ? || replay_ids`,
		uids[0], rid, uids[1], rid, uids[2], rid, uids[3], rid,
		rid,
	)

	return err
}
