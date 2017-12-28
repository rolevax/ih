package tssn

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
	"github.com/rolevax/ih/ako/ss"
	"github.com/rolevax/ih/nodoka"
	"github.com/rolevax/ih/ryuuka"
)

func (tssn *tssn) Happy(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		tssn.sweepAll()
	case *pcSeat:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Happy get pcSeat")
	case *pcAction:
		tssn.handleAction(msg.Uid, msg.Act)
	case *ccAction:
		tssn.handleActionI(msg.UserIndex, msg.Act)
	default:
		log.Fatalf("tssn.Seat unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) handleAction(uid model.Uid, act *cs.TableAction) {
	i, ok := tssn.findUser(uid)
	if !ok {
		log.Fatalf("tssn.handleAction user %d not found\n", uid)
	}
	tssn.handleActionI(i, act)
}

func (tssn *tssn) handleActionI(i int, act *cs.TableAction) {
	if act.ActStr == "RESUME" {
		tssn.onlines[i] = true
	} else if act.Nonce != tssn.nonces[i] {
		tssn.kick(i, "wrong action nonce")
		return
	}
	tssn.waits[i] = false
	outputs, err := ryuuka.SendToToki(&ss.TableAction{
		Tid:     int64(tssn.match.Users[0].Id),
		Who:     int64(i),
		ActStr:  act.ActStr,
		ActArg:  int64(act.ActArg),
		ActTile: act.ActTile,
	})
	if err != nil {
		// FUCK
		return
	}
	tssn.handleOutputs(outputs.(*ss.TableOutputs))
}

func (tssn *tssn) start() {
	gids := &tssn.gids
	log.Println("TSSN ****", gids)

	msg := &ss.TableStart{Tid: int64(tssn.match.Users[0].Id)}
	for _, gid := range tssn.gids {
		msg.Gids = append(msg.Gids, int64(gid))
	}

	outputs, err := ryuuka.SendToToki(msg)
	if err != nil {
		// FUCK
		return
	}
	tssn.handleOutputs(outputs.(*ss.TableOutputs))
}

func (tssn *tssn) handleOutputs(to *ss.TableOutputs) {
	if to.GameOver {
		tssn.gameOver = true
	}

	if len(to.Mails) > 0 {
		var nonceInced [4]bool
		for _, mail := range to.Mails {
			w := int(mail.Who)
			if w != -1 && !nonceInced[w] {
				tssn.nonces[w]++
				nonceInced[w] = true
			}
		}
	}

	for _, mail := range to.Mails {
		toWhom := int(mail.Who)
		str := mail.Content
		if toWhom == -1 {
			var msg map[string]interface{}
			if err := json.Unmarshal([]byte(str), &msg); err != nil {
				log.Fatalln("unmarshal c++ str", err)
			}
			tssn.handleSystemMail(msg, str)
		} else {
			var msg sc.TableEvent
			if err := json.Unmarshal([]byte(str), &msg); err != nil {
				log.Fatalln("unmarshal c++ str", err)
			}
			tssn.sendUserMail(toWhom, &msg)
		}
	}

	if tssn.waitClient {
		tssn.waitClient = false
		time.Sleep(800 * time.Millisecond)
	}
}

func (tssn *tssn) sendUserMail(who int, msg *sc.TableEvent) {
	msg.Nonce = tssn.nonces[who]

	if tssn.match.Users[who].Id.IsBot() {
		if msg.Event == "activated" {
			go func() {
				// simulate ai thinking time
				hesi := 300 * time.Millisecond
				actMap := msg.Args["action"].(map[string]interface{})
				if _, ok := actMap["PASS"]; ok {
					hesi = 100 * time.Millisecond
				}
				time.Sleep(hesi)

				tssn.p.Tell(&ccAction{
					UserIndex: who,
					Act: &cs.TableAction{
						ActStr: "BOT",
						Nonce:  msg.Nonce,
					},
				})
			}()
		}
		return
	}

	tssn.injectResume(who, msg)

	err := tssn.sendPeer(who, msg)
	if err != nil && msg.Event == "activated" {
		if tssn.anyOnline() {
			tssn.sweepOne(who)
		}
	}

	// not a Hong Kong reporter, don't run so fast
	// wait for the client's rendering to avoid unintentional timeout
	if err == nil && msg.Event == "discarded" {
		tssn.waitClient = true
	}
}

func (tssn *tssn) handleSystemMail(msg map[string]interface{},
	msgStr string) {
	switch msg["Type"] {
	case "round-start-log":
		al := ""
		if msg["allLast"].(bool) {
			al = "a"
		}
		log.Printf(
			"TSSN .... %v %v.%v%s d=%v depo=%v seed=%v",
			tssn.match.Users[0].Id,
			msg["round"], msg["extra"], al,
			msg["dealer"], msg["deposit"], uint(msg["seed"].(float64)),
		)
	case "table-end-stat":
		var stat model.EndTableStat
		err := json.Unmarshal([]byte(msgStr), &stat)
		if err != nil {
			log.Fatalln("table-end-stat unmarshal", err)
		}
		//tssn.injectReplay(stat.Replay)
		//TODO
		//mako.UpdateUserGirl(tssn.abcd, tssn.uids, tssn.gids, &stat)
		for w := 0; w < 4; w++ {
			nodoka.Umgr.Tell(&nodoka.MuUpdateInfo{Uid: tssn.match.Users[w].Id})
		}
	case "riichi-auto":
		who := int(msg["Who"].(float64))
		tssn.p.Tell(&ccAction{
			UserIndex: who,
			Act: &cs.TableAction{
				Nonce:  tssn.nonces[who],
				ActStr: "SPIN_OUT",
			},
		})
	case "cannot":
		who := int(msg["who"].(float64))
		actStr := msg["actStr"].(string)
		actArg := msg["actArg"].(string)
		log.Printf("TSSN EEEE %d cannot %d-%s-%s\n",
			tssn.match.Id, tssn.match.Users[who].Id, actStr, actArg)
		tssn.kick(who, "illegal table action")
	default:
		log.Fatalln("unknown system mail", msg)
	}
}

/* TODO merge into mako.EndTable(...)
func (tssn *tssn) injectReplay(replay map[string]interface{}) {
	var users [4]map[string]interface{}
	for w := 0; w < 4; w++ {
		user := make(map[string]interface{})
		user["Id"] = tssn.uids[w]
		user["Username"] = tssn.users[w].Username
		user["Level"] = tssn.users[w].Level
		user["Rating"] = tssn.users[w].Rating
		users[w] = user
	}
	replay["users"] = users
}
*/

func (tssn *tssn) injectResume(who int, msg *sc.TableEvent) {
	if msg.Event == "resume" {
		right := (who + 1) % 4
		cross := (who + 2) % 4
		left := (who + 3) % 4
		rotated := [4]*model.User{
			&tssn.match.Users[who],
			&tssn.match.Users[right],
			&tssn.match.Users[cross],
			&tssn.match.Users[left],
		}
		msg.Args["users"] = rotated
	}
}

func (tssn *tssn) sweepOne(i int) {
	outputs, err := ryuuka.SendToToki(&ss.TableSweepOne{
		Tid: int64(tssn.match.Users[0].Id),
		Who: int64(i),
	})
	if err != nil {
		// FUCK
		return
	}
	tssn.handleOutputs(outputs.(*ss.TableOutputs))
}

func (tssn *tssn) sweepAll() {
	outputs, err := ryuuka.SendToToki(&ss.TableSweepAll{
		Tid: int64(tssn.match.Users[0].Id),
	})
	if err != nil {
		// FUCK
		return
	}
	to := outputs.(*ss.TableOutputs)
	for _, who := range to.Sweepees {
		w := int(who)
		tssn.waits[w] = false
		tssn.kick(w, "happy timeout")
	}
	tssn.handleOutputs(to)
}
