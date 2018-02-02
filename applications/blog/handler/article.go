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
 *     Modify : 2018/02/02        Tong Yuehong
 */

package blog

import (
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/article"
)

// createArticleReq - the request struct that get article information by id.
type createArticleReq struct {
	Author   string   `json:"author"`
	Title    string   `json:"title" validate:"required,alphanumunicode,min=5,max=20"`
	Content  string   `json:"content" validate:"required,min=50"`
	Abstract string   `json:"abstract" validate:"required,min=10"`
	Tags     []string `json:"tag" validate:"required,dive,alphaunicode,min=2,max=6"`
	Active   bool     `json:"active" validate:"required"`
}

// getByTagReq - the request struct that get article information by tag.
type getByTagReq struct {
	Tags []string `json:"tags" validate:"required,dive,alphaunicode,min=2,max=6"`
}

// getByIdReq - the request struct that get article information by id.
type getByIdReq struct {
	ID string `json:"id" validate:"required,alphanum,len=24"`
}

// activateReq - the request struct that modify article information by id.
type activateReq struct {
	ArticleID string `json:"id" validate:"required,alphanum,len=24"`
	Title     string `json:"title" validate:"required,alphaunicode,min=6"`
	Content   string `json:"content" validate:"required,min=50"`
	Abstract  string `json:"abstract" validate:"required,min=10"`
	Active    bool   `json:"active" validate:"required"`
}

// Create - insert article.
func CreateArticle(this *server.Context) error {
	var articleInfo createArticleReq

	if err := this.JSONBody(&articleInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	id, err := article.Service.Create(&articleInfo.Author, &articleInfo.Title, &articleInfo.Abstract, &articleInfo.Content, &articleInfo.Tags)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, id)
}

// ListArticle return all article list
func ListArticle(this *server.Context) error {
	var articleList []article.Article

	if err := this.JSONBody(&articleList); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleList); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articleList, err := article.Service.List()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, articleList)
}

// ActiveList return active article list
func ArticleActiveList(this *server.Context) error {
	var activeArticleList []article.Article

	activeArticleList, err := article.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, activeArticleList)
}

// GetByTag return articles by tag.
func GetByTag(this *server.Context) error {
	var (
		articleList []article.Article
		tagArr      getByTagReq
	)

	if err := this.JSONBody(&tagArr); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleList); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articleList, err := article.Service.List()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articleList, err = article.Service.GetByTags(&tagArr.Tags)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, articleList)
}

// GetByID get by id
func GetByID(this *server.Context) error {
	var (
		articleRes article.Article
		articleID  getByIdReq
	)

	if err := this.JSONBody(&articleID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articleRes, err := article.Service.GetByID(articleID.ID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, articleRes)
}

// ModifyArticle modify article.
func ModifyArticle(this *server.Context) error {
	var articleToModify activateReq
	if err := this.JSONBody(&articleToModify); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&articleToModify); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.Service.Modify(&articleToModify.ArticleID, &articleToModify.Title, &articleToModify.Content, &articleToModify.Abstract, &articleToModify.Active)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}
