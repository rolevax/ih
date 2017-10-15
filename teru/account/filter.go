package account

import (
	"fmt"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	restful "github.com/emicklei/go-restful"
	"github.com/rolevax/ih/ako/model"
	"github.com/rolevax/ih/mako"
	"github.com/rolevax/ih/teru/msg"
)

func FilterAuth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	replyError := func(err error) {
		sc := &msg.Sc{Error: err.Error()}
		resp.WriteEntity(sc)
	}

	auth := req.Request.Header.Get("Authorization")
	split := strings.Split(auth, " ")
	if len(split) != 2 {
		replyError(fmt.Errorf("jwt: wrong auth %s", auth))
		return
	}

	uid, err := getUidFromJwt(split[1])
	if err != nil {
		replyError(err)
		return
	}

	req.Request.Header.Set("X-User-Id", uid)

	chain.ProcessFilter(req, resp)
}

func getUidFromJwt(jwtStr string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt alg: %v", token.Header["alg"])
		}

		return []byte(mako.GetAdminToken()), nil
	}

	token, err := jwt.ParseWithClaims(jwtStr, &jwt.StandardClaims{}, keyFunc)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		i, err := strconv.Atoi(claims.Audience)
		if err != nil {
			return "", err
		} else if !model.Uid(uint(i)).IsHuman() {
			return "", fmt.Errorf("uid not a human")
		} else {
			return claims.Audience, nil
		}
	} else {
		return "", fmt.Errorf("invalid")
	}
}
