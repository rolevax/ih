package mako

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/rolevax/ih/ako/model"
)

func GetCultis(uid model.Uid) []model.Culti {
	var cs []model.Culti

	err := db.Model(&cs).Where("user_id=?", uid).Select()
	if err != nil {
		log.Fatalln("mako.GetStats", err)
	}

	return cs
}

func UpdateUserGirl(uids [4]model.Uid,
	gids [4]model.Gid, args *model.EndTableStat) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	err = updateUserGirlStat(tx, uids, gids, args)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	err = updateUserRank(tx, uids, args.Ranks)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	err = updateReplay(tx, uids, args.Replay)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	tx.Commit()
}

func updateUserGirlStat(tx *pg.Tx, uids [4]model.Uid,
	gids [4]model.Gid, args *model.EndTableStat) error {
	for i := 0; i < 4; i++ {
		var aTop, aLast int // fuck golang, cannot cast bool to int
		if args.ATop && args.Ranks[i] == 1 {
			aTop = 1
		}
		if args.ALast && args.Ranks[i] == 4 {
			aLast = 1
		}

		win := args.Wins[i]
		gun := args.Guns[i]
		bark := args.Barks[i]
		riichi := args.Riichis[i]
		ready := args.Readys[i]

		winPoint := float64(args.WinSumPoints[i])
		gunPoint := float64(args.GunSumPoints[i])
		barkP := float64(args.BarkSumPoints[i])
		riichiPoint := float64(args.RiichiSumPoints[i])
		winSumTurn := float64(args.WinSumTurns[i])
		readySumTurn := float64(args.ReadySumTurns[i])

		query := db.Model(&model.Culti{}).
			Set("avg_point=(avg_point*play(ranks)+?)/(play(ranks)+1)", args.Points[i]).
			Set("a_top=a_top+?", aTop).
			Set("a_last=a_last+?", aLast).
			Set("win_point=if(win+?, (win_point*win+?)/(win+?), 0)", win, winPoint, win).
			Set("gun_point=if(gun+?, (gun_point*gun+?)/(gun+?), 0)", gun, gunPoint, gun).
			Set("bark_point=if(bark+?, (bark_point*bark+?)/(bark+?), 0)", bark, barkP, bark).
			Set("riichi_point=if(riichi+?, (riichi_point*riichi+?)/(riichi+?), 0)", riichi, riichiPoint, riichi).
			Set("win_turn=if(win+?, (win_turn*win+?)/(win+?), 0)", win, winSumTurn, win).
			Set("ready_turn=if(ready+?, (ready_turn*ready+?)/(ready+?), 0)", ready, readySumTurn, ready).
			Set("ready=ready+?", ready).
			Set("round=round+?", args.Round).
			Set("win=win+?", win).
			Set("gun=gun+?", gun).
			Set("bark=bark+?", bark).
			Set("riichi=riichi+?", riichi).
			Set("kzeykm=kzeykm+?", args.Kzeykms[i]).
			// pg index starts from 1, so use rank value directly
			Set("ranks[?]=ranks[?]+1", args.Ranks[i], args.Ranks[i])

		if yaku := args.Yakus[i]; len(yaku) > 0 {
			han := args.SumHans[i]
			for key, ct := range yaku {
				hanCol := "han_" + key
				yakuCol := "yaku_" + key
				// update avg han unless yakuman
				if sum, ok := han[key]; ok {
					set := fmt.Sprintf(
						"%s=(%s*%s+?)/(%s+?)",
						hanCol, hanCol, yakuCol, yakuCol,
					)
					query = query.Set(set, sum, ct)
				}

				set := fmt.Sprintf("%s=%s+?", yakuCol, yakuCol)
				query = query.Set(set, ct)
			}
		}

		_, err := query.
			Where("user_id=?", uids[i]).
			Where("girl_id=?", gids[i]).
			Update()

		if err != nil {
			return err
		}
	}

	return nil
}

// using ostrich algorithm:
// - read-modify-write cycle
// - assume there is no race condition
func updateUserRank(tx *pg.Tx, uids [4]model.Uid, ranks [4]int) error {
	var res []struct {
		model.User
		Play int
	}

	_, err := tx.Query(
		&res,
		`SELECT users.user_id, level, pt, rating, plays.play
		FROM users JOIN (
		    SELECT user_id, SUM(play(ranks)) AS play
		    FROM user_girl GROUP BY user_id
		) AS plays ON users.user_id=plays.user_id
		WHERE users.user_id in (?)
		ORDER BY users.user_id=? DESC,
		         users.user_id=? DESC,
		         users.user_id=? DESC,
		         users.user_id=? DESC`,
		pg.In(uids), uids[0], uids[1], uids[2], uids[3],
	)
	if err != nil {
		return err
	}

	if len(res) != 4 {
		return fmt.Errorf("updateUserRank: res len %d", len(res))
	}

	var lprs [4]*model.Lpr
	var plays [4]int
	for i := 0; i < 4; i++ {
		lprs[i] = &res[i].Lpr
		plays[i] = res[i].Play
	}

	updateLpr(&lprs, ranks, plays)

	_, err = db.Model(&res[0].User, &res[1].User, &res[2].User, &res[3].User).
		Column("level").Column("pt").Column("rating").
		Update()

	return err
}
