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
	"encoding/json"

	"github.com/astaxie/beego/validation"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/article"
)

var GlobalValid *validation.Validation

func init() {
	GlobalValid = &validation.Validation{}
}

// ArticleController - article associated handler.
type ArticleController struct {
	baseController
}

// Create a new article
func (ac *ArticleController) Create() {
	var (
		articleInfo article.Article
		err         error
		articleID   string
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleInfo)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam, constants.RespKeyData: ac.Ctx.Input.RequestBody}

		goto finish
	}

	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleID, err = article.Service.Create(&articleInfo)
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
func (ac *ArticleController) List() {
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
func (ac *ArticleController) ActiveList() {
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

// GetArticleByTag get articles by tag
func (ac *ArticleController) GetArticleByTag() {
	var (
		articleList []article.Article
		tags        []string
		err         error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &tags)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleList, err = article.Service.GetByTags(tags)
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

// GetArticleById get by id
func (ac *ArticleController) GetArticleById() {
	var (
		articleRes article.Article
		articleId  string
		err        error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleId)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	articleRes, err = article.Service.GetByID(articleId)
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

// ModifyArticle modify article
func (ac *ArticleController) ModifyArticle() {
	var (
		articleToModify article.Article
		err             error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleToModify)
	if err != nil {
		logger.Error(err)
		ac.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = article.Service.Modify(&articleToModify)
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
