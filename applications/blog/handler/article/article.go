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
 *     Modify : 2018/03/25        Chen Yanchen
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

// CreateArticle represents the article information when created.
type art struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Abstract string   `json:"abstract"`
	Tags     []string `json:"tags"`
	Image    string   `json:"image"`
}

// CreateArticle - insert article.
func CreateArticle(this *server.Context) error {
	var req art

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	AuthorID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	a := &article.Article{
		AuthorID: AuthorID,
		Title:    req.Title,
		Abstract: req.Abstract,
		Content:  req.Content,
		Tags:     req.Tags,
		Image:    req.Image,
	}

	id, err := article.ArticleService.Create(a)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, id)
}

// ArticleByID return article by articleID.
func ArticleByID(this *server.Context) error {
	var req struct {
		ID string `json:"aid" validate:"required,alphanum,len=24"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.GetByID(req.ID)
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
	var page struct {
		Page int `json:"page"`
	}

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
	var req struct {
		ArticleID string `json:"aid"`
		Status    int8   `json:"status"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	StaffID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.ModifyStatus(req.ArticleID, req.Status, StaffID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Delete delete article.
func Delete(this *server.Context) error {
	var req struct {
		ArticleID string `json:"aid"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	StaffID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.Delete(req.ArticleID, StaffID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// UpdateView update article's view.
func UpdateView(this *server.Context) error {
	var view struct {
		ArticleID string `json:"aid"`
		View      uint32 `json:"view"`
	}

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
	var req struct {
		ArticleID string   `json:"aid"`
		Title     string   `json:"title"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
		Tags      []string `json:"tags"`
		Image     string   `json:"image"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	AuditorID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	a := &article.Article{
		AuditorID: AuditorID,
		Title:     req.Title,
		Content:   req.Content,
		Abstract:  req.Abstract,
		Tags:      req.Tags,
		Image:     req.Image,
	}

	err := article.ArticleService.ModifyArticle(req.ArticleID, *a)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// GetByTag get article by tag.
func GetByTag(c *server.Context) error {
	var req struct {
		Tag string `json:"tag"`
	}
	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err := article.ArticleService.GetByTag(req.Tag)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}
