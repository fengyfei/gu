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

	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/applications/echo/staff/handler"
)

func InitRouter(server *echo.Echo, token string) {
	if server == nil {
		panic("[InitRouter]: server couldn't be nil")
	}

	core.InitJWTWithToken(token)

	// Staff
	server.POST("/staff/login", handler.Login)
	server.POST("/staff/create", handler.Create, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/modify/pwd", handler.ModifyPwd, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/modify/mobile", handler.ModifyMobile, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/modify/active", handler.ModifyActive, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/staff/dismiss", handler.Dismiss, core.MustLoginIn, core.IsActiveMiddleWare)

	// Events
	server.GET("/events/checkin", handler.CheckIn, core.MustLoginIn, core.IsActiveMiddleWare)
	server.GET("/events/checkout", handler.CheckOut, core.MustLoginIn, core.IsActiveMiddleWare)

	// Contact
	server.GET("/contact/list", handler.ContactList, core.MustLoginIn, core.IsActiveMiddleWare)
	server.GET("/contact/list/detail", handler.InfoList, core.MustLoginIn, core.IsActiveMiddleWare)
	server.POST("/contact/detail", handler.Info, core.MustLoginIn, core.IsActiveMiddleWare)
}
