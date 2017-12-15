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

package main

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/polaris/auth"
	"github.com/fengyfei/gu/applications/beego/polaris/router"
)

func init() {
	auth.InitRPC()
	router.InitRouter()
	base.InitJWTWithToken(beego.AppConfig.String("jwt::TokenKey"))
}

// startServer starts an HTTP server.
func startServer() {
	go auth.RPCClients.Ping("AuthRPC.Ping")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     strings.Split(beego.AppConfig.String("cors::hosts"), ","),
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	beego.InsertFilter("/*", beego.BeforeRouter, base.LoginFilter)
	beego.InsertFilter("/*", beego.BeforeRouter, base.ActiveFilter)
	beego.InsertFilter("/*", beego.BeforeRouter, base.PermissionFilter)

	beego.Run()
}
