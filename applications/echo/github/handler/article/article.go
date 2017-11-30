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
 *     Initial: 2017/11/22        Jia Chenhui
 */

package article

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/models/github/article"
)

type (
	// createReq - The request struct that create article information.
	createReq struct {
		Title  *string `json:"title" validate:"required,min=1,max=256"`
		URL    *string `json:"url" validate:"required,url"`
		Source *string `json:"source"`
	}

	// modifyActiveReq - The request struct that modify article status.
	modifyActiveReq struct {
		ID     string `json:"id" validate:"required,alphanum,len=24"`
		Active bool   `json:"active"`
	}

	// infoReq - The request struct that get one article detail information.
	infoReq struct {
		ID string `json:"id" validate:"required,alphanum,len=24"`
	}

	// infoResp - The more detail of article.
	infoResp struct {
		Title   string    `json:"title"`
		URL     string    `json:"url"`
		Source  string    `json:"source"`
		Active  bool      `json:"active"`
		Created time.Time `json:"created"`
	}
)

// Create - Create article information.
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

	id, err := article.Service.Create(req.Title, req.URL, req.Source)
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyID:     id,
	})
}

// ModifyActive - Modify article status.
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

	if err = article.Service.ModifyActive(&req.ID, req.Active); err != nil {
		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// List - Get all the articles.
func List(c echo.Context) error {
	var resp []infoResp = make([]infoResp, 0)

	alist, err := article.Service.List()
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	for _, a := range alist {
		info := infoResp{
			Title:   a.Title,
			URL:     a.URL,
			Source:  a.Source,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// ActiveList - Get all the active articles.
func ActiveList(c echo.Context) error {
	var resp []infoResp = make([]infoResp, 0)

	alist, err := article.Service.ActiveList()
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	for _, a := range alist {
		info := infoResp{
			Title:   a.Title,
			URL:     a.URL,
			Source:  a.Source,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// Info - Get detail information for specified article.
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

	info, err := article.Service.GetByID(&req.ID)
	if err != nil {
		if err == mgo.ErrNotFound {
			return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMongoDB, err.Error())
	}

	resp = infoResp{
		Title:   info.Title,
		URL:     info.URL,
		Source:  info.Source,
		Active:  info.Active,
		Created: info.Created,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}
