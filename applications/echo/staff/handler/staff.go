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
		HireAt   *string `json:"hireat" validate:"required"`
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
		Active bool `json:"active"  validate:"required"`
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

	uid, err := staff.Service.Login(req.Name, req.Pwd)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	key, token, err := core.NewToken(uid)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("token", token)

	return c.JSON(http.StatusOK, map[string]interface{}{key: token})

}

// Create - Create staff information.
func Create(c echo.Context) error {
	var (
		err    error
		hireAt time.Time
		req    createReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	hireAt, err = time.Parse("2006-01-02", *req.HireAt)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	err = staff.Service.Create(req.Name, req.Pwd, req.RealName, req.Mobile, req.Email, hireAt, req.Male)
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyPwd(&uid, req.OldPwd, req.NewPwd); err != nil {
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyMobile(&uid, req.Mobile); err != nil {
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyActive(&uid); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// Dismiss - Staff dismission.
func Dismiss(c echo.Context) error {
	var err error

	uid := core.UserID(c)
	if err = staff.Service.Dismiss(&uid); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// CheckIn - Staff check in.
func CheckIn(c echo.Context) error {
	var err error

	uid := core.UserID(c)
	_, ok, err := staff.Service.IsRegistered(&uid)
	if err != nil {
		if err == xorm.ErrNotExist {

			if err = staff.Service.Register(&uid); err != nil {
				return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
			}

			goto finish
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	if ok {
		return core.NewErrorWithMsg(http.StatusForbidden, "staff already registered")
	}

	if err = staff.Service.RegisterAgain(&uid); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

finish:
	return c.JSON(http.StatusOK, nil)
}

// CheckOut - Staff check out.
func CheckOut(c echo.Context) error {
	var err error

	uid := core.UserID(c)
	r, ok, err := staff.Service.IsRegistered(&uid)
	if err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	if !ok {
		return core.NewErrorWithMsg(http.StatusInternalServerError, "not register yet")
	}

	if err = staff.Service.LeaveOffice(&uid, r); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

// ContactList - Get a list of on-the-job staff.
func ContactList(c echo.Context) error {
	list, err := staff.Service.OverviewList()

	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

// InfoList - Get a list of on-the-job staff details.
func InfoList(c echo.Context) error {
	list, err := staff.Service.InfoList()

	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

// Info - Get detail information for specified staff.
func Info(c echo.Context) error {
	uid := core.UserID(c)

	info, err := staff.Service.GetByID(&uid)
	if err != nil {
		if err == xorm.ErrNotExist {
			return core.NewErrorWithMsg(http.StatusNotFound, err.Error())
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}
