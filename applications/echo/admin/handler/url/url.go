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
 *     Initial: 2017/11/13        Jia Chenhui
 */

package url

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/models/url"
)

type (
	// createReq - The request struct that create url information.
	createReq struct {
		Content *string `json:"content" validate:"required,alphanum,min=3,max=128"`
	}

	// modifyReq - The request struct that modify url information.
	modifyReq struct {
		Id      *int16  `json:"id" validate:"required,numeric"`
		Content *string `json:"content" validate:"required,alphanum,min=3,max=128"`
	}

	// modifyActiveReq - The request struct that modify url status.
	modifyActiveReq struct {
		Id     *int16 `json:"id" validate:"required,numeric"`
		Active *bool  `json:"active"  validate:"required"`
	}

	// infoReq - The request struct for get detail of specified URL.
	infoReq struct {
		Id *int16 `json:"id" validate:"required,numeric"`
	}

	// infoResp - The detail information for url.
	infoResp struct {
		Id      int16
		Content string
		Created time.Time
	}
)

// Create - Create url information.
func Create(c echo.Context) error {
	var (
		err error
		req createReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	err = url.Service.Create(conn, req.Content)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// Modify - Modify url information.
func Modify(c echo.Context) error {
	var (
		err error
		req modifyReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	if err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = url.Service.Modify(conn, req.Id, req.Content); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ModifyActive - Modify url status.
func ModifyActive(c echo.Context) error {
	var (
		err error
		req modifyActiveReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	if err = url.Service.ModifyActive(conn, req.Id, req.Active); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// InfoList - Get a list of active url details.
func InfoList(c echo.Context) error {
	var resp []infoResp

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	ulist, err := url.Service.List(conn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	for _, u := range ulist {
		info := infoResp{
			Id:      u.Id,
			Content: u.Content,
			Created: *u.Created,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, resp)
}

// Info - Get detail information for specified url.
func Info(c echo.Context) error {
	var (
		err error
		req infoReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	info, err := url.Service.GetByID(conn, req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, *info)
}
