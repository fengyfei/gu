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
 *     Initial: 2018/01/21        Chen Yanchen
 */

package user

import (
	"github.com/fengyfei/gu/libs/http/server"

	"github.com/fengyfei/gu/applications/bbs/handler/article"
	"github.com/fengyfei/gu/applications/bbs/handler/user"
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

	u.Post("/user/wechatlogin", user.WechatLogin)
	u.Post("/user/changeusername", user.ChangeUsername)
	u.Post("/user/changeavatar", user.ChangeAvatar)

	u.Post("/article/insert", article.AddArticle)
	u.Post("/article/delete", article.DeleteArt)
	u.Post("/article/module", article.GetByModuleID)
	u.Post("/article/theme", article.GetByThemeID)
	u.Post("/article/title", article.GetByTitle)
	u.Post("/article/userid", article.GetByUserId)
	u.Post("/article/updatetimes", article.UpdateTimes)

	u.Post("/module/insert", article.AddModule)
	u.Post("/module/theme/insert", article.AddTheme)
	u.Post("/module/delete", article.DeleteModule)
	u.Post("/module/theme/delete", article.DeleteTheme)
	u.Post("/module/getinfo", article.GetInfo)
	u.Get("/module/allmodule", article.GetAllModule)
	u.Post("/module/updateview", article.UpdateModuleView)

	u.Post("/comment/insert", article.AddComment)
	u.Post("/comment/delete", article.DeleteComment)
	u.Get("/comment/getinfo", article.GetCommentInfo)
}
