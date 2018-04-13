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
 *     Initial: 2018/01/26        Tong Yuehong
 */

package article

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/bbs/article"
)

type (
	categoryVisit struct {
		Num        int64  `json:"num"           validate:"required"`
		CategoryID string `json:"categoryid"`
	}

	createTag struct {
		CategoryID string `json:"categoryid"`
		Tag        string `json:"tag"`
	}

	category struct {
		CategoryID string `json:"categoryid"`
	}

	createCategory struct {
		Name string `json:"name" validate:"required"`
	}
)

// AddCategory add category.
func AddCategory(this *server.Context) error {
	var (
		create createCategory
	)

	if err := this.JSONBody(&create); err != nil {
		logger.Error("AddCategory json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&create); err != nil {
		logger.Error("AddCategory Validate()", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	category := &article.Category{
		Name: create.Name,
	}

	err := article.CategoryService.CreateCategory(category)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}

// UpdateCategoryVisit updates CategoryVisit.
func UpdateCategoryVisit(this *server.Context) error {
	var (
		visit categoryVisit
	)

	if err := this.JSONBody(&visit); err != nil {
		logger.Error("UpdateCategoryVisit json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&visit); err != nil {
		logger.Error("UpdateCategoryVisit Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CategoryService.UpdateCategoryVisit(visit.Num, visit.CategoryID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}

// AddTag add tag.
func AddTag(this *server.Context) error {
	var (
		createTag createTag
	)

	if err := this.JSONBody(&createTag); err != nil {
		logger.Error("AddTag", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&createTag); err != nil {
		logger.Error("AddTag Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CategoryService.CreateTag(createTag.CategoryID, createTag.Tag)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}

// DeleteCategory delete category.
func DeleteCategory(this *server.Context) error {
	var (
		category category
	)

	if err := this.JSONBody(&category); err != nil {
		logger.Error("DeleteCategory json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CategoryService.DeleteCategory(category.CategoryID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}

// DeleteTag delete tag.
func DeleteTag(this *server.Context) error {
	var (
		tag struct {
			CategoryID string `json:"categoryid"`
			TagID      string `json:"tagid"`
		}
	)
	if err := this.JSONBody(&tag); err != nil {
		logger.Error("DeleteTag json: ", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(tag.CategoryID) || !bson.IsObjectIdHex(tag.TagID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CategoryService.DeleteTag(tag.CategoryID, tag.TagID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}

// ListTags return all tags.
func ListTags(this *server.Context) error {
	var (
		category category
	)

	if err := this.JSONBody(&category); err != nil {
		logger.Error("ListTags json: ", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(category.CategoryID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.CategoryService.ListInfo(category.CategoryID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	tags := make([]string, len(list.Tags))
	for i, tag := range list.Tags {
		if tag.Active == true {
			tags[i] = tag.Name
		}
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, tags)
}

// AllCategories returns all categories.
func AllCategories(this *server.Context) error {
	list, err := article.CategoryService.AllCategories()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// ListRecommend return recommended category.
func ListRecommend(this *server.Context) error {
	list, err := article.CategoryService.ListRecommend()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}
