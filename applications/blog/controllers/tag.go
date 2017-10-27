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
	"encoding/json"

	"github.com/fengyfei/gu/common"
	"github.com/fengyfei/gu/models/blog/tag"
	"github.com/fengyfei/gu/models/req"
	"github.com/fengyfei/gu/pkg/log"
)

type Tag struct {
	baseController
}

// @router /blog/tag/list [get]
func (tc *Tag) List() {
	tagList, err := tag.Service.GetList()

	if err != nil {
		log.Logger.Error("Tag.List returned error: %v", err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
	} else {
		log.Logger.Debug("Get tag list success.")

		tc.Data["json"] = map[string]interface{}{
			common.RespKeyStatus: common.ErrSucceed,
			common.RespKeyData:   tagList,
		}
	}

	tc.ServeJSON()
}

// @router /blog/tag/activelist [get]
func (tc *Tag) ActiveList() {
	tagList, err := tag.Service.GetActiveList()

	if err != nil {
		log.Logger.Error("Tag.ActiveList returned error: %v", err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
	} else {
		log.Logger.Debug("Get tag active list success.")

		tc.Data["json"] = map[string]interface{}{
			common.RespKeyStatus: common.ErrSucceed,
			common.RespKeyData:   tagList,
		}
	}

	tc.ServeJSON()
}

// @router /blog/tag/info [post]
func (tc *Tag) TagInfo() {
	var info req.MDTagInfoReq

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &info)

	if err != nil {
		log.Logger.Error("Tag.TagInfo returned error: %v", err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}
	} else {
		log.Logger.Debug("Get tag information success.")

		tc.Data["json"] = map[string]interface{}{
			common.RespKeyStatus: common.ErrSucceed,
			common.RespKeyData:   info,
		}
	}

	tc.ServeJSON()
}

// @router /blog/tag/create [post]
func (tc *Tag) Create() {
	var info tag.MDCreateTag

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &info)

	if err != nil {
		log.Logger.Error("Tag.Create returned error: %v", err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}
	} else {
		id, err := tag.Service.Create(info.Tag)

		if err != nil {
			log.Logger.Error("Tag.Create returned error: %v", err)

			tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
		} else {
			log.Logger.Debug("Create tag information success.")

			tc.Data["json"] = map[string]interface{}{
				common.RespKeyStatus: common.ErrSucceed,
				common.RespKeyData:   id,
			}
		}
	}

	tc.ServeJSON()
}

// @router /blog/tag/modify [post]
func (tc *Tag) Modify() {
	var info tag.MDModifyTag

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &info)

	if err != nil {
		log.Logger.Error("Tag.Modify returned error: %v", err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}
	} else {
		err := tag.Service.Modify(&info)

		if err != nil {
			log.Logger.Error("Tag.Modify returned error: %v", err)

			tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
		} else {
			log.Logger.Debug("Modify tag success.")

			tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed}
		}
	}

	tc.ServeJSON()
}
