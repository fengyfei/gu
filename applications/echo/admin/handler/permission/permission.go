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

package permission

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/models/staff"
)

type (
	// createReq - The request struct that create permission information.
	createReq struct {
		URL    *string `json:"url" validate:"required,printascii,contains=/"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// removeReq - The request struct that remove permission information.
	removeReq struct {
		URL    *string `json:"url" validate:"required,printascii,contains=/"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// infoResp - The response struct that represents role information of URL.
	infoResp struct {
		URL    string `json:"url"`
		RoleId int16  `json:"roleid"`
	}
)

// Create - Create permission information.
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

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	if err = staff.Service.AddURLPermission(conn, req.URL, req.RoleId); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// Remove - Remove permission information.
func Remove(c echo.Context) error {
	var (
		err error
		req removeReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	if err = staff.Service.RemoveURLPermission(conn, req.URL, req.RoleId); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// List - Get a list of permission for specified URL.
func List(c echo.Context) error {
	var (
		err        error
		permission infoResp
		resp       []infoResp = make([]infoResp, 0)
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	plist, err := staff.Service.Permissions(conn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	for _, p := range plist {
		permission = infoResp{
			URL:    p.URL,
			RoleId: p.RoleId,
		}

		resp = append(resp, permission)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}
