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

package controller

import (
	"errors"
	"fmt"
	"time"

	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/polaris/mysql"
	"github.com/fengyfei/gu/applications/beego/polaris/utils"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/staff"
)

var (
	errPwdRepeat   = errors.New("the new password can't be the same as the old password")
	errPwdDisagree = errors.New("the new password and confirming password disagree")
)

// Staff - Staff associated handler.
type Staff struct {
	base.Controller
}

type (
	// loginReq - The request struct that login.
	loginReq struct {
		Name *string `json:"name" validate:"required,alphanum,min=2,max=30"`
		Pwd  *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
	}

	// createStaffReq - The request struct that create staff information.
	createStaffReq struct {
		Name     *string `json:"name" validate:"required,alphanum,min=2,max=30"`
		Pwd      *string `json:"pwd" validate:"required,printascii,excludesall=@-,min=6,max=30"`
		RealName *string `josn:"realname" validate:"required,alphanumunicode,min=2,max=20"`
		Mobile   *string `json:"mobile" validate:"required,numeric,len=11"`
		Email    *string `json:"email" validate:"required,email"`
		Male     bool    `json:"male"`
	}

	// modifyStaffReq - The request struct that modify staff information.
	modifyStaffReq struct {
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

	// activateStaffReq - The request struct that modify staff status.
	activateStaffReq struct {
		Id     int32 `json:"id" validate:"required"`
		Active bool  `json:"active"`
	}

	// dismissReq - The request struct that dismiss a staff.
	dismissReq struct {
		Id int32 `json:"id" validate:"required"`
	}

	// staffInfoReq - The request struct that get one staff detail information.
	staffInfoReq struct {
		Id int32 `json:"id" validate:"required"`
	}

	// staffInfoResp - The more detail of one particular staff.
	staffInfoResp struct {
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

	// roleListResp - The response struct that list all the roles of the specified staff.
	roleListResp struct {
		StaffId int32     `json:"staffid"`
		RoleId  int16     `json:"roleid"`
		Created time.Time `json:"created"`
	}
)

// Login - Staff login.
func (s *Staff) Login() {
	var (
		err error
		req loginReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	uid, err := staff.Service.Login(conn, req.Name, req.Pwd)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	_, token, err := utils.NewToken(uid)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
	}

	fmt.Println("token:", token)
	s.WriteStatusAndTokenJSON(constants.ErrSucceed, token)
}

// Create - Create staff information.
func (s *Staff) Create() {
	var (
		err error
		req createStaffReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	err = staff.Service.Create(conn, req.Name, req.Pwd, req.RealName, req.Mobile, req.Email, req.Male)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// Modify - Modify staff information.
func (s *Staff) Modify() {
	var (
		err error
		req modifyStaffReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	uid := s.Ctx.Request.Context().Value(utils.ClaimUID).(int32)

	if err = staff.Service.Modify(conn, uid, req.Name, req.Mobile, req.Email); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// ModifyPwd - Modify staff password.
func (s *Staff) ModifyPwd() {
	var (
		err error
		req modifyPwdReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if *req.NewPwd == *req.OldPwd {
		logger.Debug(errPwdRepeat.Error())
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if *req.NewPwd != *req.Confirm {
		logger.Debug(errPwdDisagree.Error())
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	uid := s.Ctx.Request.Context().Value(utils.ClaimUID).(int32)

	if err = staff.Service.ModifyPwd(conn, uid, req.OldPwd, req.NewPwd); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// ModifyMobile - Modify staff mobile.
func (s *Staff) ModifyMobile() {
	var (
		err error
		req modifyMobileReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	uid := s.Ctx.Request.Context().Value(utils.ClaimUID).(int32)

	if err = staff.Service.ModifyMobile(conn, uid, req.Mobile); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// ModifyActive - Modify staff status.
func (s *Staff) ModifyActive() {
	var (
		err error
		req activateStaffReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	if err = staff.Service.ModifyActive(conn, req.Id, req.Active); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// Dismiss - Dismissal of staff.
func (s *Staff) Dismiss() {
	var (
		err error
		req dismissReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	if err = staff.Service.Dismiss(conn, req.Id); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// List - Get a list of on-the-job staff details.
func (s *Staff) List() {
	var (
		err  error
		resp []staffInfoResp = make([]staffInfoResp, 0)
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	slist, err := staff.Service.List(conn)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	for _, s := range slist {
		info := staffInfoResp{
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

	s.WriteStatusAndDataJSON(constants.ErrSucceed, resp)
}

// Info - Get detail information for specified staff.
func (s *Staff) Info() {
	var (
		err error
		req staffInfoReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	info, err := staff.Service.GetByID(conn, req.Id)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	resp := staffInfoResp{
		Id:        info.Id,
		Name:      info.Name,
		RealName:  info.RealName,
		Mobile:    info.Mobile,
		Email:     info.Email,
		Male:      info.Male,
		CreatedAt: *info.CreatedAt,
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, resp)
}

// AddRole - Add a role to staff.
func (s *Staff) AddRole() {
	var (
		err error
		req addRoleReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	if err = staff.Service.AddRole(conn, req.StaffId, req.RoleId); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// RemoveRole - Remove role from staff.
func (s *Staff) RemoveRole() {
	var (
		err error
		req removeRoleReq
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	if err = staff.Service.RemoveRole(conn, req.StaffId, req.RoleId); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
}

// RoleList - List all the roles of the specified staff.
func (s *Staff) RoleList() {
	var (
		err  error
		req  roleListReq
		resp []roleListResp = make([]roleListResp, 0)
	)

	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	if err = s.Validate(&req); err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	rlist, err := staff.Service.AssociatedRoleList(conn, req.StaffId)
	if err != nil {
		logger.Error(err)
		s.WriteStatusAndDataJSON(constants.ErrMysql, nil)
	}

	for _, r := range rlist {
		relation := roleListResp{
			StaffId: r.StaffId,
			RoleId:  r.RoleId,
			Created: *r.Created,
		}

		resp = append(resp, relation)
	}

	s.WriteStatusAndDataJSON(constants.ErrSucceed, resp)
}
