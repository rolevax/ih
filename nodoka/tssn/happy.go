package tssn

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rolevax/ih/ako/cs"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/ako/sc"
	"github.com/rolevax/ih/ako/ss"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/ryuuka"
)

func (tssn *tssn) Happy(ctx actor.Context) {
	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("hisa internal: %v", v)
			tssn.handleTokiCrash(err)
		}
	}()

	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		// somehow reset to 0 everytime this msg comes
		// diff from what doc says, may be their bug
		ctx.SetReceiveTimeout(recvTimeout)
		tssn.sweepAll()
	case *pcSeat:
		i, _ := tssn.findUser(msg.Uid)
		tssn.kick(i, "tssn.Happy get pcSeat")
	case *pcAction:
		tssn.handleAction(msg.Uid, msg.Act)
	case *ccAction:
		tssn.handleActionI(msg.UserIndex, msg.Act)
	default:
		log.Fatalf("tssn.Happy unexpected %T\n", msg)
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
	}
	tssn.waits[i] = false
	outputs, err := ryuuka.SendToToki(&ss.TableAction{
		Tid:     int64(tssn.match.Users[0].Id),
		Who:     int64(i),
		ActStr:  act.ActStr,
		ActArg:  int64(act.ActArg),
		ActTile: act.ActTile,
		Nonce:   int64(act.Nonce),
	})
	if err != nil {
		tssn.handleTokiCrash(err)
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
		tssn.handleTokiCrash(err)
		return
	}
	tssn.handleOutputs(outputs.(*ss.TableOutputs))
}

func (tssn *tssn) handleOutputs(to *ss.TableOutputs) {
	for _, mail := range to.Mails {
		toWhom := int(mail.Who)
		str := mail.Content
		var msg sc.TableEvent
		if err := json.Unmarshal([]byte(str), &msg); err != nil {
			log.Fatalln("unmarshal c++ str", err)
		}

		if toWhom == -1 {
			tssn.handleSystemMail(&msg)
		} else {
			tssn.sendUserMail(toWhom, &msg)
		}
	}

	if tssn.waitClient {
		tssn.waitClient = false
		time.Sleep(800 * time.Millisecond)
	}
}

func (tssn *tssn) sendUserMail(who int, msg *sc.TableEvent) {
	if tssn.match.Users[who].Id.IsBot() {
		log.Println("server-side bot feature deleted")
		return
	}

	tssn.injectResume(who, msg)

	err := tssn.sendPeer(who, msg)
	if msg.Event == "activated" {
		if err != nil {
			// clear off-liner's activation
			if tssn.anyOnline() {
				tssn.sweepOne(who)
			}
		} else {
			tssn.waits[who] = true
		}
	}

	// not a Hong Kong reporter, don't run so fast
	// wait for the client's rendering to avoid unintentional timeout
	if err == nil && msg.Event == "discarded" {
		tssn.waitClient = true
	}
}

func (tssn *tssn) handleSystemMail(msg *sc.TableEvent) {
	args := msg.Args

	switch msg.Event {
	case "game-over":
		tssn.gameOver = true
	case "round-start-log":
		tssn.handleRoundStartLog(msg.Args)
	case "table-end-stat":
		tssn.handleTableEndStat(msg)
	case "riichi-auto":
		tssn.handleRiichiAuto(msg.Args)
	case "action-expired":
		who := int(args["Who"].(float64))
		log.Printf(
			"TSSN EEEE %d action-expired by %d\n",
			tssn.match.Users[0].Id, tssn.match.Users[who].Id,
		)
	case "action-illegal":
		who := int(args["Who"].(float64))
		log.Printf(
			"TSSN EEEE %d action-illegal by %d\n",
			tssn.match.Users[0].Id, tssn.match.Users[who].Id,
		)
		tssn.kick(who, "illegal table action")
	case "table-tan90":
		tssn.handleTokiCrash(fmt.Errorf("table tan90"))
	default:
		log.Fatalln("unknown system mail", msg)
	}
}

