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

	jwtgo "github.com/dgrijalva/jwt-go"
	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/core"
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

// WechatLogin login by wechat
func WechatLogin(c *server.Context) error {
	var (
		wechatCode struct {
			UserName string `json:"user_name"`
			Avatar   string `json:"avatar"`
			Sex      uint8  `json:"sex"`
			Phone    string `json:"phone" validate:"required,numeric,len=11"`
			Code     string `json:"code" validate:"required"`
		}

		wechatData struct {
			SessionKey string `json:"session_key"`
			UnionID    string `json:"unionid"`
		}

		userData struct {
			Token    string `json:"token"`
			UserName string `json:"user_name"`
			Phone    string `json:"phone"`
			Avatar   string `json:"avatar"`
			Sex      uint8  `json:"sex"`
		}
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

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, err := user.UserService.WeChatLogin(conn, wechatData.UnionID, wechatCode.UserName, wechatCode.Avatar)
	if err != nil {
		logger.Error("Error in Wechat Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	token, err := util.NewToken(u.UserID, wechatData.SessionKey, u.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData.Token = token
	userData.UserName = u.UserName
	userData.Phone = u.Phone
	userData.Sex = u.Sex
	userData.Avatar = u.Avatar

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// AddPhone wechat add a phone number
func AddPhone(c *server.Context) error {
	var (
		wechatPhone struct {
			Phone string `json:"phone" validate:"required,numeric,len=11"`
		}
	)

	err := c.JSONBody(&wechatPhone)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&wechatPhone)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !util.IsValidPhone(wechatPhone.Phone) {
		logger.Error("Invalid phone number")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.AddPhone(conn, userID, wechatPhone.Phone)
	if err != nil {
		logger.Error("Error in adding a phone number:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// ChangeInfo change information
func ChangeInfo(c *server.Context) error {
	var (
		change struct {
			UserName string `json:"user_name" validate:"required,alphaunicode,min=2,max=30"`
			Sex      uint8  `json:"sex"`
			Avatar   string `json:"avatar"`
		}
	)

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
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	if len(change.Avatar) > 0 {
		change.Avatar, err = user.SavePicture(change.Avatar, "avatar/", userID)
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
		err = user.UserService.ChangeAvatar(conn, userID, change.Avatar)
	}

	err = user.UserService.ChangeInfo(conn, userID, change.UserName, change.Sex)

	if err != nil {
		logger.Error("Error in changing informantion:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneRegister register by phone
func PhoneRegister(c *server.Context) error {
	var (
		register struct {
			Phone    string `json:"phone" validate:"required,numeric,len=11"`
			Password string `json:"password" validate:"required,min=6,max=30"`
			UserName string `json:"user_name" validate:"required,alphaunicode,min=2,max=30"`
		}
	)

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

	err = user.UserService.PhoneRegister(conn, register.UserName, register.Phone, register.Password)
	if err != nil {
		logger.Error("Error in registering by phone:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// PhoneLogin login by phone
func PhoneLogin(c *server.Context) error {
	var (
		login struct {
			Phone    string `json:"phone" validate:"required,numeric,len=11"`
			Password string `json:"password" validate:"required,min=6,max=30"`
		}
		userData struct {
			Token    string `json:"token"`
			UserName string `json:"user_name"`
			Phone    string `json:"phone"`
			Avatar   string `json:"avatar"`
			Sex      uint8  `json:"sex"`
		}
	)

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
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	u, err := user.UserService.PhoneLogin(conn, login.Phone, login.Password)
	if err != nil {
		logger.Error("Error in Phone Login:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrAccount, nil)
	}

	token, err := util.NewToken(u.UserID, "", u.IsAdmin)
	if err != nil {
		logger.Error("Error in getting token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userData.Token = token
	userData.UserName = u.UserName
	userData.Phone = u.Phone
	userData.Sex = u.Sex
	userData.Avatar = u.Avatar

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, &userData)
}

// ChangePassword change password
func ChangePassword(c *server.Context) error {
	var (
		change struct {
			OldPass string `json:"old_pass" validate:"required,min=6,max=30"`
			NewPass string `json:"new_pass" validate:"required,min=6,max=30"`
		}
	)

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
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = user.UserService.ChangePassword(conn, userID, change.OldPass, change.NewPass)
	if err != nil {
		logger.Error("Error in changing password:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
