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
	"github.com/fengyfei/gu/models/blog/bio"
	"github.com/fengyfei/gu/pkg/log"
)

type Bio struct {
	baseController
}

// @router /blog/bio/info [get]
func (tc *Bio) BioInfo() {
	info, err := bio.Service.GetBio()

	if err != nil {
		log.GlobalLogReporter.Error(err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
	} else {
		log.GlobalLogReporter.Debug("Get bio information success.")

		tc.Data["json"] = map[string]interface{}{
			common.RespKeyStatus: common.ErrSucceed,
			common.RespKeyData:   info,
		}
	}

	tc.ServeJSON()
}

// @router /blog/bio/create [post]
func (tc *Bio) Create() {
	var info bio.MDCreateBio

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &info)

	if err != nil {
		log.GlobalLogReporter.Error(err)

		tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrInvalidParam}
	} else {
		err := bio.Service.Create(&info)

		if err != nil {
			log.GlobalLogReporter.Error(err)

			tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrMongoDB}
		} else {
			log.GlobalLogReporter.Debug("Create bio information success.")

			tc.Data["json"] = map[string]interface{}{common.RespKeyStatus: common.ErrSucceed}
		}
	}

	tc.ServeJSON()
}
