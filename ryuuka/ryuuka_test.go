package ryuuka

import (
	"fmt"
	"testing"
	"time"

	"github.com/rolevax/ih/ako/ss"
)

func TestSend(t *testing.T) {
	time.Sleep(1 * time.Second) // wait until connected

	msg := &ss.TableStart{
		Tid:  2333,
		Gids: []int64{0, 0, 0, 0},
	}
	resp, err := SendToToki(msg)
	if err != nil {
		fmt.Println("error", err)
	} else {
		fmt.Println(resp.(*ss.TableOutputs))
	}
}
