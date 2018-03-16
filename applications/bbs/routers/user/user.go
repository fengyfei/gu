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

	u.Post("/bbs/user/wechatlogin", user.WechatLogin)
	u.Post("/bbs/user/addphone", user.AddPhone)
	u.Post("/bbs/user/changeinfo", user.ChangeUserInfo)
	u.Post("/bbs/user/changepassword", user.ChangePassword)
	u.Post("/bbs/user/phoneregister", user.PhoneRegister)
	u.Post("/bbs/user/phonelogin", user.PhoneLogin)
	u.Get("/bbs/user/homepage", user.BbsUserInfo)

	u.Post("/bbs/article/insert", article.AddArticle)
	u.Post("/bbs/article/delete", article.DeleteArt)
	u.Post("/bbs/article/module", article.GetByModuleID)
	u.Post("/bbs/article/theme", article.GetByThemeID)
	u.Post("/bbs/article/title", article.GetByTitle)
	u.Post("/bbs/article/userid", article.GetByUserID)
	u.Post("/bbs/article/updatetimes", article.UpdateTimes)
	u.Post("/bbs/article/recommend", article.Recommend)

	u.Post("/bbs/module/insert", article.AddModule)
	u.Post("/bbs/module/theme/insert", article.AddTheme)
	u.Post("/bbs/module/delete", article.DeleteModule)
	u.Post("/bbs/module/theme/delete", article.DeleteTheme)
	u.Post("/bbs/module/oneinfo", article.ModuleInfo)
	u.Get("/bbs/module/allinfo", article.AllModules)
	u.Get("/bbs/module/recommend", article.ListRecommend)
	u.Post("/bbs/module/updateview", article.UpdateModuleView)

	u.Post("/bbs/comment/insert", article.AddComment)
	u.Post("/bbs/comment/delete", article.DeleteComment)
	u.Post("/bbs/comment/listinfo", article.CommentInfo)
	u.Post("/bbs/comment/userreply", article.UserReply)
	u.Post("/bbs/comment/article", article.GetByArticle)
}
