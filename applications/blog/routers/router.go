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
	"github.com/fengyfei/gu/applications/blog/handler/staff"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/applications/blog/handler/project"
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

	u.Post("/blog/tag/create", article.CreateTag)
	u.Get("/blog/tag/list", article.ListTags)
	u.Get("/blog/tag/activelist", article.TagActiveList)
	u.Post("/blog/tag/info", article.TagInfo)
	u.Post("/blog/tag/modify", article.ModifyTag)

	u.Post("/staff/login", staff.Login)

	u.Post("/blog/article/create", article.CreateArticle)
	u.Post("/blog/article/allcreated", article.ListCreated)
	u.Post("/blog/article/delete", article.Delete)
	u.Post("/blog/article/modify", article.ModifyArticle)
	u.Post("/blog/article/modifystatus", article.ModifyStatus)
	u.Post("/blog/article/updateview", article.UpdateView)
	u.Post("/blog/article/approval", article.ListApproval)
	u.Post("/blog/article/getbyid", article.ArticleByID)
	u.Post("/blog/article/getbytag", article.GetByTag)

	u.Post("/blog/project/create", project.Create)
	u.Post("/blog/project/delete", project.Delete)
	u.Post("/blog/project/modify", project.Modify)
	u.Post("/blog/project/getid", project.GetID)
	u.Post("/blog/project/getbyid", project.GetByID)
	u.Get("/blog/project/list", project.AbstractList)
}
