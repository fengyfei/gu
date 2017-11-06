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
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/models/staff"
)

func MustLoginIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwtgo.Token)
		claims := user.Claims.(jwtgo.MapClaims)
		uid := claims[ClaimUID]

		c.Set(ClaimUID, uid)

		return next(c)
	}
}

func IsActiveMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := mysql.Pool.Get()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer mysql.Pool.Release(conn)

		uid := UserID(c)
		ok, err := staff.Service.IsActive(conn, &uid)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if ok {
			return next(c)
		}

		return c.JSON(http.StatusConflict, err.Error())
	}
}
