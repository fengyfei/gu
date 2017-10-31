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
 *     Initial: 2017/10/27        ShiChao
 */

package controllers

import (
	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/article"
)

// Article - article associated handler.
type Article struct {
	base.Controller
}

// createArticleReq - the request struct that get article information by id.
type createArticleReq struct {
	Author   string   `json:"author" validate:"required"`
	Title    string   `json:"title" validate:"required,alphanumunicode,len=24"`
	Content  string   `json:"content" validate:"required"`
	Abstract string   `json:"abstract" validate:"required"`
	Tags      []string `json:"tag" validate:"required"`
	Active   bool     `json:"active" validate:"required"`
}

// getByTagReq - the request struct that get article information by tag.
type getByTagReq struct {
	Tags []string `json:"tags" validate:"required"`
}

// getByIdReq - the request struct that get article information by id.
type getByIdReq struct {
	ID string `json:"id" validate:"required"`
}

// modifyArticleReq - the request struct that modify article information by id.
type modifyArticleReq struct {
	ArticleID string `json:"id" validate:"required"`
	Title     string `json:"title" validate:"required,alphanumunicode,len=24"`
	Content   string `json:"content" validate:"required"`
	Abstract  string `json:"abstract" validate:"required"`
	Active    bool   `json:"active" validate:"required"`
}

// Create a new article
func (ac *Article) Create() {
	var (
		articleInfo createArticleReq
		err         error
		articleID   string
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleInfo)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam, constants.RespKeyData: ac.Ctx.Input.RequestBody}

		goto finish
	}

	err = ac.Validate(&articleInfo)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleID, err = article.Service.Create(&articleInfo.Author, &articleInfo.Title, &articleInfo.Abstract, &articleInfo.Content, &articleInfo.Tags)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: articleID}
	logger.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// List gets all article list
func (ac *Article) List() {
	var (
		articleList []article.Article
		err         error
	)

	articleList, err = article.Service.List()
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: articleList}
	logger.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// ActiveList gets active article list
func (ac *Article) ActiveList() {
	var (
		activeArticleList []article.Article
		err               error
	)

	activeArticleList, err = article.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: activeArticleList}
	logger.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// GetByTag get articles by tag
func (ac *Article) GetByTag() {
	var (
		articleList []article.Article
		tagArr      getByTagReq
		err         error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &tagArr)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleList, err = article.Service.GetByTags(&tagArr.Tags)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: articleList}
	logger.Info("get article by tags success")

finish:
	ac.ServeJSON(true)
}

// GetByID get by id
func (ac *Article) GetByID() {
	var (
		articleRes article.Article
		articleID  getByIdReq
		err        error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleID)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleRes, err = article.Service.GetByID(articleID.ID)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: articleRes}
	logger.Info("get article by id success")

finish:
	ac.ServeJSON(true)
}

// Modify modify article
func (ac *Article) Modify() {
	var (
		articleToModify modifyArticleReq
		err             error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleToModify)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = article.Service.Modify(&articleToModify.ArticleID, &articleToModify.Title, &articleToModify.Content, &articleToModify.Abstract, &articleToModify.Active)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}
	logger.Info("modify article by tags success")

finish:
	ac.ServeJSON(true)
}
