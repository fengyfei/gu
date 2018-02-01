/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/11/15        ShiChao
 */

package util

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"

	libctx "context"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/fengyfei/gu/libs/constants"
)

const (
	tokenKey = "techcat_shop"
)

var (
	wechatLoginUrl = "/shop/user/wechatlogin"
	registerUrl    = "/shop/user/register"
	loginUrl       = "/shop/user/login"
)

func NewToken(userId uint, isAdmin bool) (string, error) {
	claims := make(jwt.MapClaims)
	claims["userid"] = userId
	claims["admin"] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 480).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenKey))
}

func Jwt(ctx *context.Context) {
	var (
		err         error
		token       *jwt.Token
		tokenString string
		claims      jwt.MapClaims
		ok          bool
		c           libctx.Context
		cc          libctx.Context
	)

	url := ctx.Request.URL.String()
	if url == wechatLoginUrl || url == loginUrl || url == registerUrl {
		return
	}

	authString := ctx.Input.Header("Authorization")
	//beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		goto errFinish
	}
	tokenString = kv[1]

	// Parse token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		goto errFinish
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		goto errFinish
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		goto errFinish
	}

	c = libctx.WithValue(libctx.Background(), "userId", int32(claims["userid"].(float64)))
	cc = libctx.WithValue(c, "isAdmin", bool(claims["admin"].(bool)))
	ctx.Request = ctx.Request.WithContext(cc)
	return

errFinish:
	ctx.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}, false, false)
}
