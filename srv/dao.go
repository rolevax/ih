package srv

import (
	"log"
	"database/sql"
	"errors"
	"strconv"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)

type dao struct {
	db		*sql.DB
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
		`select girl_id,rank1,rank2,rank3,rank4,avg_point
		from user_girl where user_id=?`, uid)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r statRow
		err := rows.Scan(&r.GirlId,
			&r.Ranks[0], &r.Ranks[1], &r.Ranks[2], &r.Ranks[3],
			&r.AvgPoint)
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
		// fuck mariadb, cannot use virtual columns in "on dup key update"
		// ("play" will always be null somehow)
		// so manually typing (rank1+rank2+rank3+rank4) everywhere
		format := `insert into user_girl
			(user_id, girl_id, %s, avg_point)
			values (?, ?, 1, ?)
			on duplicate key update
			avg_point=(avg_point*(rank1+rank2+rank3+rank4)+?)
				/(rank1+rank2+rank3+rank4+1),
			%s=%s+1`;
		stmt := fmt.Sprintf(format, rankCol, rankCol, rankCol)
		_, err := tx.Exec(stmt,
			uids[i], gids[i], args.Points[i],
			args.Points[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func updateUserRank(tx *sql.Tx, uids [4]uid, ranks [4]int, bt bookType) error {
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

func updateGirlRank(tx *sql.Tx, gids [4]gid, ranks [4]int, bt bookType) error {
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



