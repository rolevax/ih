package srv

import (
)

// user id
type uid uint

// girl id, signed-int for compatibility to libsaki
type gid int

// level, pt, and rating
type lpr struct {
	Level		int
	Pt			int
	Rating		float64
}

type user struct {
	Id			uid
	Username	string
	lpr
}

type girl struct {
	Id			gid
	lpr
}

type bookType int

func (b bookType) index() int {
	return int(b)
}

func (b bookType) valid() bool {
	i := int(b)
	return 0 <= i && i < 4
}

type reqTypeOnly struct {
	Type		string
}

type reqAuth struct {
	Type		string
	Version		string
	Username	string
	Password	string
}

type reqAction struct {
	Nonce		int
	ActStr		string
	ActArg		string
}

type reqBook struct {
	BookType	bookType
}

type reqChoose struct {
	GirlIndex	int
}

type respTypeOnly struct {
	Type		string
}

type respAuthFail struct {
	Type		string
	Ok			bool
	Reason		string
}

func newRespAuthFail(str string) interface{} {
	return respAuthFail{"auth", false, str}
}

type respAuthOk struct {
	Type	string
	Ok		bool
	User	*user
}

func newRespAuthOk(u *user) interface{} {
	return respAuthOk{"auth", true, u}
}

type bookEntry struct {
	Bookable	bool
	Book		int
	Play		int
}

type respLookAround struct {
	Type		string
	Conn		int
	Books		[4]bookEntry
}

func newRespLookAround(conn int) *respLookAround {
	resp := new(respLookAround)
	resp.Type = "look-around"
	resp.Conn = conn
	return resp
}

type statRow struct {
	GirlId			gid
	Ranks			[4]int
	AvgPoint		float64
	ATop			int
	ALast			int
	Round			int
	Win				int
	Gun				int
	Bark			int
	Riichi			int
	WinPoint		float64
	GunPoint		float64
	BarkPoint		float64
	RiichiPoint		float64
	Ready			int
	ReadyTurn		float64
	WinTurn			float64
	Rci				int
	Ipt				int
	Tmo				int
	Tny				int
	Pnf				int
	Y1y				int
	Y2y				int
	Y3y				int
	Jk1				int
	Jk2				int
	Jk3				int
	Jk4				int
	Bk1				int
	Bk2				int
	Bk3				int
	Bk4				int
	Ipk				int
	Rns				int
	Hai				int
	Hou				int
	Ckn				int
	Ss1				int
	It1				int
	Ct1				int
	Wri				int
	Ss2				int
	It2				int
	Ct2				int
	Toi				int
	Ctt				int
	Sak				int
	Skt				int
	Stk				int
	Hrt				int
	S3g				int
	H1t				int
	Jc2				int
	Mnh				int
	Jc3				int
	Rpk				int
	C1t				int
	Mnc				int
	X13				int
	Xd3				int
	X4a				int
	Xt1				int
	Xs4				int
	Xd4				int
	Xcr				int
	Xr1				int
	Xth				int
	Xch				int
	X4k				int
	X9r				int
	W13				int
	W4a				int
	W9r				int
	Kzeykm			int
	RciHan			float64
	IptHan			float64
	TmoHan			float64
	TnyHan			float64
	PnfHan			float64
	Y1yHan			float64
	Y2yHan			float64
	Y3yHan			float64
	Jk1Han			float64
	Jk2Han			float64
	Jk3Han			float64
	Jk4Han			float64
	Bk1Han			float64
	Bk2Han			float64
	Bk3Han			float64
	Bk4Han			float64
	IpkHan			float64
	RnsHan			float64
	HaiHan			float64
	HouHan			float64
	CknHan			float64
	Ss1Han			float64
	It1Han			float64
	Ct1Han			float64
	WriHan			float64
	Ss2Han			float64
	It2Han			float64
	Ct2Han			float64
	ToiHan			float64
	CttHan			float64
	SakHan			float64
	SktHan			float64
	StkHan			float64
	HrtHan			float64
	S3gHan			float64
	H1tHan			float64
	Jc2Han			float64
	MnhHan			float64
	Jc3Han			float64
	RpkHan			float64
	C1tHan			float64
	MncHan			float64
	Dora			int
	Uradora			int
	Akadora			int
	Kandora			int
	Kanuradora		int
}

type respUpdateUser struct {
	Type		string
	User		*user
	Stats		[]statRow
}

func newRespUpdateUser(user *user) *respUpdateUser {
	resp := new(respUpdateUser)
	resp.Type = "update-user"
	resp.User = user
	resp.Stats = sing.Dao.GetStats(user.Id)
	return resp
}



