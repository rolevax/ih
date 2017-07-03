package cs

import (
	"encoding/json"
	"fmt"
)

func FromJson(breq []byte) (interface{}, error) {
	var tonly TypeOnly
	if err := json.Unmarshal(breq, &tonly); err != nil {
		return nil, err
	}

	switch tonly.Type {
	case "look-around":
		return &LookAround{}, nil
	case "heartbeat":
		return &HeartBeat{}, nil
	case "book":
		var cs Book
		err := json.Unmarshal(breq, &cs)
		if err == nil && !cs.BookType.Valid() {
			err = fmt.Errorf("invalid bktype %v", cs.BookType)
		}
		return &cs, err
	case "unbook":
		return &Unbook{}, nil
	case "get-replay-list":
		return &GetReplayList{}, nil
	case "get-replay":
		var cs GetReplay
		err := json.Unmarshal(breq, &cs)
		return &cs, err
	case "choose":
		var cs Choose
		err := json.Unmarshal(breq, &cs)
		if err == nil && !(0 <= cs.GirlIndex && cs.GirlIndex <= 2) {
			err = fmt.Errorf("invalid choose idx %d", cs.GirlIndex)
		}
		return &cs, err
	case "ready":
		return &Ready{}, nil
	case "t-action":
		var cs Action
		err := json.Unmarshal(breq, &cs)
		return &cs, err
	default:
		return nil, fmt.Errorf("model.From unknown type %s", tonly.Type)
	}
}
