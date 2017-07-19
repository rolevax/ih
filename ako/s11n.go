package ako

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
)

type Decoder map[string]reflect.Type

func NewDecoder(msgs []interface{}) Decoder {
	d := Decoder{}

	for _, msg := range msgs {
		name := dash(reflect.TypeOf(msg).Name())
		d[name] = reflect.TypeOf(msg)
	}

	return d
}

func (d Decoder) FromJson(breq []byte) (interface{}, error) {
	tonly := struct{ Type string }{}
	if err := json.Unmarshal(breq, &tonly); err != nil {
		return nil, err
	}

	cs, err := d.interfaceOf(tonly.Type)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(breq, cs)
	if err != nil {
		return nil, err
	}

	return cs, err
}

func (d Decoder) interfaceOf(dashType string) (interface{}, error) {
	typ, ok := d[dashType]
	if !ok {
		return nil, fmt.Errorf("ako: no reg for %s", dashType)
	}
	return reflect.New(typ).Interface(), nil
}

func ToJson(msg interface{}) []byte {
	if reflect.TypeOf(msg).Kind() != reflect.Ptr {
		log.Fatalf("ako.ToJson want pointer get %T", msg)
	}

	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("ako.ToJson", err)
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(jsonb, &m)
	if err != nil {
		log.Fatal("ako.ToJson", err)
	}
	m["Type"] = dash(reflect.TypeOf(msg).Elem().Name())

	jsonb, err = json.Marshal(m)
	if err != nil {
		log.Fatal("ako.ToJson", err)
	}

	return jsonb
}

func dash(camel string) string {
	sp := camelcase.Split(camel)
	return strings.ToLower(strings.Join(sp, "-"))
}
