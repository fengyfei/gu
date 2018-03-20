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
 *     Initial: 2017/10/27        ShiChao
 *     Modify : 2018/02/02        Tong Yuehong
 */

package article

import (
	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/article"
)

// getByIdReq - the request struct that get article information by id.
type (
	getByIdReq struct {
		ID string `json:"id" validate:"required,alphanum,len=24"`
	}
)

// CreateArticle - insert article.
func CreateArticle(this *server.Context) error {
	var articleInfo article.CreateArticle

	if err := this.JSONBody(&articleInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articleInfo.AuthorID = int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	id, err := article.ArticleService.Create(articleInfo)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, id)
}

// ArticleByID return article by articleID.
func ArticleByID(this *server.Context) error {
	var (
		articleID getByIdReq
	)

	if err := this.JSONBody(&articleID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.GetByID(articleID.ID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, articles)
}

// ListCreated return articles which are waiting for checking.
func ListCreated(this *server.Context) error {
	articles, err := article.ArticleService.ListCreated()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, articles)
}

// ListApproval returns the articles which are passed.
func ListApproval(this *server.Context) error {
	var (
		page struct {
			Page int `json:"page"`
		}
	)

	if err := this.JSONBody(&page); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.ListApproval(page.Page)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, articles)
}

// ModifyStatus modify the article status.
func ModifyStatus(this *server.Context) error {
	var (
		req struct {
			ArticleID string `json:"articleID"`
			StaffID   int32  `json:"staffID"`
			Status    int8   `json:"status"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	req.StaffID = int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.ModifyStatus(req.ArticleID, req.Status, req.StaffID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Delete delete article.
func Delete(this *server.Context) error {
	var (
		req struct {
			ArticleID string `json:"articleID"`
			StaffID   int32  `json:"staffID"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	req.StaffID = int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.Delete(req.ArticleID, req.StaffID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)

}

// UpdateView update article's view.
func UpdateView(this *server.Context) error {
	var (
		view struct {
			ArticleID string `json:"articleID"`
			View      int32  `json:"view"`
		}
	)

	if err := this.JSONBody(&view); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.UpdateView(&view.ArticleID, view.View)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// ModifyArticle modify article.
func ModifyArticle(this *server.Context) error {
	var (
		modify struct {
			ArticleID string                `json:"articleID"`
			Article   article.CreateArticle `json:"article"`
		}
	)

	if err := this.JSONBody(&modify); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	modify.Article.AuthorID = int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.ModifyArticle(modify.ArticleID, modify.Article)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}
