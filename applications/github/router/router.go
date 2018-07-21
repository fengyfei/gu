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
 *     Initial: 2017/12/28        Jia Chenhui
 */

package router

import (
	"github.com/TechCatsLab/apix/http/server"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/github/filter"
	"github.com/fengyfei/gu/applications/github/handler/article"
	"github.com/fengyfei/gu/applications/github/handler/issues"
	"github.com/fengyfei/gu/applications/github/handler/news"
	"github.com/fengyfei/gu/applications/github/handler/repos"
	"github.com/fengyfei/gu/applications/github/handler/trending"
)

var (
	Router *server.Router
)

func init() {
	Router = server.NewRouter()
	register(Router)
}

func register(r *server.Router) {

	// JWT middleware does not affect these route.
	core.URLMap["/api/v1/techcats/article/list"] = struct{}{}
	core.URLMap["/api/v1/techcats/article/activelist"] = struct{}{}
	core.URLMap["/api/v1/techcats/article/info"] = struct{}{}
	core.URLMap["/api/v1/techcats/repos/list"] = struct{}{}
	core.URLMap["/api/v1/techcats/repos/activelist"] = struct{}{}
	core.URLMap["/api/v1/techcats/repos/info"] = struct{}{}
	core.URLMap["/api/v1/techcats/repos/readme"] = struct{}{}
	core.URLMap["/api/v1/techcats/trending/lang"] = struct{}{}
	core.URLMap["/api/v1/techcats/issues/list"] = struct{}{}
	core.URLMap["/api/v1/techcats/issues/get"] = struct{}{}

	// Article
	r.Post("/api/v1/techcats/article/create", article.Create, filter.LoginFilter)
	r.Post("/api/v1/techcats/article/modify/active", article.ModifyActive, filter.LoginFilter)
	r.Get("/api/v1/techcats/article/list", article.List)
	r.Get("/api/v1/techcats/article/activelist", article.ActiveList)
	r.Post("/api/v1/techcats/article/info", article.Info)

	// News
	r.Post("/api/v1/techcats/news/info", news.Everything, filter.LoginFilter)

	// just for Wechat Mini Programs
	r.Get("/api/v1/techcats/repos/activelist", repos.ActiveListV1)
	// Repos New
	r.Post("/api/v1/techcats/repos/create", repos.Create, filter.LoginFilter)
	r.Post("/api/v1/techcats/repos/modify/active", repos.ModifyActive, filter.LoginFilter)
	r.Get("/api/v1/techcats/repos/list", repos.List)
	r.Post("/api/v1/techcats/repos/info", repos.Info)
	r.Post("/api/v1/techcats/repos/readme", repos.ReadmeURL)

	// Issues
	r.Post("/api/v1/techcats/issues/list", issues.List)
	r.Post("/api/v1/techcats/issues/get", issues.Get)

	// Trending
	r.Post("/api/v1/techcats/trending/lang", trending.LangInfo)
}
