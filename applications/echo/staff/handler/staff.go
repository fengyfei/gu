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
 *     Initial: 2017/11/01        Jia Chenhui
 */

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/applications/echo/staff/mysql"
	"github.com/fengyfei/gu/models/staff"
)

type (
	// loginReq - The request struct that login.
	loginReq struct {
		Name *string `json:"name" validate:"required,alphanum,min=6,max=30"`
		Pwd  *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
	}

	// createReq - The request struct that create staff information.
	createReq struct {
		Name     *string `json:"name" validate:"required,alphanum,min=6,max=30"`
		Pwd      *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		RealName *string `josn:"realname" validate:"required,alphanumunicode,min=2,max=20"`
		Mobile   *string `json:"mobile" validate:"required,numeric,len=11"`
		Email    *string `json:"email" validate:"required,email"`
		Male     bool    `json:"male"`
	}

	// modifyPwdReq - The request struct that modify staff password.
	modifyPwdReq struct {
		OldPwd *string `json:"oldpass" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		NewPwd *string `json:"newpass" validate:"required,printascii,excludesall=@-,min=6,max=30"`
	}

	// modifyMobileReq - The request struct that modify staff mobile.
	modifyMobileReq struct {
		Mobile *string `json:"mobile" validate:"required,numeric,len=11"`
	}

	// modifyActiveReq - The request struct that modify staff status.
	modifyActiveReq struct {
		Active *bool `json:"active"  validate:"required"`
	}

	// overviewResp - Overview of a staff.
	overviewResp struct {
		Id       int32
		RealName string
	}

	// infoResp - The more detail of one particular staff.
	infoResp struct {
		Id        int32
		Name      string
		RealName  string
		Mobile    string
		Email     string
		Male      bool
		CreatedAt time.Time
	}
)

// Login - Staff login.
func Login(c echo.Context) error {
	var (
		err error
		req loginReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	uid, err := staff.Service.Login(conn, req.Name, req.Pwd)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	key, token, err := core.NewToken(uid)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("token:", token)

	return c.JSON(http.StatusOK, map[string]interface{}{key: token})

}

// Create - Create staff information.
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

	if err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	err = staff.Service.Create(conn, req.Name, req.Pwd, req.RealName, req.Mobile, req.Email, req.Male)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ModifyPwd - Modify staff password.
func ModifyPwd(c echo.Context) error {
	var (
		err error
		req modifyPwdReq
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyPwd(conn, &uid, req.OldPwd, req.NewPwd); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ModifyMobile - Modify staff mobile.
func ModifyMobile(c echo.Context) error {
	var (
		err error
		req modifyMobileReq
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyMobile(conn, &uid, req.Mobile); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ModifyActive - Modify staff status.
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyActive(conn, &uid, req.Active); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// Dismiss - Staff dismission.
func Dismiss(c echo.Context) error {
	var err error

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	uid := core.UserID(c)
	if err = staff.Service.Dismiss(conn, &uid); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ContactList - Get a list of on-the-job staff.
func OverviewList(c echo.Context) error {
	var resp []overviewResp

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	slist, err := staff.Service.List(conn)
	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	for _, s := range slist {
		info := overviewResp{
			Id:       s.Id,
			RealName: s.RealName,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, resp)
}

// InfoList - Get a list of on-the-job staff details.
func InfoList(c echo.Context) error {
	var resp []infoResp

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	slist, err := staff.Service.List(conn)
	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	for _, s := range slist {
		info := infoResp{
			Id:        s.Id,
			Name:      s.Name,
			RealName:  s.RealName,
			Email:     s.Email,
			Male:      s.Male,
			CreatedAt: s.CreatedAt,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, resp)
}

// Info - Get detail information for specified staff.
func Info(c echo.Context) error {
	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}
	defer mysql.Pool.Release(conn)

	uid := core.UserID(c)
	info, err := staff.Service.GetByID(conn, &uid)
	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, *info)
}
