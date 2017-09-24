package mako

import (
	"log"

	"github.com/rolevax/ih/ako/model"
)

func GetCpoints() []model.CpointEntry {
	var res []model.CpointEntry

	err := db.Model(&res).Order("cpoint DESC").Select()
	if err != nil {
		log.Fatalln("mako.GetCpoints", err)
	}

	return res
}
