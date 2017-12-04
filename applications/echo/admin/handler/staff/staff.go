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

package staff

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/models/staff"
)

var (
	errPwdRepeat   = errors.New("the new password can't be the same as the old password")
	errPwdDisagree = errors.New("the new password and confirming password disagree")
)

type (
	// loginReq - The request struct that login.
	loginReq struct {
		Name *string `json:"name" validate:"required,alphanum,min=2,max=30"`
		Pwd  *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
	}

	// createReq - The request struct that create staff information.
	createReq struct {
		Name     *string `json:"name" validate:"required,alphanum,min=2,max=30"`
		Pwd      *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		RealName *string `josn:"realname" validate:"required,alphanumunicode,min=2,max=20"`
		Mobile   *string `json:"mobile" validate:"required,numeric,len=11"`
		Email    *string `json:"email" validate:"required,email"`
		Male     bool    `json:"male"`
	}

	// modifyReq - The request struct that modify staff information.
	modifyReq struct {
		Name   *string `json:"name" validate:"required,alphanum,min=2,max=30"`
		Mobile *string `json:"mobile" validate:"required,numeric,len=11"`
		Email  *string `json:"email" validate:"required,email"`
	}

	// modifyPwdReq - The request struct that modify staff password.
	modifyPwdReq struct {
		OldPwd  *string `json:"oldpwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		NewPwd  *string `json:"newpwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		Confirm *string `json:"confirm" validate:"required,printascii,excludesall=@-,min=6,max=30"`
	}

	// modifyMobileReq - The request struct that modify staff mobile.
	modifyMobileReq struct {
		Mobile *string `json:"mobile" validate:"required,numeric,len=11"`
	}

	// activateReq - The request struct that modify staff status.
	activateReq struct {
		Id     int32 `json:"id" validate:"required"`
		Active bool  `json:"active"`
	}

	// dismissReq - The request struct that dismiss a staff.
	dismissReq struct {
		Id int32 `json:"id" validate:"required"`
	}

	// infoReq - The request struct that get one staff detail information.
	infoReq struct {
		Id int32 `json:"id" validate:"required"`
	}

	// infoResp - The more detail of one particular staff.
	infoResp struct {
		Id        int32     `json:"id"`
		Name      string    `json:"name"`
		RealName  string    `json:"realname"`
		Mobile    string    `json:"mobile"`
		Email     string    `json:"email"`
		Male      bool      `json:"male"`
		Resigned  bool      `json:"resigned"`
		CreatedAt time.Time `json:"createdat"`
	}

	// addRoleReq - The request struct that add role to staff.
	addRoleReq struct {
		StaffId int32 `json:"staffid" validate:"required"`
		RoleId  int16 `json:"roleid" validate:"required"`
	}

	// removeRoleReq - The request struct that remove role from staff.
	removeRoleReq struct {
		StaffId int32 `json:"staffid" validate:"required"`
		RoleId  int16 `json:"roleid" validate:"required"`
	}

	// roleListReq - The request struct that list all the roles of the specified staff.
	roleListReq struct {
		StaffId int32 `json:"staffid" validate:"required"`
	}

	// roleListResp - The response struct that list all the roles of the specifide staff.
	roleListResp struct {
		StaffId int32     `json:"staffid"`
		RoleId  int16     `json:"roleid"`
		Created time.Time `json:"created"`
	}
)

// Login - Staff login.
func Login(c echo.Context) error {
	var (
		err error
		req loginReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	uid, err := staff.Service.Login(conn, req.Name, req.Pwd)
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	_, token, err := core.NewToken(uid)
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrInternalServerError, err.Error())
	}

	fmt.Println("token:", token)

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyToken:  token,
	})

}

// Create - Create staff information.
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

	err = staff.Service.Create(conn, req.Name, req.Pwd, req.RealName, req.Mobile, req.Email, req.Male)
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// Modify - Modify staff information.
func Modify(c echo.Context) error {
	var (
		err error
		req modifyReq
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

	uid := core.UserID(c)
	if err = staff.Service.Modify(conn, uid, req.Name, req.Mobile, req.Email); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// ModifyPwd - Modify staff password.
func ModifyPwd(c echo.Context) error {
	var (
		err error
		req modifyPwdReq
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, err.Error())
	}

	if *req.NewPwd == *req.OldPwd {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, errPwdRepeat.Error())
	}

	if *req.NewPwd != *req.Confirm {
		return core.NewErrorWithMsg(constants.ErrInvalidParam, errPwdDisagree.Error())
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	uid := core.UserID(c)
	if err = staff.Service.ModifyPwd(conn, uid, req.OldPwd, req.NewPwd); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// ModifyMobile - Modify staff mobile.
func ModifyMobile(c echo.Context) error {
	var (
		err error
		req modifyMobileReq
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

	uid := core.UserID(c)
	if err = staff.Service.ModifyMobile(conn, uid, req.Mobile); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// ModifyActive - Modify staff status.
func ModifyActive(c echo.Context) error {
	var (
		err error
		req activateReq
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

	if err = staff.Service.ModifyActive(conn, req.Id, req.Active); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// Dismiss - Staff dismission.
func Dismiss(c echo.Context) error {
	var (
		err error
		req dismissReq
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

	if err = staff.Service.Dismiss(conn, req.Id); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// List - Get a list of on-the-job staff details.
func List(c echo.Context) error {
	var resp []infoResp = make([]infoResp, 0)

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	slist, err := staff.Service.List(conn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	for _, s := range slist {
		info := infoResp{
			Id:        s.Id,
			Name:      s.Name,
			RealName:  s.RealName,
			Mobile:    s.Mobile,
			Email:     s.Email,
			Male:      s.Male,
			Resigned:  s.Resigned,
			CreatedAt: *s.CreatedAt,
		}

		resp = append(resp, info)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// Info - Get detail information for specified staff.
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

	conn, err := mysql.Pool.Get()
	if err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}
	defer mysql.Pool.Release(conn)

	info, err := staff.Service.GetByID(conn, req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	resp = infoResp{
		Id:        info.Id,
		Name:      info.Name,
		RealName:  info.RealName,
		Mobile:    info.Mobile,
		Email:     info.Email,
		Male:      info.Male,
		CreatedAt: *info.CreatedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

// AddRole - Add a role to staff.
func AddRole(c echo.Context) error {
	var (
		err error
		req addRoleReq
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

	if err = staff.Service.AddRole(conn, req.StaffId, req.RoleId); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// RemoveRole - Remove role from staff.
func RemoveRole(c echo.Context) error {
	var (
		err error
		req removeRoleReq
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

	if err = staff.Service.RemoveRole(conn, req.StaffId, req.RoleId); err != nil {
		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
	})
}

// RoleList - List all the roles of the specified staff.
func RoleList(c echo.Context) error {
	var (
		err      error
		req      roleListReq
		relation roleListResp
		resp     []roleListResp = make([]roleListResp, 0)
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

	rlist, err := staff.Service.AssociatedRoleList(conn, req.StaffId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
		}

		return core.NewErrorWithMsg(constants.ErrMysql, err.Error())
	}

	for _, r := range rlist {
		relation = roleListResp{
			StaffId: r.StaffId,
			RoleId:  r.RoleId,
			Created: *r.Created,
		}

		resp = append(resp, relation)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}
