package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type dao struct {
	db *sql.DB
}

func newDao() *dao {
	dao := new(dao)

	// just enjoying hard coding password
	db, err := sql.Open("mysql",
		"sakilogy:@k052a9@tcp(127.0.0.1:3306)/sakilogy")
	if err != nil {
		log.Fatalln(err)
	}

	if db.Ping() != nil {
		log.Fatalln("ping DB failed", err)
	}

	dao.db = db

	return dao
}

func (dao *dao) Close() {
	dao.db.Close()
}

func (dao *dao) Login(username, password string) (*ussn, error) {
	ussn := new(ussn)

	err := dao.db.QueryRow(
		`select user_id, username, level, pt, rating
		from users where username=? && password=?`,
		username, password).
		Scan(&ussn.user.Id, &ussn.user.Username, &ussn.user.Level,
			&ussn.user.Pt, &ussn.user.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		log.Fatalln("dao.login", err)
	}

	return ussn, nil
}

func (dao *dao) SignUp(username, password string) (*ussn, error) {
	var exist bool
	err := dao.db.QueryRow(
		"select exists(select 1 from users where username=?)",
		username).Scan(&exist)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	if exist {
		return nil, errors.New("用户名已存在")
	}

	_, err = dao.db.Exec(
		"insert into users (username, password) values (?,?)",
		username, password)

	if err != nil {
		log.Fatalln("dao.SignUp", err)
	}

	return dao.Login(username, password)
}

func (dao *dao) GetUser(uid uid) *user {
	user := new(user)

	err := dao.db.QueryRow(
		`select user_id, username, level, pt, rating
		from users where user_id=?`, uid).
		Scan(&user.Id, &user.Username, &user.Level,
			&user.Pt, &user.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalln("dao.GetUser", err)
	}

	return user
}

