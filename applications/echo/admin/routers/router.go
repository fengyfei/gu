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

package routers

import (
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/handler/staff"
	"github.com/fengyfei/gu/applications/echo/core"
)

func InitRouter(server *echo.Echo, token string) {
	if server == nil {
		panic("[InitRouter]: server couldn't be nil")
	}

	core.InitJWTWithToken(token)

	// JWT middleware does not affect these route.
	core.URLMap["/api/v1/staff/login"] = struct{}{}
	core.URLMap["/api/v1/staff/create"] = struct{}{}

	// Staff
	server.POST("/api/v1/staff/login", staff.Login)
	server.POST("/api/v1/staff/create", staff.Create)
	server.POST("/api/v1/staff/modify/info", staff.Modify, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/api/v1/staff/modify/pwd", staff.ModifyPwd, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/api/v1/staff/modify/mobile", staff.ModifyMobile, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/api/v1/staff/modify/active", staff.ModifyActive, core.MustLoginIn)
	server.POST("/api/v1/staff/dismiss", staff.Dismiss, core.MustLoginIn, core.IsActiveMiddleWare)
	server.GET("/api/v1/staff/overview/list", staff.OverviewList, core.MustLoginIn, core.IsActiveMiddleWare)
	server.GET("/api/v1/staff/detail/list", staff.InfoList, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/api/v1/staff/detail", staff.Info, core.MustLoginIn, core.IsActiveMiddleWare)
}
