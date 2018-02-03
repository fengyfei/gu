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
 *     Initial: 2018/01/21        Chen Yanchen
 */

package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/bbs/util"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/user"
)

var (
	// APPID
	APPID = ""

	// SECRET
	SECRET   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	typeUser = false
)

type (
	// WechatLoginReq
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
)

// WechatLogin
func WechatLogin(this *server.Context) error {
	var (
		wechatUser WechatLoginReq
		userId     uint64
		url        string
		wechatData wechatLogin
		wechatRes  *http.Response
		con        []byte
		token      string
	)

	if err := this.JSONBody(&wechatUser); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := this.Validate(&wechatUser)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", APPID, SECRET, wechatUser.WechatCode)

	wechatRes, err = http.Get(url)
	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	con, _ = ioutil.ReadAll(wechatRes.Body)
	err = json.Unmarshal(con, &wechatData)
	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	if wechatData.data.errmsg != "" {
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	userId, err = user.UserServer.WeChatLogin(&wechatUser.UserName, &wechatData.data.unionid)
	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	token, err = util.NewToken(userId, typeUser)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, token)
}

// PhoneRegister register by phoneNumber
func PhoneRegister(c *server.Context) error {
	var (
		req user.PhoneRegister
		err error
	)

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = user.UserServer.PhoneRegister(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneLogin login by phone
func PhoneLogin(c *server.Context) error {
	var (
		req   user.PhoneLogin
		err   error
		token string
		uid   uint64
	)

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	uid, err = user.UserServer.PhoneLogin(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err = util.NewToken(uid, typeUser)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, token)
}

// ChangeUsername
func ChangeUsername(this *server.Context) error {
	var req struct {
		NewName string `json:"newname" validate:"required,alphanum,min=6,max=30"`
	}

	userID := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64)
	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	err := user.UserServer.ChangeName(uint64(userID), &req.NewName)
	if err != nil {
		logger.Error(err)
		core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// ChangeAvatar
func ChangeAvatar(this *server.Context) error {
	var req struct {
		NewAvatar string `json:"newname" validate:"required,alphanum,min=6,max=30"`
	}
	userID := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	avatarId, err := user.UserServer.ChangeAvatar(uint64(userID), &req.NewAvatar)
	if err != nil {
		logger.Error(err)
		core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, avatarId)
}
