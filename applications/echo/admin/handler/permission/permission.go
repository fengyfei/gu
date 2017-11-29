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
		URL    *string `json:"url" validate:"required,url"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// removeReq - The request struct that remove permission information.
	removeReq struct {
		URL    *string `json:"url" validate:"required,url"`
		RoleId int16   `json:"roleid" validate:"required"`
	}

	// listReq - The request struct that get a list of permission for specified URL.
	listReq struct {
		URL *string `json:"url" validate:"required,url"`
	}

	// infoResp - The response struct that represents role information of URL.
	infoResp struct {
		RoleId int16 `json:"roleid"`
	}
)

// Create - Create permission information.
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

	err = staff.Service.AddURLPermission(conn, req.URL, req.RoleId)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// Remove - Remove permission information.
func Remove(c echo.Context) error {
	var (
		err error
		req removeReq
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

	if err = staff.Service.RemoveURLPermission(conn, req.URL, req.RoleId); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// List - Get a list of permission for specified URL.
func List(c echo.Context) error {
	var (
		err   error
		req   listReq
		role  infoResp
		roles []infoResp
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

	resp, err := staff.Service.URLPermissions(conn, req.URL)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	for roleId := range resp {
		role = infoResp{
			RoleId: roleId,
		}
		roles = append(roles, role)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   roles,
	})
}
