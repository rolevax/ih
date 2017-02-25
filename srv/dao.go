package srv

import (
	"log"
	"database/sql"
	"errors"
	_"github.com/go-sql-driver/mysql"
)

type dao struct {
	db		*sql.DB
}

func newDao() *dao {
	dao := new(dao)

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
		`select user_id, username, level, pt, rating,
		rank1, rank2, rank3, rank4
		from users where username=? && password=?`,
		username, password).
		Scan(&ussn.user.Id, &ussn.user.Username, &ussn.user.Level,
			 &ussn.user.Pt, &ussn.user.Rating,
			 &ussn.user.Ranks[0], &ussn.user.Ranks[1],
			 &ussn.user.Ranks[2], &ussn.user.Ranks[3])

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
		`select user_id, username, level, pt, rating,
		rank1, rank2, rank3, rank4
		from users where user_id=?`, uid).
		Scan(&user.Id, &user.Username, &user.Level, 
			 &user.Pt, &user.Rating,
			 &user.Ranks[0], &user.Ranks[1],
			 &user.Ranks[2], &user.Ranks[3])

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
		`select user_id, username, level, pt, rating,
		rank1, rank2, rank3, rank4
		from users where user_id in (?,?,?,?)`,
		uids[0], uids[1], uids[2], uids[3])
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := new(user)
		err := rows.Scan(&user.Id,
			&user.Username, &user.Level,
			&user.Pt, &user.Rating,
			&user.Ranks[0], &user.Ranks[1],
			&user.Ranks[2], &user.Ranks[3])
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

	rows, err := dao.db.Query(
		`select girl_id from girls order by rating desc`)
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

func (dao *dao) GetGirls(gids *[4]gid) [4]*girl {
	var girls [4]*girl

	rows, err := dao.db.Query(
		`select girl_id, level, pt, rating,
		rank1, rank2, rank3, rank4
		from girls where girl_id in (?,?,?,?)`,
		gids[0], gids[1], gids[2], gids[3])
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		girl := new(girl)
		err := rows.Scan(&girl.Id,
			&girl.Level, &girl.Pt, &girl.Rating,
			&girl.Ranks[0], &girl.Ranks[1],
			&girl.Ranks[2], &girl.Ranks[3])
		if err != nil {
			log.Fatalln(err)
		}
		for w := 0; w < 4; w++ {
			if gids[w] == girl.Id {
				girls[w] = girl
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	return girls
}

func (dao *dao) SetUsersRank(users *[4]*user) {
	for _, user := range users {
		if user == nil {
			log.Fatalln("dao.SetUsersRank: nil user")
		}
	}

	tx, err := dao.db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	stmt := `update users
		set level=?, pt=?, rating=?, rank1=?, rank2=?, rank3=?, rank4=?
		where user_id=?`
	for w := 0; w < 4; w++ {
		u := users[w]
		_, err = dao.db.Exec(stmt, u.Level, u.Pt, u.Rating,
			u.Ranks[0], u.Ranks[1], u.Ranks[2], u.Ranks[3],
			u.Id)

		if err != nil {
			log.Println("dao.SetUsersRank", err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
}

func (dao *dao) SetGirlsRank(girls *[4]*girl) {
	for _, girl := range girls {
		if girl == nil {
			log.Fatalln("dao.SetGirlRank: nil girl")
		}
	}

	tx, err := dao.db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	stmt := `update girls
		set level=?, pt=?, rating=?, rank1=?, rank2=?, rank3=?, rank4=?
		where girl_id=?`
	for w := 0; w < 4; w++ {
		g := girls[w]
		_, err = dao.db.Exec(stmt, g.Level, g.Pt, g.Rating,
			g.Ranks[0], g.Ranks[1], g.Ranks[2], g.Ranks[3],
			g.Id)

		if err != nil {
			log.Println("dao.SetGirlsRank", err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
}



