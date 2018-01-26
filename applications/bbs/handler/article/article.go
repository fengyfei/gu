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
 *     Initial: 2018/01/24        Tong Yuehong
 */

package article

import (
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs/article"
)

type req struct {
	Art    article.CreateArticle `json:"art"`
	UserId uint64                `json:"userId"`
}

// AddArticle - add article.
func AddArticle(this *server.Context) error {
	var reqAdd req

	if err := this.JSONBody(&reqAdd); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	id, err := article.ArticleService.InsertArticle(reqAdd.Art, reqAdd.UserId)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, id)
}

// GetByModuleID gets articles by ModuleId.
func GetByModuleID(this *server.Context) error {
	var req struct {
		artId    string
		moduleId string
	}
	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByModuleID(req.artId, req.moduleId)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByThemeID - gets articles by ThemeId.
func GetByThemeID(this *server.Context) error {
	var req struct {
		artId   string
		ThemeId string
	}
	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByModuleID(req.artId, req.ThemeId)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// DeleteArt deletes article
func DeleteArt(this *server.Context) error {
	var title struct {
		Title string
	}

	if err := this.JSONBody(&title); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.DeleteArt(title.Title)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, 0)
}