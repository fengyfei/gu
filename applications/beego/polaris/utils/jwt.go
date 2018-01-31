/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/11/01        Jia Chenhui
 */

package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	beegoCtx "github.com/astaxie/beego/context"
	jwtgo "github.com/dgrijalva/jwt-go"
)

const (
	invalidUID = int32(-1)

	tokenExpireInHour = 48

	ClaimUID     = "uid"
	claimExpire  = "exp"
	respTokenKey = "token"
)

var (
	tokenHMACKey string
)

func init() {
	initJWTWithToken(beego.AppConfig.String("jwt::TokenKey"))
}

// initJWTWithToken initialize the HMAC token.
func initJWTWithToken(token string) {
	tokenHMACKey = token
}

// NewToken generates a JWT token.
func NewToken(uid int32) (string, string, error) {
	token := jwtgo.New(jwtgo.SigningMethodHS256)

	claims := token.Claims.(jwtgo.MapClaims)
	claims[ClaimUID] = uid
	claims[claimExpire] = time.Now().Add(time.Hour * tokenExpireInHour).Unix()

	t, err := token.SignedString([]byte(tokenHMACKey))
	return respTokenKey, t, err
}

// getUID check whether the token is valid, it returns user identity if valid.
func getUID(ctx *beegoCtx.Context) (int32, error) {
	claims, err := checkJWT(ctx)
	if err != nil {
		return invalidUID, err
	}

	rawUID := claims[ClaimUID].(float64)
	uid := int32(rawUID)

	return uid, nil
}

// bindUID bind the uid to the context.
func bindUID(ctx *beegoCtx.Context, uid int32) {
	newCtx := context.WithValue(context.Background(), ClaimUID, uid)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

// checkJWT check whether the token is valid, it returns claims if valid.
func checkJWT(ctx *beegoCtx.Context) (jwtgo.MapClaims, error) {
	var (
		err error
	)

	authString := ctx.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")

	if len(kv) != 2 || kv[0] != "Bearer" {
		err = errors.New("invalid token authorization string")
		return nil, err
	}

	tokenString := kv[1]

	token, _ := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenHMACKey), nil
	})

	claims, ok := token.Claims.(jwtgo.MapClaims)

	if !ok || !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}
