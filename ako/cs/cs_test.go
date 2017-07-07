package cs

import (
	"fmt"
	"testing"
)

func TestToJson(t *testing.T) {
	cs := &LookAround{}
	json := string(ToJson(cs))
	fmt.Println(string(json))
}
