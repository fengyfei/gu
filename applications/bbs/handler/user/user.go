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
 *     Modify : 2018/03/04
 */

package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/applications/bbs/util"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/user"
)

const (
	WechatURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	APPID     = ""
	SECRET    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
)

// WechatLogin
func WechatLogin(this *server.Context) error {
	var (
		wechatCode user.WechatCode
		wechatData user.WechatData
	)

	if err := this.JSONBody(&wechatCode); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&wechatCode); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	url := fmt.Sprintf(WechatURL, APPID, SECRET, wechatCode.Code)
	wechatRes, err := http.Get(url)

	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	if wechatRes.StatusCode != http.StatusOK {
		logger.Error("Can't get session key from weixin server: response status code", wechatRes.StatusCode)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	err = json.NewDecoder(wechatRes.Body).Decode(&wechatData)
	if err != nil {
		logger.Error("Error in parsing response:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	wechatLogin := user.WechatLogin{
		UnionID: wechatData.UnionID,
	}

	// connect to mysql
	conn, err := initialize.Pool.Get()
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	u, avatar, err := user.UserServer.WeChatLogin(conn, &wechatLogin)
	if err != nil {
		logger.Error("Wechat login failed.")
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(u.UserID, wechatData.SessionKey, false)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	userData := user.UserData{
		Token:    token,
		UserName: u.UserName,
		Phone:    u.Phone,
		Avatar:   avatar.Avatar,
		Sex:      u.Sex,
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, userData)
}

// PhoneRegister register by phoneNumber
func PhoneRegister(c *server.Context) error {
	var register user.PhoneRegister

	if err := c.JSONBody(&register); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&register); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := initialize.Pool.Get()
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.PhoneRegister(conn, &register)
	if err != nil {
		logger.Error("Register failed.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneLogin login by phone
func PhoneLogin(c *server.Context) error {
	var phoneLogin user.PhoneLogin

	if err := c.JSONBody(&phoneLogin); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&phoneLogin); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := initialize.Pool.Get()
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, avatar, err := user.UserServer.PhoneLogin(conn, &phoneLogin)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err := util.NewToken(u.UserID, "", false)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData := user.UserData{
		Token:    token,
		UserName: u.UserName,
		Phone:    u.Phone,
		Sex:      u.Sex,
		Avatar:   avatar.Avatar,
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// Change User Info
func ChangeUserInfo(this *server.Context) error {
	var changeInfo user.ChangeInfo

	if err := this.JSONBody(&changeInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&changeInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	token, err := util.Parse(this)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwtgo.MapClaims)
	userid := uint32(claims[util.UserId].(float64))

	conn, err := initialize.Pool.Get()
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	err = user.UserServer.ChangeInfo(conn, userid, &changeInfo)
	if err != nil {
		logger.Error(err)
		core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Change password
func ChangePassword(c *server.Context) error {
	var change user.ChangePass

	if err := c.JSONBody(&change); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&change); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	token, err := util.Parse(c)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwtgo.MapClaims)
	userid := uint32(claims[util.UserId].(float64))

	conn, err := initialize.Pool.Get()
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.ChangePassword(conn, userid, &change)
	if err != nil {
		logger.Error("Error in changing password:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
