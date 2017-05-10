package model

import (
	"encoding/json"
	"fmt"
)

func FromJson(breq []byte) (interface{}, error) {
	var tonly CsTypeOnly
	if err := json.Unmarshal(breq, &tonly); err != nil {
		return nil, err
	}

	switch tonly.Type {
	case "look-around":
		return &CsLookAround{}, nil
	case "heartbeat":
		return &CsHeartBeat{}, nil
	case "book":
		var cs CsBook
		err := json.Unmarshal(breq, &cs)
		if err == nil && !cs.BookType.Valid() {
			err = fmt.Errorf("invalid bktype %v", cs.BookType)
		}
		return &cs, err
	case "unbook":
		return &CsUnbook{}, nil
	case "get-replay-list":
		return &CsGetReplayList{}, nil
	case "get-replay":
		var cs CsGetReplay
		err := json.Unmarshal(breq, &cs)
		return &cs, err
	case "choose":
		var cs CsChoose
		err := json.Unmarshal(breq, &cs)
		if err == nil && !(0 <= cs.GirlIndex && cs.GirlIndex <= 2) {
			err = fmt.Errorf("invalid choose idx %d", cs.GirlIndex)
		}
		return &cs, err
	case "ready":
		return &CsReady{}, nil
	case "t-action":
		var cs CsAction
		err := json.Unmarshal(breq, &cs)
		return &cs, err
	default:
		return nil, fmt.Errorf("model.From unknown type %s", tonly.Type)
	}
}
