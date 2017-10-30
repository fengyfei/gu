/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 */

package controllers

import (
	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/tag"
)

// Tag - tag associated handlers
type Tag struct {
	base.Controller
}

// infoReq - the request struct that get tag information by id.
type infoReq struct {
	TagID string `json:"tagid" validate:"required"`
}

// createReq - the request struct that create tag information.
type createReq struct {
	Tag string `json:"tag" validate:"required"`
}

// modifyReq - the request struct that modify the tag information.
type modifyReq struct {
	TagID  string `json:"tagid" validate:"required"`
	Tag    string `json:"tag" validate:"required"`
	Active *bool  `json:"active" validate:"required"`
}

// List all tags;
func (tc *Tag) List() {
	tagList, err := tag.Service.GetList()

	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}
	} else {
		logger.Debug("Get tag list success.")

		tc.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: constants.ErrSucceed,
			constants.RespKeyData:   tagList,
		}
	}

	tc.ServeJSON()
}

// ActiveList returns all active tags.
func (tc *Tag) ActiveList() {
	tagList, err := tag.Service.GetActiveList()

	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}
	} else {
		logger.Debug("Get tag active list success.")

		tc.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: constants.ErrSucceed,
			constants.RespKeyData:   tagList,
		}
	}

	tc.ServeJSON()
}

// Info for specific tag
func (tc *Tag) Info() {
	var (
		req  infoReq
		resp tag.Tag
	)

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: constants.ErrInvalidParam,
			constants.RespKeyData:   req,
		}
		goto finish
	}

	err = tc.Validate(&req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	resp, err = tag.Service.GetByID(&req.TagID)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}
		goto finish
	}

	logger.Debug("Get tag information success.")

	tc.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	}

finish:
	tc.ServeJSON()
}

// Create a new tag.
func (tc *Tag) Create() {
	var (
		req  createReq
		resp string
	)

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: constants.ErrInvalidParam,
			constants.RespKeyData:   req,
		}
		goto finish
	}

	err = tc.Validate(&req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	resp, err = tag.Service.Create(&req.Tag)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}
		goto finish
	}

	logger.Debug("Create tag information success.")

	tc.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	}

finish:
	tc.ServeJSON()
}

// Modify a specific tag.
func (tc *Tag) Modify() {
	var req modifyReq

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: constants.ErrInvalidParam,
			constants.RespKeyData:   req,
		}
		goto finish
	}

	err = tc.Validate(&req)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = tag.Service.Modify(&req.TagID, &req.Tag, req.Active)
	if err != nil {
		logger.Error(err)

		tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMongoDB}
		goto finish
	}

	logger.Debug("Modify tag success.")

	tc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	tc.ServeJSON()
}
