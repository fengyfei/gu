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
	"errors"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/models/staff"
)

var (
	errPermissionNotMatch = errors.New("user permissions and url permissions do not match")
	errNoRole             = errors.New("this user or url has no role")
)

// IsLogin check if the user is logged in.
func IsLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwtgo.Token)
		claims := user.Claims.(jwtgo.MapClaims)
		uid := claims[ClaimUID]

		c.Set(ClaimUID, uid)

		return next(c)
	}
}

// IsActive check whether the user is active.
func IsActive(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := mysql.Pool.Get()
		if err != nil {
			return NewErrorWithMsg(constants.ErrMysql, err.Error())
		}
		defer mysql.Pool.Release(conn)

		uid := UserID(c)
		ok, err := staff.Service.IsActive(conn, uid)

		if err != nil {
			return NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		if ok {
			return next(c)
		}

		return c.JSON(http.StatusConflict, err.Error())
	}
}

// IsPermissionMatch check whether the user permissions match the URL permissions.
func IsPermissionMatch(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := mysql.Pool.Get()
		if err != nil {
			return NewErrorWithMsg(constants.ErrMysql, err.Error())
		}
		defer mysql.Pool.Release(conn)

		url := c.Request().RequestURI
		urlRoles, err := staff.Service.URLPermissions(conn, &url)

		if err != nil {
			return NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		// If there is no permission record, pass the validation.
		if len(urlRoles) == 0 {
			return next(c)
		}

		uid := UserID(c)
		userRoles, err := staff.Service.AssociatedRoleMap(conn, uid)

		if err != nil {
			return NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		// If the user does not have a role, return the error directly.
		if len(userRoles) == 0 {
			err = errNoRole
			return NewErrorWithMsg(constants.ErrForbidden, err.Error())
		}

		for urlR := range urlRoles {
			if userRoles[urlR] {
				return next(c)
			}
		}

		err = errPermissionNotMatch

		return c.JSON(http.StatusOK, map[string]interface{}{
			constants.RespKeyStatus: constants.ErrForbidden,
		})
	}
}
