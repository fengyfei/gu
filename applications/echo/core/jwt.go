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

package core

import (
	"strconv"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	tokenExpireInHour = 48

	ClaimUID    = "uid"
	ClaimExpire = "exp"

	respTokenKey = "token"
)

var (
	TokenHMACKey string
	urlMap       map[string]struct{}
)

func CustomJWT(tokenkey string) echo.MiddlewareFunc {
	jwtconf := middleware.JWTConfig{
		Skipper:    CustomSkipper,
		SigningKey: []byte(tokenkey),
	}

	return middleware.JWTWithConfig(jwtconf)
}

func CustomSkipper(c echo.Context) bool {
	if _, ok := urlMap[c.Request().RequestURI]; ok {
		return true
	}

	return false
}

// InitJWTWithToken initialize the HMAC token.
func InitJWTWithToken(token string) {
	TokenHMACKey = token
}

// NewToken generates a JWT token.
func NewToken(uid int32) (string, string, error) {
	token := jwtgo.New(jwtgo.SigningMethodHS256)

	claims := token.Claims.(jwtgo.MapClaims)
	claims[ClaimUID] = strconv.FormatInt(int64(uid), 10)
	claims[ClaimExpire] = time.Now().Add(time.Hour * tokenExpireInHour).Unix()

	t, err := token.SignedString([]byte(TokenHMACKey))
	return respTokenKey, t, err
}

// UserID returns user identity.
func UserID(c echo.Context) int32 {
	uid := c.Get(ClaimUID).(int64)

	return int32(uid)
}

func init() {
	urlMap = make(map[string]struct{})

	urlMap["/staff/login"] = struct{}{}
	urlMap["/staff/create"] = struct{}{}
}
