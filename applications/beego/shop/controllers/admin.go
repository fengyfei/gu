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
 *     Initial: 2017/11/18        ShiChao
 */

package controllers

import (
	"encoding/json"
	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/models/shop/user"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/applications/beego/shop/util"
	"github.com/fengyfei/gu/libs/orm"
)

var (
	typeAdmin = true
)

type (
	AdminController struct {
		base.Controller
	}

	loginReq struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
	}

	adminChangePassReq struct {
		OldPass string `json:"oldPass" validate:"required,min=6,max=30"`
		NewPass string `json:"newPass" validate:"required,min=6,max=30"`
	}
)

func (this *UserController) Login() {
	var (
		loginReq phoneLoginReq
		err      error
		token    string
		uid      uint
	)

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &loginReq)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = this.Validate(&loginReq)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	uid, err = user.Service.PhoneLogin(conn, &loginReq.Phone, &loginReq.Password)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	token, err = util.NewToken(uid, typeAdmin)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyToken: token}

finish:
	this.ServeJSON(true)
}

func (this *UserController) ChangeAdminPassword() {
	var (
		req  adminChangePassReq
		conn orm.Connection
		err  error
	)
	adminId := this.Ctx.Request.Context().Value("userId").(uint)
	isAdmin := this.Ctx.Request.Context().Value("isAdmin").(bool)
	if !isAdmin {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}
		goto finish
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = this.Validate(&req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = user.Service.ChangePassword(conn, adminId, &req.OldPass, &req.NewPass)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}
