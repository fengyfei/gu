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
 *     Initial: 2017/10/22        Feng Yifei
 */

package base

import (
	"github.com/astaxie/beego"
	"gopkg.in/go-playground/validator.v9"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
	"github.com/fengyfei/gu/libs/constants"
)

var(
	errInputData = errors.New("input error")
	errExpired = errors.New("expired")
	errToken = errors.New("token error")
)

// Controller wraps general functionality.
type Controller struct {
	beego.Controller
}

// Validate the parameters.
func (base Controller) Validate(val interface{}) error {
	validator := validator.New()

	return validator.Struct(val)
}

func (c *Controller) parseToken() (t *jwt.Token, e error) {
	authString := c.Ctx.Input.Header("Authorization")
	//beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil, errInputData
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("mykey"), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That‘s not even a token
				return nil, errInputData
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, errExpired
			} else {
				// Couldn‘t handle this token
				return nil, errInputData
			}
		} else {
			// Couldn‘t handle this token
			return nil, errInputData
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil, errInputData
	}
	beego.Debug("Token:", token)

	return token, nil
}

func (this *Controller) ParseToken() (int32, bool, error){
	token, err := this.parseToken()
	if err != nil {
		return 0, false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errToken
	}

	return int32(claims["userid"].(float64)), bool(claims["admin"].(bool)), nil
}