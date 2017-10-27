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
	"github.com/fengyfei/gu/models/blog/article"
	"github.com/astaxie/beego/validation"
	"github.com/fengyfei/gu/common"
	"github.com/astaxie/beego/logs"
)

var GlobalValid *validation.Validation

func init() {
	GlobalValid = &validation.Validation{}
}

type ArticleController struct {
	baseController
}

// add article
func (ac *ArticleController) AddArticle() {
	var (
		articleInfo article.MDCreateArticle
		err         error
		articleId   string
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleInfo)
	if err != nil {
		logs.Error("article unmarshal err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam, "data": ac.Ctx.Input.RequestBody}

		goto finish
	}

	if err != nil {
		logs.Error("article validation err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}

		goto finish
	}

	articleId, err = article.Service.Create(&articleInfo)
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed, common.RespKeyData: articleId}
	logs.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// get article list
func (ac *ArticleController) ListAll() {
	var (
		articleList []article.MDArticle
		err         error
	)

	articleList, err = article.Service.GetList()
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed, common.RespKeyData: articleList}
	logs.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// get active article list
func (ac *ArticleController) ActiveList() {
	var (
		activeArticleList []article.MDArticle
		err               error
	)

	activeArticleList, err = article.Service.GetActiveList()
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed, common.RespKeyData: activeArticleList}
	logs.Info("add article success")

finish:
	ac.ServeJSON(true)
}

// get articles by tag
func (ac *ArticleController) GetArticleByTag() {
	var (
		articleList []article.MDArticle
		tags        []string
		err         error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &tags)
	if err != nil {
		logs.Error("article unmarshal err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}

		goto finish
	}

	articleList, err = article.Service.GetByTags(tags)
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed, common.RespKeyData: articleList}
	logs.Info("get article by tags success")

finish:
	ac.ServeJSON(true)
}

// get by id
func (ac *ArticleController) GetArticleById() {
	var (
		articleRes article.MDArticle
		articleId  string
		err        error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleId)
	if err != nil {
		logs.Error("article unmarshal err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}

		goto finish
	}

	articleRes, err = article.Service.GetByID(articleId)
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed, common.RespKeyData: articleRes}
	logs.Info("get article by id success")

finish:
	ac.ServeJSON(true)
}

// modify article
func (ac *ArticleController) ModifyArticle() {
	var (
		articleToModify article.MDModifyArticle
		err             error
	)

	err = json.Unmarshal(ac.Ctx.Input.RequestBody, &articleToModify)
	if err != nil {
		logs.Error("article unmarshal err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}

		goto finish
	}

	err = article.Service.Modify(&articleToModify)
	if err != nil {
		logs.Error("article mongo err: ", err)
		ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}

		goto finish
	}
	ac.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed}
	logs.Info("modify article by tags success")

finish:
	ac.ServeJSON(true)
}
