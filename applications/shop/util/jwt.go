/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/02/02        Li Zebang
 */

package util

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/libs/http/server"
)

const (
	tokenKey = "techcat_shop"

	UserID     = "userid"
	SessionKey = "session_key"
	IsAdmin    = "is_admin"
)

var (
	wechatLoginUrl = "/shop/user/wechatlogin"
	registerUrl    = "/shop/user/register"
	loginUrl       = "/shop/user/login"
)

func NewToken(userID uint32, sessionKey string, isAdmin bool) (string, error) {
	claims := make(jwt.MapClaims)
	claims[UserID] = userID
	claims[SessionKey] = sessionKey
	claims[IsAdmin] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 480).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenKey))
}

func Parse(ctx *server.Context) (*jwt.Token, error) {
	var (
		tokenString string
		token       *jwt.Token
		err         error
	)

	url := ctx.Request().URL.String()
	if url == wechatLoginUrl || url == loginUrl || url == registerUrl {
		return nil, errors.New("don't use token on " + url)
	}

	authString := ctx.Request().Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return nil, errors.New("authString invalid " + authString)
	}

	tokenString = kv[1]
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalid " + tokenString)
	}

	return token, nil
}
