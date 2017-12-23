package book

import (
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/nodoka"
)

type MatchState struct {
	Result model.MatchResult
	UserCt int
}

var (
	matchStates [model.RuleMax]MatchState
	matchIdGen  int
)

func init() {
	for i, _ := range matchStates {
		matchStates[i].Result.RuleId = model.RuleId(i)
	}
}

func (ms *MatchState) Add(msg *nodoka.MbMatchJoin) *model.MatchResult {
	res := (*model.MatchResult)(nil)

	for i := 0; i < ms.UserCt; i++ {
		if ms.Result.Users[i].Id == msg.User.Id {
			return nil
		}
	}

	ms.Result.Users[ms.UserCt] = msg.User
	ms.UserCt++

	if ms.UserCt == 4 {
		ms.UserCt = 0
		res = &model.MatchResult{}
		*res = ms.Result // copy
		res.Id = matchIdGen
		matchIdGen++
	}

	return res
}

func (ms *MatchState) RemoveUid(uid model.Uid) {
	for i := 0; i < ms.UserCt; i++ {
		if ms.Result.Users[i].Id == uid {
			ms.Result.Users[i] = ms.Result.Users[ms.UserCt-1]
			ms.UserCt--
			return
		}
	}
}
