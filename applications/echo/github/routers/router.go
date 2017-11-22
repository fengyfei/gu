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
 *     Initial: 2017/11/17        Jia Chenhui
 */

package routers

import (
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/github/handler/article"
	"github.com/fengyfei/gu/applications/echo/github/handler/repos"
)

func InitRouter(server *echo.Echo) {
	if server == nil {
		panic("[InitRouter]: server couldn't be nil")
	}

	// Repos
	server.POST("/api/v1/repos/create", repos.Create)
	server.POST("/api/v1/repos/modify/active", repos.ModifyActive)
	server.GET("/api/v1/repos/list", repos.List)
	server.GET("/api/v1/repos/activelist", repos.ActiveList)
	server.POST("/api/v1/repos/info", repos.Info)

	// Article
	server.POST("/api/v1/article/create", article.Create)
	server.POST("/api/v1/article/modify/active", article.ModifyActive)
	server.GET("/api/v1/article/list", article.List)
	server.GET("/api/v1/article/activelist", article.ActiveList)
	server.POST("/api/v1/article/info", article.Info)
}
