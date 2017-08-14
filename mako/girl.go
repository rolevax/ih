package mako

import (
	"github.com/rolevax/ih/ako/model"
)

// deprecated
func GetRankedGids() []model.Gid {
	var gids []model.Gid

	// excluding doge
	for i := 0; i < 22; i++ {
		gids = append(gids, 0)
	}

	return gids
}
