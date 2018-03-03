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
 *     Initial: 2018/02/01        Shi Ruitao
 */

package main

import (
	"github.com/fengyfei/gu/applications/shop/conf"
	"github.com/fengyfei/gu/applications/shop/routers"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/http/server/middleware"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/models/user"
)

var (
	ep *server.Entrypoint
)

// startServer starts a HTTP server.
func startServer() {
	serverConfig := &server.Configuration{
		Address: conf.ShopConfig.Address,
	}

	user.InitMongo("avatar", "user", conf.ShopConfig.MongoURL)

	mysql.InitPool()

	ep = server.NewEntrypoint(serverConfig, nil)

	// add middlewares
	ep.AttachMiddleware(middleware.NegroniRecoverHandler())
	ep.AttachMiddleware(middleware.NegroniLoggerHandler())
	ep.AttachMiddleware(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowedOrigins: conf.ShopConfig.CorsHosts,
		AllowedMethods: []string{server.GET, server.POST},
	}))

	if err := ep.Start(routers.Router.Handler()); err != nil {
		logger.Error(err)
		return
	}

	ep.Wait()
}
