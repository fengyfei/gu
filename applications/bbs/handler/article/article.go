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
	jwtgo "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/bbs/article"
)

type (
	title struct {
		Title string `json:"title"`
	}
)

// AddArticle - add article.
func AddArticle(this *server.Context) error {
	var (
		reqAdd article.CreateArticle
	)

	if err := this.JSONBody(&reqAdd); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	userID := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(uint32)
	if err := this.Validate(&reqAdd); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	id, err := article.ArticleService.Insert(reqAdd, userID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, id)
}

// GetByModuleID gets articles by ModuleID.
func GetByModuleID(this *server.Context) error {
	var (
		req struct {
			Page   int    `json:"page"`
			Module string `json:"module"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByModuleID(req.Page, req.Module)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByThemeID - gets articles by ThemeID.
func GetByThemeID(this *server.Context) error {
	var (
		req struct {
			Page   int    `json:"page"`
			Module string `json:"module"`
			Theme  string `json:"theme"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByThemeID(req.Page, req.Module, req.Theme)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByTitle - gets articles by ThemeID.
func GetByTitle(this *server.Context) error {
	var (
		title title
	)

	if err := this.JSONBody(&title); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByTitle(title.Title)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByUserID - gets articles by ThemeID.
func GetByUserID(this *server.Context) error {
	var (
		user struct {
			UserID string `json:"userID"`
		}
	)

	if err := this.JSONBody(&user); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByTitle(user.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// DeleteArt deletes article
func DeleteArt(this *server.Context) error {
	var (
		title title
	)

	if err := this.JSONBody(&title); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.Delete(title.Title)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// UpdateTimes update times.
func UpdateTimes(this *server.Context) error {
	var (
		times struct {
			Num   int64
			ArtID string
		}
	)

	if err := this.JSONBody(&times); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(times.ArtID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.UpdateTimes(times.Num, times.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Recommend return popular articles.
func Recommend(this *server.Context) error {
	var (
		req struct {
			Page int `json:"page"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.Recommend(req.Page)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}
