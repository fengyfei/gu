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
 *     Initial: 2018/02/01        Jia Chenhui
 */

package core

import (
	"net/http"
	"time"

	"github.com/TechCatsLab/apix/http/server"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/libs/logger"
)

const (
	InvalidUID = int32(-1)

	tokenExpireInHour = 48

	claimUID     = "uid"
	claimExpire  = "exp"
	respTokenKey = "token"
	ClaimsKey    = "user"
)

var (
	// TokenHMACKey used to validate token.
	// Note: This variable must be initialized before attach JWT middleware
	// to HTTP Entrypoint.
	TokenHMACKey string

	// Defines a URL map to skip JWT validation on the specified route.
	URLMap map[string]struct{}
)

func init() {
	URLMap = make(map[string]struct{})
}

// InitHMACKey initialize TokenHMACKey.
func InitHMACKey(key string) {
	TokenHMACKey = key
}

// Skipper -
func Skipper(path string) bool {
	if _, ok := URLMap[path]; ok {
		return true
	}

	return false
}

// JwtErrHandler -
func JwtErrHandler(w http.ResponseWriter, r *http.Request, err string) {
	logger.Error(err)
	http.Error(w, err, 401)
}

// NewToken generates a JWT token based on uid.
func NewToken(uid int32) (string, string, error) {
	token := jwtgo.New(jwtgo.SigningMethodES256)

	claims := token.Claims.(jwtgo.MapClaims)
	claims[claimUID] = uid
	claims[claimExpire] = time.Now().Add(time.Hour * tokenExpireInHour).Unix()

	t, err := token.SignedString([]byte(TokenHMACKey))
	return respTokenKey, t, err
}

// GetUID extract uid from context.
func GetUID(c *server.Context) int32 {
	claims := c.Request().Context().Value(ClaimsKey)
	mapClaims := claims.(jwtgo.MapClaims)

	if rawUID, ok := mapClaims[claimUID].(float64); ok {
		return int32(rawUID)
	}

	return InvalidUID
}
