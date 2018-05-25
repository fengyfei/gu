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
 *     Initial: 2017/10/22        Feng Yifei
 *     Modify : 2018/02/02        Tong Yuehong
 *     Modify : 2018/03/25        Chen Yanchen
 */

package router

import (
	"github.com/fengyfei/gu/applications/blog/handler/article"
	"github.com/fengyfei/gu/applications/blog/handler/project"
	"github.com/fengyfei/gu/applications/blog/handler/staff"
	"github.com/fengyfei/gu/libs/http/server"
)

var (
	Router *server.Router
)

func init() {
	Router = server.NewRouter()
	InitRouter(Router)
}

func InitRouter(u *server.Router) {
	if u == nil {
		panic("[InitRouter]: server couldn't be nil")
	}

	u.Post("/api/v1/blog/tag/create", article.CreateTag)
	u.Get("/api/v1/blog/tag/list", article.ListTags)
	u.Get("/api/v1/blog/tag/activelist", article.TagActiveList)
	u.Post("/api/v1/blog/tag/info", article.TagInfo)
	u.Post("/api/v1/blog/tag/modify", article.ModifyTag)

	u.Post("/api/v1/staff/login", staff.Login)
	u.Post("/api/v1/staff/info", staff.Info)

	u.Post("/api/v1/blog/article/create", article.CreateArticle)
	u.Get("/api/v1/blog/article/allcreated", article.ListCreated)
	u.Post("/api/v1/blog/article/delete", article.Delete)
	u.Post("/api/v1/blog/article/modify", article.ModifyArticle)
	u.Post("/api/v1/blog/article/modifystatus", article.ModifyStatus)
	u.Post("/api/v1/blog/article/updateview", article.UpdateView)
	u.Post("/api/v1/blog/article/approval", article.ListApproval)
	u.Post("/api/v1/blog/article/getbyid", article.ArticleByID)
	u.Post("/api/v1/blog/article/getbyauthorid", article.GetByAuthorID)
	u.Post("/api/v1/blog/article/getbytag", article.GetByTag)

	u.Post("/api/v1/blog/project/create", project.Create)
	u.Post("/api/v1/blog/project/delete", project.Delete)
	u.Post("/api/v1/blog/project/modify", project.Modify)
	u.Post("/api/v1/blog/project/getid", project.GetID)
	u.Post("/api/v1/blog/project/getbyid", project.GetByID)
	u.Get("/api/v1/blog/project/list", project.AbstractList)
}
