package account

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/msg"
)

func PostAuth(request *restful.Request, response *restful.Response) {
	slow()

	sc := &msg.ScAuth{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAccountAuth{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	user, err := mako.Login(cs.Username, cs.Password)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  fmt.Sprint(user.Id),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
	})

	signed, err := token.SignedString([]byte(mako.GetAdminToken()))
	if err != nil {
		sc.Error = err.Error()
		return
	}

	sc.Jwt = signed
	sc.User = user
}

func PostCreate(request *restful.Request, response *restful.Response) {
	slow()

	sc := &msg.Sc{}
	defer response.WriteEntity(sc)

	cs := &msg.CsAccountCreate{}
	err := request.ReadEntity(cs)
	if err != nil {
		sc.Error = err.Error()
		return
	}

	err = mako.SignUp(cs.Username, cs.Password)
	if err != nil {
		sc.Error = err.Error()
	}
}

func GetCPoints(request *restful.Request, response *restful.Response) {
	sc := &msg.ScCpoints{}
	defer response.WriteEntity(sc)

	sc.Entries = mako.GetCPoints()
}

func slow() {
	time.Sleep(1 * time.Second)
}
