package model

type MatchState struct {
	Wait int
	Play int
}

type MatchResult struct {
	Id     int // id of the result set
	RuleId RuleId
	Users  [4]User
}

func (mr *MatchResult) Uids() [4]Uid {
	uids := [4]Uid{}
	for i, u := range mr.Users {
		uids[i] = u.Id
	}
	return uids
}

// rotate perspective
func (mr *MatchResult) RightPers() *MatchResult {
	next := &MatchResult{}

	*next = *mr

	user0 := mr.Users[0]
	for i := 0; i < 3; i++ {
		next.Users[i] = mr.Users[i+1]
	}
	next.Users[3] = user0

	return next
}
