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
 *     Initial: 2017/11/13        Jia Chenhui
 */

package controller

import (
	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/polaris/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/staff"
)

// Permission - Permission associated handler.
type Permission struct {
	base.Controller
}

type (
	// createPmsReq - The request struct that create permission information.
	createPmsReq struct {
		URL    *string `json:"url" validate:"required,printascii,contains=/"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// removePmsReq - The request struct that remove permission information.
	removePmsReq struct {
		URL    *string `json:"url" validate:"required,printascii,contains=/"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// pmsInfoResp - The response struct that represents role information of URL.
	pmsInfoResp struct {
		URL    string `json:"url"`
		RoleId int16  `json:"roleid"`
	}
)

// Create - Create permission information.
func (p *Permission) Create() {
	var (
		err  error
		req  createPmsReq
		conn orm.Connection
	)

	if err = json.Unmarshal(p.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	if err = p.Validate(&req); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	if err = staff.Service.AddURLPermission(conn, req.URL, req.RoleId); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	p.ServeJSON(true)
}

// Remove - Remove permission information.
func (p *Permission) Remove() {
	var (
		err  error
		req  removePmsReq
		conn orm.Connection
	)

	if err = json.Unmarshal(p.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	if err = p.Validate(&req); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	if err = staff.Service.RemoveURLPermission(conn, req.URL, req.RoleId); err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	p.ServeJSON(true)
}

// List - Get a list of permission for specified URL.
func (p *Permission) List() {
	var (
		err        error
		permission pmsInfoResp
		conn       orm.Connection
		plist      []staff.Permission
		resp       []pmsInfoResp = make([]pmsInfoResp, 0)
	)

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	plist, err = staff.Service.Permissions(conn)
	if err != nil {
		logger.Error(err)
		p.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	for _, p := range plist {
		permission = pmsInfoResp{
			URL:    p.URL,
			RoleId: p.RoleId,
		}

		resp = append(resp, permission)
	}

	p.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	}

finish:
	p.ServeJSON(true)
}