func (tssn *tssn) handleRoundStartLog(args sc.VarMap) {
	al := ""
	if args["allLast"].(bool) {
		al = "a"
	}
	log.Printf(
		"TSSN .... %v %v.%v%s d=%v depo=%v seed=%v",
		tssn.match.Users[0].Id,
		args["round"], args["extra"], al,
		args["dealer"], args["deposit"], uint(args["seed"].(float64)),
	)
}

func (tssn *tssn) handleTableEndStat(msg *sc.TableEvent) {
	stat := &model.EndTableStat{}
	bytes, err := json.Marshal(msg.Args)
	if err != nil {
		log.Fatalln("table-end-stat marshal", err)
	}
	err = json.Unmarshal(bytes, stat)
	if err != nil {
		log.Fatalln("table-end-stat unmarshal", err)
	}
	tssn.injectReplay(stat.Replay)
	tssn.endTableAward(stat)
	err = mako.UpdateUserGirl(tssn.match.Uids(), stat)
	if err != nil {
		log.Println("table-end-stat db:", err)
	}
}

func (tssn *tssn) handleRiichiAuto(args sc.VarMap) {
	who := int(args["Who"].(float64))
	tssn.p.Tell(&ccAction{
		UserIndex: who,
		Act: &cs.TableAction{
			Nonce:  int(args["Nonce"].(float64)),
			ActStr: "SPIN_OUT",
		},
	})
}

func (tssn *tssn) injectReplay(replay map[string]interface{}) {
	var users [4]map[string]interface{}
	for w := 0; w < 4; w++ {
		user := make(map[string]interface{})
		user["Id"] = tssn.match.Users[w].Id
		user["Username"] = tssn.match.Users[w].Username
		users[w] = user
	}
	replay["users"] = users
}

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
		tssn.handleTokiCrash(err)
		return
	}
	tssn.handleOutputs(outputs.(*ss.TableOutputs))
}

func (tssn *tssn) sweepAll() {
	log.Println("TSSN swpa", tssn.match.Users[0].Id)

	outputs, err := ryuuka.SendToToki(&ss.TableSweepAll{
		Tid: int64(tssn.match.Users[0].Id),
	})
	if err != nil {
		tssn.handleTokiCrash(err)
		return
	}
	to := outputs.(*ss.TableOutputs)
	for w, waiting := range tssn.waits {
		if waiting {
			tssn.waits[w] = false
			tssn.kick(w, "happy timeout, sweep all")
		}
	}
	tssn.handleOutputs(to)
}

func (tssn *tssn) handleTokiCrash(err error) {
	log.Println("TSSN toki", tssn.match.Users[0].Id, err)

	for w := 0; w < 4; w++ {
		tssn.addFoodChange(w, &model.FoodChange{
			Delta:  32000,
			Reason: "发现服务器Bug奖励",
		})
	}
	tssn.p.Stop()
}

func (tssn *tssn) endTableAward(stat *model.EndTableStat) {
	for w := 0; w < 4; w++ {
		switch stat.Ranks[w] {
		case 1:
			tssn.addFoodChange(w, &model.FoodChange{
				Delta:  8000,
				Reason: "获得第一名",
			})
			if stat.ATop {
				tssn.addFoodChange(w, &model.FoodChange{
					Delta:  8000,
					Reason: "三杀",
				})
			}
		case 2:
			tssn.addFoodChange(w, &model.FoodChange{
				Delta:  1500,
				Reason: "获得第二名",
			})
		case 3:
			tssn.addFoodChange(w, &model.FoodChange{
				Delta:  1300,
				Reason: "获得第三名",
			})
		case 4:
			tssn.addFoodChange(w, &model.FoodChange{
				Delta:  1000,
				Reason: "获得第四名",
			})
			if stat.ALast {
				tssn.addFoodChange(w, &model.FoodChange{
					Delta:  1000,
					Reason: "骑马补贴",
				})
			}
		}
	}
}
