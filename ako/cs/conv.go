package cs

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
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

func ToJson(cs interface{}) []byte {
	if reflect.TypeOf(cs).Kind() != reflect.Ptr {
		log.Fatalf("cs.ToJson want pointer get %T", cs)
	}

	jsonb, err := json.Marshal(cs)
	if err != nil {
		log.Fatal("cs.ToJson", err)
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(jsonb, &m)
	if err != nil {
		log.Fatal("cs.ToJson", err)
	}
	m["Type"] = dash(reflect.TypeOf(cs).Elem().Name())

	jsonb, err = json.Marshal(m)
	if err != nil {
		log.Fatal("cs.ToJson", err)
	}

	return jsonb
}

func dash(camel string) string {
	sp := camelcase.Split(camel)
	return strings.ToLower(strings.Join(sp, "-"))
}
