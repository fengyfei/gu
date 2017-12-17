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
 *     Initial: 2017/11/09        Jia Chenhui
 */

package controller

import (
	"time"

	json "github.com/json-iterator/go"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/polaris/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/staff"
)

// Role - Role associated handler.
type Role struct {
	base.Controller
}

type (
	// createRoleReq - The request struct that create role information.
	createRoleReq struct {
		Name  *string `json:"name" validate:"required,alpha,min=4,max=64"`
		Intro *string `json:"intro" validate:"required,alphanumunicode,min=4,max=256"`
	}

	// modifyRoleReq - The request struct that modify role information.
	modifyRoleReq struct {
		Id    int16   `json:"id" validate:"required"`
		Name  *string `json:"name" validate:"required,alpha,min=4,max=64"`
		Intro *string `json:"intro" validate:"required,alphanumunicode,min=4,max=256"`
	}

	// activateRoleReq - The request struct that modify role status.
	activateRoleReq struct {
		Id     int16 `json:"id" validate:"required"`
		Active bool  `json:"active"`
	}

	// roleInfoReq - The request struct for get detail of specified role.
	roleInfoReq struct {
		Id int16 `json:"id" validate:"required"`
	}

	// roleInfoResp - The detail information for role.
	roleInfoResp struct {
		Id      int16     `json:"id"`
		Name    string    `json:"name"`
		Intro   string    `json:"intro"`
		Active  bool      `json:"active"`
		Created time.Time `json:"created"`
	}
)

// Create - Create role information.
func (r *Role) Create() {
	var (
		err  error
		req  createRoleReq
		conn orm.Connection
	)

	if err = json.Unmarshal(r.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	if err = r.Validate(&req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	if err = staff.Service.CreateRole(conn, req.Name, req.Intro); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	r.WriteJSON(constants.ErrSucceed)

finish:
	r.ServeJSON(true)
}

// Modify - Modify role information.
func (r *Role) Modify() {
	var (
		err  error
		req  modifyRoleReq
		conn orm.Connection
	)

	if err = json.Unmarshal(r.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	if err = r.Validate(&req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	if err = staff.Service.ModifyRole(conn, req.Id, req.Name, req.Intro); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	r.WriteJSON(constants.ErrSucceed)

finish:
	r.ServeJSON(true)
}

// ModifyActive - Modify role status.
func (r *Role) ModifyActive() {
	var (
		err  error
		req  activateRoleReq
		conn orm.Connection
	)

	if err = json.Unmarshal(r.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	if err = r.Validate(&req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	if err = staff.Service.ModifyRoleActive(conn, req.Id, req.Active); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	r.WriteJSON(constants.ErrSucceed)

finish:
	r.ServeJSON(true)
}

// List - Get a list of active role details.
func (r *Role) List() {
	var (
		err   error
		conn  orm.Connection
		rlist []staff.Role
		resp  []roleInfoResp = make([]roleInfoResp, 0)
	)

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	rlist, err = staff.Service.RoleList(conn)
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	for _, r := range rlist {
		info := roleInfoResp{
			Id:      r.Id,
			Name:    r.Name,
			Intro:   r.Intro,
			Active:  r.Active,
			Created: *r.Created,
		}

		resp = append(resp, info)
	}

	r.WriteJSON(constants.ErrSucceed, resp)

finish:
	r.ServeJSON(true)
}

// Info - Get detail information for specified role.
func (r *Role) Info() {
	var (
		err  error
		req  roleInfoReq
		info *staff.Role
		conn orm.Connection
	)

	if err = json.Unmarshal(r.Ctx.Input.RequestBody, &req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	if err = r.Validate(&req); err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrInvalidParam)

		goto finish
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	info, err = staff.Service.GetRoleByID(conn, req.Id)
	if err != nil {
		logger.Error(err)
		r.WriteJSON(constants.ErrMysql)

		goto finish
	}

	r.WriteJSON(constants.ErrSucceed, *info)

finish:
	r.ServeJSON(true)
}
