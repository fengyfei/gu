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
 *     Initial: 2017/10/22        Feng Yifei
 *     Modify : 2018/02/02        Tong Yuehong
 */

package router

import (
	"github.com/fengyfei/gu/applications/blog/handler"
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

	u.Post("/blog/tag/create", blog.CreateTag)
	u.Post("/blog/tag/list", blog.ListTags)
	u.Post("/blog/tag/activelist", blog.TagActiveList)
	u.Post("/blog/tag/info", blog.TagInfo)
	u.Post("/blog/tag/modify", blog.ModifyTag)

	u.Post("/blog/article/create", blog.CreateArticle)
	u.Post("/blog/article/list", blog.ListArticle)
	u.Post("/blog/article/activelist", blog.ArticleActiveList)
	u.Post("/blog/article/getbytag", blog.GetByTag)
	u.Post("/blog/article/getbyid", blog.GetByID)

	u.Post("/blog/article/modify", blog.ModifyArticle)
}
