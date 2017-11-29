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
 *     Initial: 2017/11/17        Jia Chenhui
 */

package repos

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/models/github/repos"
)

type (
	// createReq - The request struct that create repos information.
	createReq struct {
		Avatar *string   `json:"avatar" validate:"required,url"`
		Name   *string   `json:"name" validate:"required,printascii,excludesall=;0x2C"`
		Image  *string   `json:"image"`
		Intro  *string   `json:"intro"`
		Lang   *[]string `json:"lang"`
	}

	// modifyActiveReq - The request struct that modify repos status.
	modifyActiveReq struct {
		ID     string `json:"id" validate:"required,alphanum,len=24"`
		Active bool   `json:"active"`
	}

	// infoReq - The request struct that get one repos detail information.
	infoReq struct {
		ID string `json:"id" validate:"required,alphanum,len=24"`
	}

	// infoResp - The more detail of repos.
	infoResp struct {
		Avatar  string    `json:"avatar"`
		Name    string    `json:"name"`
		Image   string    `json:"image"`
		Intro   string    `json:"intro"`
		Lang    []string  `json:"lang"`
		Created time.Time `json:"created"`
		Active  bool      `json:"active"`
	}
)

// Create - Create repos information.
func Create(c echo.Context) error {
	var (
		err error
		req createReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	id, err := repos.Service.Create(req.Avatar, req.Name, req.Image, req.Intro, req.Lang)
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyID:     id,
	})
}

// ModifyActive - Modify repos status.
func ModifyActive(c echo.Context) error {
	var (
		err error
		req modifyActiveReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = repos.Service.ModifyActive(&req.ID, req.Active); err != nil {
		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// List - Get all the repos.
func List(c echo.Context) error {
	var resp []infoResp

	rlist, err := repos.Service.List()
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	for _, r := range rlist {
		info := infoResp{
			Avatar:  r.Avatar,
			Name:    r.Name,
			Image:   r.Image,
			Intro:   r.Intro,
			Lang:    r.Lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// ActiveList - Get all the active repos.
func ActiveList(c echo.Context) error {
	var resp []infoResp

	rlist, err := repos.Service.ActiveList()
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	for _, r := range rlist {
		info := infoResp{
			Avatar:  r.Avatar,
			Name:    r.Name,
			Image:   r.Image,
			Intro:   r.Intro,
			Lang:    r.Lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// Info - Get detail information for specified repos.
func Info(c echo.Context) error {
	var (
		err  error
		req  infoReq
		resp infoResp
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	info, err := repos.Service.GetByID(&req.ID)
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	resp = infoResp{
		Avatar:  info.Avatar,
		Name:    info.Name,
		Image:   info.Image,
		Intro:   info.Intro,
		Lang:    info.Lang,
		Created: info.Created,
		Active:  info.Active,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}
