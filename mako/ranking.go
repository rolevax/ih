package mako

import (
	"log"

	"github.com/mjpancake/ih/ako/model"
)

func average(d *[4]float64) float64 {
	return (d[0] + d[1] + d[2] + d[3]) / 4.0
}

func updateLpr(lprs *[4]*model.Lpr, ranks [4]int, plays [4]int) {
	sumRating := 0.0
	for w := 0; w < 4; w++ {
		if lprs[w] == nil {
			log.Fatalln("updateLpr nil ptr")
		}
		sumRating += lprs[w].Rating
	}
	avgRating := sumRating / 4.0
	bases := [4]float64{30.0, 10.0, -10.0, -30.0}

	for w := 0; w < 4; w++ {
		playCoeff := 0.2
		play := plays[w]
		if play < 400 {
			playCoeff = (1 - float64(play)*0.002)
		}
		diffCoeff := (avgRating - lprs[w].Rating) / 40.0
		delta := playCoeff * (bases[ranks[w]-1] + diffCoeff)
		lprs[w].Rating += delta
	}

	for w := 0; w < 4; w++ {
		switch ranks[w] {
		case 1:
			updateTopPt(&lprs[w].Pt, &lprs[w].Level)
		case 2:
			update2ndPt(&lprs[w].Pt, &lprs[w].Level)
		case 3:
			// no change
		case 4:
			updateLastPt(&lprs[w].Pt, &lprs[w].Level)
		default:
			log.Fatalln("db.updateLpr: invalid rank", ranks[w])
		}
	}
}

func updateTopPt(pt *int, level *int) {
	*pt += 75
	updateLevel(pt, level)
}

func update2ndPt(pt *int, level *int) {
	*pt += 30
	updateLevel(pt, level)
}

func updateLastPt(pt *int, level *int) {
	diffs := [20]int{
		0, 0, 0, 0, 0,
		0, 0, 0, -15, -30,
		-45, -60, -75, -90, -105,
		-120, -135, -150, -165, -180}
	*pt += diffs[*level]
	updateLevel(pt, level)
}

func updateLevel(pt *int, level *int) {
	maxs := [20]int{
		30, 30, 30, 60, 60,
		60, 90, 100, 100, 100,
		400, 800, 1200, 1600, 2000,
		2400, 2800, 3200, 3600, 4000}

	starts := [20]int{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		200, 400, 600, 800, 1000,
		1200, 1400, 1600, 1800, 2000}

	if *pt >= maxs[*level] {
		if *level+1 < len(maxs) {
			(*level)++
			*pt = starts[*level]
		} else {
			*pt = maxs[*level]
		}
	} else if *pt < 0 {
		if *level >= 10 {
			(*level)--
			*pt = starts[*level]
		} else {
			*pt = 0
		}
	}
}
