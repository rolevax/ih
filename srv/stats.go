package srv

import (
	"log"
)

func average(d *[4]float64) float64 {
	return (d[0] + d[1] + d[2] + d[3]) / 4.0
}

func updateLpr(lprs *[4]*lpr, plays [4]int, bt bookType) {
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
			playCoeff = (1 - float64(play) * 0.002)
		}
		diffCoeff := (avgRating - lprs[w].Rating) / 40.0
		delta := playCoeff * (bases[w] + diffCoeff)
		lprs[w].Rating += delta
	}

	updateTopPt(&lprs[0].Pt, &lprs[0].Level, bt)
	update2ndPt(&lprs[1].Pt, &lprs[1].Level, bt)
	updateLastPt(&lprs[3].Pt, &lprs[3].Level)
}

func updateTopPt(pt *int, level *int, bookType bookType) {
	switch bookType {
	case 0:
		*pt += 45
	case 1:
		*pt += 60
	case 2:
		*pt += 75
	case 3:
		*pt += 90
	default:
		log.Fatalln("updateTopPt: unknown bookType")
	}
	updateLevel(pt, level)
}

func update2ndPt(pt *int, level *int, bookType bookType) {
	switch bookType {
	case 0:
		*pt += 0
	case 1:
		*pt += 15
	case 2:
		*pt += 30
	case 3:
		*pt += 45
	default:
		log.Fatalln("update2ndPt: unknown bookType")
	}
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
		if *level + 1 < len(maxs) {
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



