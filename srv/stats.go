package srv

import (
	"log"
)

func statRank(ordered *[4]uid) {
	users := sing.Dao.GetUsers(ordered)
	sumRating := 0.0
	for w := 0; w < 4; w++ {
		if users[w] == nil {
			log.Fatalln("uid", ordered[w], "not in DB")
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

	updateTopPt(users[0])
	update2ndPt(users[1])
	updateLastPt(users[3])

	for r := 0; r < 4; r++ {
		users[r].Ranks[r]++
	}

	sing.Dao.SetUsersRank(&users)
	for w := 0; w < 4; w++ {
		sing.UssnMgr.UpdateInfo(users[w])
	}
}

func average(d *[4]float64) float64 {
	return (d[0] + d[1] + d[2] + d[3]) / 4.0
}

func updateTopPt(user *user) {
	user.Pt += 45
	updateLevel(user)
}

func update2ndPt(user *user) {
	user.Pt += 0
	updateLevel(user)
}

func updateLastPt(user *user) {
	diffs := [20]int{
		0, 0, 0, 0, 0,
		0, 0, 0, -15, -30,
		-45, -60, -75, -90, -105,
		-120, -135, -150, -165, -180}
	user.Pt += diffs[user.Level]
	updateLevel(user)
}

func updateLevel(user *user) {
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

	if user.Pt >= maxs[user.Level] {
		if user.Level + 1 < len(maxs) {
			user.Level++
			user.Pt = starts[user.Level]
		} else {
			user.Pt = maxs[user.Level]
		}
	} else if user.Pt < 0 {
		if user.Level >= 10 {
			user.Level--
			user.Pt = starts[user.Level]
		} else {
			user.Pt = 0
		}
	}
}



