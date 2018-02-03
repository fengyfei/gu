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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	models "github.com/fengyfei/gu/models/shop/user"
)

const (
	WechatURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	APPID     = ""
	SECRET    = ""

	typeUser = false
)

// Login by wechat
func WechatLogin(c *server.Context) error {
	var (
		wechatCode models.WechatCode
		wechatData models.WechatData
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

	body, err := ioutil.ReadAll(wechatResp.Body)
	err = json.Unmarshal(body, &wechatData)
	if err != nil {
		logger.Error("Error in parsing response:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrWechatAuth, nil)
	}

	wechatLogin := models.WechatLogin{
		OpenID:  wechatData.OpenID,
		UnionID: wechatData.UnionID,
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	user, err := models.Service.WechatLogin(conn, &wechatLogin)
	if err != nil {
		logger.Error("Error in Wechat Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(user.ID, wechatData.SessionKey, user.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData := models.UserData{
		Token:    token,
		UserName: user.UserName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Sex:      user.Sex,
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// Add a phone number
func AddPhone(c *server.Context) error {
	var phone models.AddPhone

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

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userid := uint(claims["userid"].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.AddPhone(conn, userid, &phone)
	if err != nil {
		logger.Error("Error in adding a phone number:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Change information
func ChangeInfo(c *server.Context) error {
	var change models.ChangeInfo

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

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userid := uint(claims["userid"].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.ChangeInfo(conn, userid, &change)
	if err != nil {
		logger.Error("Error in changing informantion:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Register by phone
func PhoneRegister(c *server.Context) error {
	var register models.PhoneRegister

	err := c.JSONBody(&register)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&register)
	if err != nil || !util.IsValidPhone(register.Phone) {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.PhoneRegister(conn, &register)
	if err != nil {
		logger.Error("Error in registering by phone:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Login by phone
func PhoneLogin(c *server.Context) error {
	var login models.PhoneLogin

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

	user, err := models.Service.PhoneLogin(conn, &login)
	if err != nil {
		logger.Error("Error in Phone Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err := util.NewToken(user.ID, "", user.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData := models.UserData{
		Token:    token,
		UserName: user.UserName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Sex:      user.Sex,
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// Change password
func ChangePassword(c *server.Context) error {
	var change models.ChangePass

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

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userid := uint(claims["userid"].(float64))

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.ChangePassword(conn, userid, &change)
	if err != nil {
		logger.Error("Error in changing password:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
