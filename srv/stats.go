package srv

import (
	"log"
)

func statUserRank(ordUids *[4]uid) {
	users := sing.Dao.GetUsers(ordUids)
	sumRating := 0.0
	for w := 0; w < 4; w++ {
		if users[w] == nil {
			log.Fatalln("uid", ordUids[w], "not in DB")
		}
		sumRating += users[w].Rating
	}
	avgRating := sumRating / 4.0
	bases := [4]float64{30.0, 10.0, -10.0, -30.0}

	for w := 0; w < 4; w++ {
		playCoeff := 0.2
		ranks := &users[w].Ranks
		play := ranks[0] + ranks[1] + ranks[2] + ranks[3]
		if play < 400 {
			playCoeff = (1 - float64(play) * 0.002)
		}
		diffCoeff := (avgRating - users[w].Rating) / 40.0
		delta := playCoeff * (bases[w] + diffCoeff)
		users[w].Rating += delta
	}

	updateTopPt(&users[0].Pt, &users[0].Level)
	update2ndPt(&users[1].Pt, &users[0].Level)
	updateLastPt(&users[3].Pt, &users[3].Level)

	for r := 0; r < 4; r++ {
		users[r].Ranks[r]++
	}

	sing.Dao.SetUsersRank(&users)
	for w := 0; w < 4; w++ {
		sing.UssnMgr.UpdateInfo(users[w])
	}
}

func statGirlRank(ordGids *[4]gid) {
	girls := sing.Dao.GetGirls(ordGids)
	sumRating := 0.0
	for w := 0; w < 4; w++ {
		if girls[w] == nil {
			log.Fatalln("gid", ordGids[w], "not in DB")
		}
		sumRating += girls[w].Rating
	}
	avgRating := sumRating / 4.0
	bases := [4]float64{30.0, 10.0, -10.0, -30.0}

	for w := 0; w < 4; w++ {
		playCoeff := 0.2
		ranks := &girls[w].Ranks
		play := ranks[0] + ranks[1] + ranks[2] + ranks[3]
		if play < 400 {
			playCoeff = (1 - float64(play) * 0.002)
		}
		diffCoeff := (avgRating - girls[w].Rating) / 40.0
		delta := playCoeff * (bases[w] + diffCoeff)
		girls[w].Rating += delta
	}

	updateTopPt(&girls[0].Pt, &girls[0].Level)
	update2ndPt(&girls[1].Pt, &girls[1].Level)
	updateLastPt(&girls[3].Pt, &girls[1].Level)

	for r := 0; r < 4; r++ {
		girls[r].Ranks[r]++
	}

	sing.Dao.SetGirlsRank(&girls)
}

func average(d *[4]float64) float64 {
	return (d[0] + d[1] + d[2] + d[3]) / 4.0
}

func updateTopPt(pt *int, level *int) {
	*pt += 45
	updateLevel(pt, level)
}

func update2ndPt(pt *int, level *int) {
	*pt += 0
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



