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
 *     Initial: 2018/01/31        Jia Chenhui
 */

package filter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TechCatsLab/apix/http/server"
	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
)

const (
	claimUID = "uid"
)

var (
	loginURI = "/api/v1/office/staff/login"
)

// LoginFilter check if the user is logged in.
func LoginFilter(c *server.Context) bool {
	if c.Request().RequestURI != loginURI {
		uid, _ := getUID(c)
		if uid == core.InvalidUID {
			core.WriteStatusAndDataJSON(c, constants.ErrPermission, nil)
			return false
		}
	}

	return true
}

// getUID check whether the token is valid, it returns user identity if valid.
func getUID(ctx *server.Context) (int32, error) {
	claims, err := checkJWT(ctx)
	if err != nil {
		return core.InvalidUID, err
	}

	rawUID := claims[claimUID].(float64)
	uid := int32(rawUID)

	return uid, nil
}

// checkJWT check whether the token is valid, it returns claims if valid.
func checkJWT(ctx *server.Context) (jwtgo.MapClaims, error) {
	var (
		err error
	)

	authString := ctx.Request().Header.Get("Authorization")
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

		return []byte(core.TokenHMACKey), nil
	})

	claims, ok := token.Claims.(jwtgo.MapClaims)

	if !ok || !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}
