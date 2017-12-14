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
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/applications/beego/shop/util"
	"github.com/fengyfei/gu/libs/orm"
)

var (
	APPID    = ""
	SECRET   = ""
	typeUser = false
)

type (
	UserController struct {
		base.Controller
	}

	WechatLoginReq struct {
		UserName   string `json:"userName" validate:"required,alphanum,min=6,max=30"`
		WechatCode string `json:"wechatCode" validate:"required"`
	}

	wechatLogin struct {
		data wechatLoginData
	}

	wechatLoginData struct {
		errmsg  string
		unionid string
	}

	phoneRegisterReq struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
		NickName string `json:"name" validate:"required,alphaunicode,min=2,max=30"`
	}

	phoneLoginReq struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
	}

	changePassReq struct {
		OldPass string `json:"oldPass" validate:"required,min=6,max=30"`
		NewPass string `json:"newPass" validate:"required,min=6,max=30"`
	}
)

func (u *UserController) WechatLogin() {
	var
	(
		wechatUser WechatLoginReq
		err        error
		userId     int32
		url        string
		wechatData wechatLogin
		wechatRes  *http.Response
		con        []byte
		token      string
	)
	json.Unmarshal(u.Ctx.Input.RequestBody, &wechatUser)

	conn, err := mysql.Pool.Get()
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = u.Validate(&wechatUser)
	if err != nil {
		logger.Error(err)
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", APPID, SECRET, wechatUser.WechatCode)

	wechatRes, err = http.Get(url)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	con, _ = ioutil.ReadAll(wechatRes.Body)
	err = json.Unmarshal(con, &wechatData)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	if wechatData.data.errmsg != "" {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	userId, err = user.Service.WechatLogin(conn, &wechatUser.UserName, &wechatData.data.unionid)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	token, err = util.NewToken(userId, typeUser)
	if err != nil {
		logger.Error(err)
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}
	u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyToken: token}

finish:
	u.ServeJSON(true)
}

// Register by phoneNumber
func (this *UserController) PhoneRegister() {
	var (
		registerReq phoneRegisterReq
		err         error
	)

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &registerReq)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = this.Validate(&registerReq)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = user.Service.PhoneRegister(conn, &registerReq.Phone, &registerReq.Password, &registerReq.NickName)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)

}

func (this *UserController) PhoneLogin() {
	var (
		loginReq phoneLoginReq
		err      error
		token    string
		uid      int32
	)

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &loginReq)
	fmt.Println(loginReq)
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
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrAccount}
		goto finish
	}

	token, err = util.NewToken(uid, typeUser)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyToken: token}

finish:
	this.ServeJSON(true)
}

func (this *UserController) ChangePassword() {
	var (
		req    changePassReq
		userId int32
		conn   orm.Connection
		err    error
	)

	userId = this.Ctx.Request.Context().Value("userId").(int32)

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

	err = user.Service.ChangePassword(conn, userId, &req.OldPass, &req.NewPass)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}
