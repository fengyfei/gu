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
 *    Initial : 2018/03/25        Tong Yuehong
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
	WechatURL   = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	UserInfoURL = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	APPID       = ""
	SECRET      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
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
		Code string `json:"code" validate:"required"`
	}

	wechatData struct {
		AccessToken string
		OpenID      string
	}

	userData struct {
		OpenID   string `json:"openid"`
		UserName string `json:"nickname"`
		Avatar   string `json:"headimgurl"`
		Sex      uint8  `json:"sex"`
	}

	wechatPhone struct {
		Phone string `json:"phone"`
	}

	changeInfo struct {
		UserName string `json:"username"`
		Sex    uint8 `json:"sex"`
	}

	changeAvatar struct {
		Avatar string `json:"avatar"`
	}

	changePass struct {
		OldPass string `json:"oldpass" validate:"required,min=6,max=30"`
		NewPass string `json:"newpass" validate:"required,min=6,max=30"`
	}

	userid struct {
		UserID uint32 `json:"userid"`
	}
)

// WechatLogin
func WechatLogin(this *server.Context) error {
	var (
		wechatCode wechatCode
		wechatData wechatData
		userData   userData
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

	url = fmt.Sprintf(UserInfoURL, wechatData.AccessToken, wechatData.OpenID)
	tokenRes, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	if tokenRes.StatusCode != http.StatusOK {
		logger.Error("Can't get session key from weixin server: response status code", wechatRes.StatusCode)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	err = json.NewDecoder(tokenRes.Body).Decode(&userData)
	if err != nil {
		logger.Error("Error in parsing response:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrWechatAuth, nil)
	}

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	userData.OpenID = wechatData.OpenID
	u, err := user.UserService.WeChatLogin(conn, userData.OpenID, userData.UserName, userData.Avatar)
	if err != nil {
		logger.Error("Wechat login failed.")
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(u.UserID, u.IsAdmin)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, token)
}

// AddPhone add phone.
func AddPhone(c *server.Context) error {
	var (
		phone wechatPhone
	)

	if err := c.JSONBody(&phone); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.ValidatePhone(phone.Phone) {
		logger.Error("Invalid phone.")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userid := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64))
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

// ChangeUserInfo  changes user's information.
func ChangeUserInfo(c *server.Context) error {
	var (
		info changeInfo
	)

	if err := c.JSONBody(&info); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&info); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userid := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64))

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.ChangeInfo(conn, userid, info.UserName, info.Sex)
	if err != nil {
		logger.Error(err)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// ChangePassword changes password.
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

// ChangeAvatar changes avatar.
func ChangeAvatar(c *server.Context) error {
	var avatar changeAvatar

	if err := c.JSONBody(&avatar); err != nil {
		logger.Error("JsonBody Error.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&avatar); err != nil {
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

	path, err := user.SavePicture(avatar.Avatar, "./bbs/", userid)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
	}

	err = user.UserService.ChangeAvatar(conn, userid, path)

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// BbsUserInfo get user info in bbs.
func BbsUserInfo(c *server.Context) error {
	var (
		userid userid
	)

	if err := c.JSONBody(&userid); err != nil {
		logger.Error("JsonBody Error.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := initialize.Pool.Get()
	defer initialize.Pool.Release(conn)
	if err != nil {
		logger.Error("Can not connected mysql.", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, err := user.UserService.GetUserByID(conn, userid.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	artNum, err := article.ArticleService.ArtNum(userid.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	comNum, err := article.CommentService.CommentNum(userid.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	bbsUserInfo := &bbsUserInfo{
		UserID:     userid.UserID,
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
