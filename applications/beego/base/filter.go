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

package base

import (
	"github.com/astaxie/beego/context"

	"github.com/fengyfei/gu/applications/beego/polaris/auth"
	"github.com/fengyfei/gu/applications/beego/polaris/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/permission"
	"github.com/fengyfei/gu/models/staff"
)

var (
	loginURI string = "/api/v1/office/staff/login"
)

// LoginFilter check if the user is logged in.
func LoginFilter(c *context.Context) {
	if c.Request.RequestURI != loginURI {
		uid, err := UserID(c)
		if err != nil || uid == nil {
			logger.Error(err)
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}, false, false)
		}
	}
}

// ActiveFilter check whether the user is active.
func ActiveFilter(c *context.Context) {
	if c.Request.RequestURI != loginURI {
		conn, err := mysql.Pool.Get()
		if err != nil {
			logger.Error(err)
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}, false, false)
		}

		uid, err := UserID(c)
		if err != nil {
			logger.Error(err)
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}, false, false)
		}

		ok, err := staff.Service.IsActive(conn, *uid)
		if err != nil {
			logger.Error(err)
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}, false, false)
		}

		if !ok {
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}, false, false)
		}
	}
}

// IsPermissionMatch check whether the user permissions match the URL permissions.
func PermissionFilter(c *context.Context) {
	var (
		ok bool
	)

	if c.Request.RequestURI != loginURI {
		uid, _ := UserID(c)
		if uid == nil {
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}, false, false)
		}

		args := permission.Args{
			URL: c.Request.RequestURI,
			UId: *uid,
		}

		rpcClient, err := auth.RPCClients.Get(auth.RPCAddress)
		if err != nil || rpcClient == nil {
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrForbidden}, false, false)
		}

		err = rpcClient.Call("AuthRPC.Verify", &args, &ok)
		if err != nil || !ok {
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrForbidden}, false, false)
		}
	}
}

// JWTFilter check whether the token is valid,
// and if the token is valid, bind it to the context.
func JWTFilter(c *context.Context) {
	if c.Request.RequestURI != loginURI {
		uid, err := UserID(c)
		if err != nil {
			c.Output.JSON(map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}, false, false)
		}

		bindUID(c, *uid)
	}
}
