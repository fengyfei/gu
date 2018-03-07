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
 *     Initial: 2018/02/01        Shi Ruitao
 *     Modify:  2018/02/01        Li Zebang
 */

package handler

import (
	"fmt"
	"net/http"

	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/core"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/user"
)

const (
	WechatURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	APPID     = ""
	SECRET    = ""
)

// Login by wechat
func WechatLogin(c *server.Context) error {
	var (
		wechatCode user.WechatCode
		wechatData user.WechatData
	)

	err := c.JSONBody(&wechatCode)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&wechatCode)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	url := fmt.Sprintf(WechatURL, APPID, SECRET, wechatCode.Code)
	wechatResp, err := http.Get(url)
	if err != nil {
		logger.Error("Error in getting session key from weixin server:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrWechatAuth, nil)
	}

	if wechatResp.StatusCode != http.StatusOK {
		logger.Error("Can't get session key from weixin server: response status code", wechatResp.StatusCode)
		return core.WriteStatusAndDataJSON(c, constants.ErrWechatAuth, nil)
	}

	err = json.NewDecoder(wechatResp.Body).Decode(&wechatData)
	if err != nil {
		logger.Error("Error in parsing response:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrWechatAuth, nil)
	}

	wechatLogin := user.WechatLogin{
		UnionID: wechatData.UnionID,
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, avatar, err := user.UserServer.WeChatLogin(conn, &wechatLogin)
	if err != nil {
		logger.Error("Error in Wechat Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(u.UserID, wechatData.SessionKey, u.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData := user.UserData{
		Token:    token,
		UserName: u.UserName,
		Sex:      u.Sex,
		Avatar:   avatar.Avatar,
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// Wechat add a phone number
func AddPhone(c *server.Context) error {
	var phone user.WechatPhone

	err := c.JSONBody(&phone)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&phone)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.IsValidPhone(phone.Phone) {
		logger.Error("Invalid phone number")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.AddPhone(conn, userID, &phone)
	if err != nil {
		logger.Error("Error in adding a phone number:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Change information
func ChangeInfo(c *server.Context) error {
	var change user.ChangeInfo

	err := c.JSONBody(&change)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&change)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.ChangeInfo(conn, userID, &change)
	if err != nil {
		logger.Error("Error in changing informantion:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Register by phone
func PhoneRegister(c *server.Context) error {
	var register user.PhoneRegister

	err := c.JSONBody(&register)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&register)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.IsValidPhone(register.Phone) {
		logger.Error("Invalid phone number")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.PhoneRegister(conn, &register)
	if err != nil {
		logger.Error("Error in registering by phone:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Login by phone
func PhoneLogin(c *server.Context) error {
	var login user.PhoneLogin

	err := c.JSONBody(&login)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&login)
	if err != nil || !util.IsValidPhone(login.Phone) {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, avatar, err := user.UserServer.PhoneLogin(conn, &login)
	if err != nil {
		logger.Error("Error in Phone Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err := util.NewToken(u.UserID, "", u.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
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

// Change password
func ChangePassword(c *server.Context) error {
	var change user.ChangePass

	err := c.JSONBody(&change)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&change)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserServer.ChangePassword(conn, userID, &change)
	if err != nil {
		logger.Error("Error in changing password:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
