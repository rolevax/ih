package mako

import (
	"log"

	"github.com/mjpancake/ih/ako/model"
)

func GetRankedGids() []model.Gid {
	var gids []model.Gid

	// excluding doge
	rows, err := db.Query(
		`select girl_id from girls where girl_id<>0 order by rating desc`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var gid model.Gid
		err := rows.Scan(&gid)
		if err != nil {
			log.Fatalln(err)
		}
		gids = append(gids, gid)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return gids
}