func (dao *dao) GetUsers(uids *[4]uid) [4]*user {
	var users [4]*user

	rows, err := dao.db.Query(
		`select user_id, username, level, pt, rating
		from users where user_id in (?,?,?,?)`,
		uids[0], uids[1], uids[2], uids[3])
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := new(user)
		err := rows.Scan(&user.Id, &user.Username,
			&user.Level, &user.Pt, &user.Rating)
		if err != nil {
			log.Fatalln(err)
		}
		for w := 0; w < 4; w++ {
			if uids[w] == user.Id {
				users[w] = user
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return users
}

func (dao *dao) GetRankedGids() []gid {
	var gids []gid

	// excluding doge
	rows, err := dao.db.Query(
		`select girl_id from girls where girl_id<>0 order by rating desc`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var gid gid
		err := rows.Scan(&gid)
		if err != nil {
			log.Fatalln(err)
		}
		gids = append(gids, gid)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return gids
}

func (dao *dao) GetStats(uid uid) []statRow {
	var stats []statRow

	// excluding doge
	rows, err := dao.db.Query(
		`select girl_id,rank1,rank2,rank3,rank4,
		avg_point,a_top,a_last,
		round,win,gun,bark,riichi,
		win_point,gun_point,bark_point,riichi_point,
		ready,ready_turn,win_turn,
		yaku_rci, yaku_ipt, yaku_tmo, yaku_tny, yaku_pnf,
		yaku_y1y, yaku_y2y, yaku_y3y,
		yaku_jk1, yaku_jk2, yaku_jk3, yaku_jk4,
		yaku_bk1, yaku_bk2, yaku_bk3, yaku_bk4,
		yaku_ipk, yaku_rns, yaku_hai, yaku_hou, yaku_ckn,
		yaku_ss1, yaku_it1, yaku_ct1,
		yaku_wri, yaku_ss2, yaku_it2, yaku_ct2,
		yaku_toi, yaku_ctt, yaku_sak, yaku_skt,
		yaku_stk, yaku_hrt, yaku_s3g, yaku_h1t, yaku_jc2,
		yaku_mnh, yaku_jc3, yaku_rpk, yaku_c1t, yaku_mnc,
		yaku_x13, yaku_xd3, yaku_x4a, yaku_xt1, yaku_xs4, yaku_xd4,
		yaku_xcr, yaku_xr1, yaku_xth, yaku_xch, yaku_x4k, yaku_x9r,
		yaku_w13, yaku_w4a, yaku_w9r, kzeykm,
		han_rci, han_ipt, han_tmo, han_tny, han_pnf,
		han_y1y, han_y2y, han_y3y,
		han_jk1, han_jk2, han_jk3, han_jk4,
		han_bk1, han_bk2, han_bk3, han_bk4,
		han_ipk, han_rns, han_hai, han_hou, han_ckn,
		han_ss1, han_it1, han_ct1,
		han_wri, han_ss2, han_it2, han_ct2,
		han_toi, han_ctt, han_sak, han_skt,
		han_stk, han_hrt, han_s3g, han_h1t, han_jc2,
		han_mnh, han_jc3, han_rpk, han_c1t, han_mnc,
		yaku_dora, yaku_uradora, yaku_akadora,
		yaku_kandora, yaku_kanuradora
		from user_girl where user_id=?`, uid)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r statRow
		err := rows.Scan(&r.GirlId,
			&r.Ranks[0], &r.Ranks[1], &r.Ranks[2], &r.Ranks[3],
			&r.AvgPoint, &r.ATop, &r.ALast,
			&r.Round, &r.Win, &r.Gun, &r.Bark, &r.Riichi,
			&r.WinPoint, &r.GunPoint, &r.BarkPoint, &r.RiichiPoint,
			&r.Ready, &r.ReadyTurn, &r.WinTurn,
			&r.Rci, &r.Ipt, &r.Tmo, &r.Tny, &r.Pnf,
			&r.Y1y, &r.Y2y, &r.Y3y,
			&r.Jk1, &r.Jk2, &r.Jk3, &r.Jk4,
			&r.Bk1, &r.Bk2, &r.Bk3, &r.Bk4,
			&r.Ipk, &r.Rns, &r.Hai, &r.Hou, &r.Ckn,
			&r.Ss1, &r.It1, &r.Ct1,
			&r.Wri, &r.Ss2, &r.It2, &r.Ct2,
			&r.Toi, &r.Ctt, &r.Sak, &r.Skt,
			&r.Stk, &r.Hrt, &r.S3g, &r.H1t, &r.Jc2,
			&r.Mnh, &r.Jc3, &r.Rpk, &r.C1t, &r.Mnc,
			&r.X13, &r.Xd3, &r.X4a, &r.Xt1, &r.Xs4, &r.Xd4,
			&r.Xcr, &r.Xr1, &r.Xth, &r.Xch, &r.X4k, &r.X9r,
			&r.W13, &r.W4a, &r.W9r, &r.Kzeykm,
			&r.RciHan, &r.IptHan, &r.TmoHan, &r.TnyHan, &r.PnfHan,
			&r.Y1yHan, &r.Y2yHan, &r.Y3yHan,
			&r.Jk1Han, &r.Jk2Han, &r.Jk3Han, &r.Jk4Han,
			&r.Bk1Han, &r.Bk2Han, &r.Bk3Han, &r.Bk4Han,
			&r.IpkHan, &r.RnsHan, &r.HaiHan, &r.HouHan, &r.CknHan,
			&r.Ss1Han, &r.It1Han, &r.Ct1Han,
			&r.WriHan, &r.Ss2Han, &r.It2Han, &r.Ct2Han,
			&r.ToiHan, &r.CttHan, &r.SakHan, &r.SktHan,
			&r.StkHan, &r.HrtHan, &r.S3gHan, &r.H1tHan, &r.Jc2Han,
			&r.MnhHan, &r.Jc3Han, &r.RpkHan, &r.C1tHan, &r.MncHan,
			&r.Dora, &r.Uradora, &r.Akadora, &r.Kandora, &r.Kanuradora)
		if err != nil {
			log.Fatalln(err)
		}
		stats = append(stats, r)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return stats
}

func (dao *dao) UpdateUserGirl(bt bookType, uids [4]uid, gids [4]gid,
	args *systemEndTableStat) {
	tx, err := dao.db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	err = updateUserGirlStat(tx, uids, gids, args)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	err = updateUserRank(tx, uids, args.Ranks, bt)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	err = updateGirlRank(tx, gids, args.Ranks, bt)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	tx.Commit()
}

func updateUserGirlStat(tx *sql.Tx, uids [4]uid, gids [4]gid,
	args *systemEndTableStat) error {
	for i := 0; i < 4; i++ {
		rankCol := "rank" + strconv.Itoa(args.Ranks[i])

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
		winPoint := float64(args.WinSumPoints[i])
		gunPoint := float64(args.GunSumPoints[i])
		barkPoint := float64(args.BarkSumPoints[i])
		riichiPoint := float64(args.RiichiSumPoints[i])
		winSumTurn := float64(args.WinSumTurns[i])
		var winAvg, winAvgTurn, gunAvg, barkAvg, riichiAvg float64
		if win != 0 {
			winAvg = winPoint / float64(win)
			winAvgTurn = winSumTurn / float64(win)
		}
		if gun != 0 {
			gunAvg = gunPoint / float64(gun)
		}
		if bark != 0 {
			barkAvg = barkPoint / float64(bark)
		}
		if riichi != 0 {
			riichiAvg = riichiPoint / float64(riichi)
		}

		ready := args.Readys[i]
		readySumTurn := float64(args.ReadySumTurns[i])
		var readyAvgTurn float64
		if ready != 0 {
			readyAvgTurn = readySumTurn / float64(ready)
		}

		kzeykm := args.Kzeykms[i]

		// fuck mariadb, cannot use virtual columns in "on dup key update"
		// ("play" will always be null somehow)
		// so manually typing (rank1+rank2+rank3+rank4) everywhere
		format := `insert into user_girl
			(user_id, girl_id, %s, avg_point, a_top, a_last,
				round, win, gun, bark, riichi,
				win_point, gun_point, bark_point, riichi_point,
				ready, ready_turn, win_turn, kzeykm)
			values (?, ?, 1, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?)
			on duplicate key update
			avg_point=(avg_point*(rank1+rank2+rank3+rank4)+?)
				/(rank1+rank2+rank3+rank4+1),
			a_top=a_top+?,a_last=a_last+?,
			win_point=if(win+?, (win_point*win+?)/(win+?), 0),
			gun_point=if(gun+?, (gun_point*gun+?)/(gun+?), 0),
			bark_point=if(bark+?, (bark_point*bark+?)/(bark+?), 0),
			riichi_point=if(riichi+?, (riichi_point*riichi+?)
				/(riichi+?), 0),
			win_turn=if(win+?, (win_turn*win+?)/(win+?), 0),
			ready_turn=if(ready+?, (ready_turn*ready+?)/(ready+?), 0),
			ready=ready+?,
			round=round+?,win=win+?,gun=gun+?,bark=bark+?,riichi=riichi+?,
			kzeykm=kzeykm+?,
			%s=%s+1`
		stmt := fmt.Sprintf(format, rankCol, rankCol, rankCol)
		_, err := tx.Exec(stmt,
			// "values" part
			uids[i], gids[i], args.Points[i], aTop, aLast,
			args.Round, win, gun, bark, riichi,
			winAvg, gunAvg, barkAvg, riichiAvg,
			ready, readyAvgTurn, winAvgTurn, kzeykm,
			// "update" part
			args.Points[i],
			aTop, aLast,
			win, winPoint, win,
			gun, gunPoint, gun,
			bark, barkPoint, bark,
			riichi, riichiPoint, riichi,
			win, winSumTurn, win,
			ready, readySumTurn, ready,
			ready,
			args.Round, win, gun, bark, riichi,
			kzeykm)
		if err != nil {
			return err
		}

		// update yaku stat
		yaku := args.Yakus[i]
		if len(yaku) > 0 {
			setter := ""
			han := args.SumHans[i]
			for key, ct := range yaku {
				hanCol := "han_" + key
				yakuCol := "yaku_" + key
				// update avg han unless yakuman
				if sum, ok := han[key]; ok {
					next := fmt.Sprintf("%s=(%s*%s+%d)/(%s+%d)",
						hanCol, hanCol, yakuCol, sum, yakuCol, ct)
					if setter == "" {
						setter = next
					} else {
						setter += ("," + next)
					}
				}

				next := fmt.Sprintf("%s=%s+%d", yakuCol, yakuCol, ct)
				if setter == "" {
					setter = next
				} else {
					setter += ("," + next)
				}
			}
			stmt := "update user_girl set " + setter +
				" where user_id=? and girl_id=?"
			_, err = tx.Exec(stmt, uids[i], gids[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateUserRank(tx *sql.Tx, uids [4]uid, ranks [4]int,
	bt bookType) error {
	var users [4]*user
	var plays [4]int

	rows, err := tx.Query(
		`select users.user_id, level, pt, rating, plays.play
		from users join
		(select user_id, sum(play) as play from user_girl group by user_id)
		as plays on users.user_id=plays.user_id
		where users.user_id in (?,?,?,?)`,
		uids[0], uids[1], uids[2], uids[3])
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		user := new(user)
		var play int
		err := rows.Scan(
			&user.Id, &user.Level, &user.Pt, &user.Rating, &play)
		if err != nil {
			return err
		}
		for w := 0; w < 4; w++ {
			if uids[w] == user.Id {
				users[w] = user
				plays[w] = play
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == nil {
			return errors.New("updateUserRank: nil user")
		}
	}

	var lprs [4]*lpr
	for i := 0; i < 4; i++ {
		lprs[i] = &users[i].lpr
	}

	updateLpr(&lprs, ranks, plays, bt)

	stmt := `update users set level=?, pt=?, rating=? where user_id=?`
	for w := 0; w < 4; w++ {
		u := users[w]
		_, err = tx.Exec(stmt, u.Level, u.Pt, u.Rating, u.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateGirlRank(tx *sql.Tx, gids [4]gid, ranks [4]int,
	bt bookType) error {
	var girls [4]*girl
	var plays [4]int

	rows, err := tx.Query(
		`select girls.girl_id, level, pt, rating, plays.play
		from girls join
		(select girl_id, sum(play) as play from user_girl group by girl_id)
		as plays on girls.girl_id=plays.girl_id
		where girls.girl_id in (?,?,?,?)`,
		gids[0], gids[1], gids[2], gids[3])
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		girl := new(girl)
		var play int
		err := rows.Scan(
			&girl.Id, &girl.Level, &girl.Pt, &girl.Rating, &play)
		if err != nil {
			return err
		}
		for w := 0; w < 4; w++ {
			if gids[w] == girl.Id {
				girls[w] = girl
				plays[w] = play
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	for _, girl := range girls {
		if girl == nil {
			return errors.New("updateGirlRank: nil girl")
		}
	}

	var lprs [4]*lpr
	for i := 0; i < 4; i++ {
		lprs[i] = &girls[i].lpr
	}

	updateLpr(&lprs, ranks, plays, bt)

	stmt := `update girls set level=?, pt=?, rating=? where girl_id=?`
	for w := 0; w < 4; w++ {
		g := girls[w]
		_, err = tx.Exec(stmt, g.Level, g.Pt, g.Rating, g.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
