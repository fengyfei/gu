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
 *     Initial: 2017/10/27        Jia Chenhui
 *     Modify : 2018/02/02        Tong Yuehong
 *     Modify : 2018/03/25        Chen Yanchen
 */

package article

import (
	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/tag"
)

type (
	// Tag - tag associated handlers
	Tag struct {
		base.Controller
	}

	// tagInfoReq - the request struct that get tag information by id.
	tagInfoReq struct {
		TagID string `json:"tagid" validate:"required,alphanum,len=24"`
	}

	// createTagReq - the request struct that create tag information.
	createTagReq struct {
		Tag string `json:"tag" validate:"required,min=2,max=6"`
	}

	// modifyTagReq - the request struct that modify the tag information.
	modifyTagReq struct {
		TagID  string `json:"tagid" validate:"required,alphanum,len=24"`
		Tag    string `json:"tag"`
		Active *bool  `json:"active"`
	}
)

// ListTags list all tags.
func ListTags(this *server.Context) error {
	tagList, err := tag.TagService.GetList()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, tagList)
}

// TagActiveList returns all active tags.
func TagActiveList(this *server.Context) error {
	tagList, err := tag.TagService.GetActiveTags()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, tagList)
}

// TagInfo return the tag's information by id.
func TagInfo(this *server.Context) error {
	var (
		req tagInfoReq
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	resp, err := tag.TagService.GetByID(&req.TagID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, resp)
}

func CreateTag(this *server.Context) error {
	var (
		req createTagReq
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	resp, err := tag.TagService.Create(&req.Tag)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, resp)
}

// Modify update a specific tag.
func ModifyTag(this *server.Context) error {
	var (
		req modifyTagReq
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := tag.TagService.Modify(&req.TagID, &req.Tag, req.Active)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, nil)
}
