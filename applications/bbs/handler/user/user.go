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
 *     Modify : 2018/03/25        Tong Yuehong
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
	"github.com/fengyfei/gu/models/bbs/article"
	"github.com/fengyfei/gu/models/user"
	"time"
)

const (
	WechatURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	APPID     = ""
	SECRET    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
)

type (
	bbsUserInfo struct {
		UserID     uint32
		Username   string
		Phone      string
		Avatar     string
		Sex        uint8
		ArticleNum int
		CommentNum int
		LastLogin  time.Time
	}

	wechatCode struct {
		UserName string `json:"username"`
		Avatar   string `json:"avatar"`
		Sex      uint8  `json:"sex"`
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Code     string `json:"code" validate:"required"`
	}

	wechatData struct {
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
	}

	wechatLogin struct {
		UnionID    string `json:"unionid"`
		SessionKey string `json:"session_key"`
	}

	userData struct {
		Token    string `json:"token"`
		UserName string `json:"username"` // todo
		Phone    string `json:"phone"`    // todo
		Avatar   string `json:"avatar"`
		Sex      uint8  `json:"sex"`
	}

	wechatPhone struct {
		Phone string `json:"phone" validate:"required,numeric,len=11"`
	}

	phoneRegister struct {
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
		UserName string `json:"username" validate:"required,alphaunicode,min=2,max=30"`
	}

	phoneLogin struct {
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
	}

	changeInfo struct {
		UserName string `json:"username"`
		Sex      uint8  `json:"sex"`
		Avatar   string `json:"avatar"`
	}

	changePass struct {
		OldPass string `json:"oldpass" validate:"required,min=6,max=30"`
		NewPass string `json:"newpass" validate:"required,min=6,max=30"`
	}
)

// WechatLogin
func WechatLogin(this *server.Context) error {
	var (
		wechatCode wechatCode
		wechatData wechatData
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

	wechatLogin := wechatLogin{
		UnionID: wechatData.UnionID,
	}

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	u, err := user.UserService.WeChatLogin(conn, wechatLogin.UnionID, wechatCode.UserName, wechatCode.Avatar)
	if err != nil {
		logger.Error("Wechat login failed.")
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(u.UserID, wechatData.SessionKey, u.IsAdmin)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	userData := userData{
		Token:    token,
		UserName: u.UserName,
		Phone:    u.Phone,
		Avatar:   u.Avatar,
		Sex:      u.Sex,
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, userData)
}

// Add phoneNum
func AddPhone(c *server.Context) error {
	var phone wechatPhone

	if err := c.JSONBody(&phone); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&phone); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.ValidatePhone(phone.Phone) {
		logger.Error("Invalid phone.")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userid := c.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(uint32)
	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.AddPhone(conn, userid, phone.Phone)
	if err != nil {
		logger.Error("Add phoneNumber failed.")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneRegister register by phoneNumber
func PhoneRegister(c *server.Context) error {
	var register phoneRegister

	if err := c.JSONBody(&register); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&register); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.ValidatePhone(register.Phone) {
		logger.Error("Phone not right.")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}
	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.PhoneRegister(conn, register.UserName, register.Phone, register.Password)
	if err != nil {
		logger.Error("Register failed.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneLogin login by phone
func PhoneLogin(c *server.Context) error {
	var phoneLogin phoneLogin

	if err := c.JSONBody(&phoneLogin); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&phoneLogin); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, err := user.UserService.PhoneLogin(conn, phoneLogin.Phone, phoneLogin.Password)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err := util.NewToken(u.UserID, "", false)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData := userData{
		Token:    token,
		UserName: u.UserName,
		Phone:    u.Phone,
		Sex:      u.Sex,
		Avatar:   u.Avatar,
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, userData)
}

// Change User Info
func ChangeUserInfo(this *server.Context) error {
	var changeInfo changeInfo

	if err := this.JSONBody(&changeInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&changeInfo); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	userid := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(uint32)

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	err = user.UserService.ChangeInfo(conn, userid, changeInfo.UserName, changeInfo.Sex)
	if err != nil {
		logger.Error(err)
	}

	if len(changeInfo.Avatar) > 0 {
		changeInfo.Avatar, err = user.SavePicture(changeInfo.Avatar, "avatar/")
			if err != nil {
				logger.Error(err)
				return core.WriteStatusAndDataJSON(this, constants.ErrInternalServerError, nil)
			}
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Change password
func ChangePassword(c *server.Context) error {
	var change changePass

	if err := c.JSONBody(&change); err != nil {
		logger.Error("JsonBody Error.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&change); err != nil {
		logger.Error("Validate Error.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userid := c.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(uint32)

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.ChangePassword(conn, userid, change.OldPass, change.NewPass)
	if err != nil {
		logger.Error("Error in changing password:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// BbsUserInfo get user info in bbs.
func BbsUserInfo(c *server.Context) error {
	userid := c.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(uint32)

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, err := user.UserService.GetUserByID(conn, userid)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	artNum, err := article.ArticleService.ArtNum(userid)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	comNum, err := article.CommentService.CommentNum(userid)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	bbsUserInfo := &bbsUserInfo{
		UserID:     userid,
		Username:   u.UserName,
		Phone:      u.Phone,
		Avatar:     u.Avatar,
		Sex:        u.Sex,
		ArticleNum: artNum,
		CommentNum: comNum,
		LastLogin:  u.LastLogin,
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, bbsUserInfo)
}
