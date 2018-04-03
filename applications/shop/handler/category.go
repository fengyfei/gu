/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/04/02        Shi Ruitao
 */

package handler

import (
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	category "github.com/fengyfei/gu/models/shop/category"
)

var (
	ErrNotAdmin = errors.New("non-administrators can not operate.")
)

type (
	pid struct {
		Id uint32 `json:"id"`
	}
)

// Add adds a new category.
func AddCategory(c *server.Context) error {
	var (
		add struct {
			Category string `json:"category" validate:"required,min=1,max=12"`
			ParentID uint32 `json:"parent_id"`
		}
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("Permission denied:", ErrNotAdmin)
		return core.WriteStatusAndDataJSON(c, constants.ErrPermission, ErrNotAdmin)
	}

	err := c.JSONBody(&add)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&add)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = category.Service.Add(conn, add.Category, add.ParentID)
	if err != nil {
		logger.Error("Error in add a category:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// GetMainCategory gets all the parentid.
func GetMainCategories(c *server.Context) error {
	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysqlErr:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err := category.Service.GetMainCategory(conn)
	if err != nil {
		logger.Error("Get main category error:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// GetSubCategory gets subcategory which parentid is pid.
func GetSubCategories(c *server.Context) error {
	var (
		err error
		pid pid
		res []category.Category
	)

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysql.Pool.Get():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&pid)
	if err != nil {
		logger.Error("jsonErr", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err = category.Service.GetSubCategory(conn, pid.Id)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// Delete deletes the category.
func Delete(c *server.Context) error {
	var (
		err error
		pid pid
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("Permission denied:", ErrNotAdmin)
		return core.WriteStatusAndDataJSON(c, constants.ErrPermission, ErrNotAdmin)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysqlErr:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&pid)
	if err != nil {
		logger.Error("jsonErr:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = category.Service.Delete(conn, pid.Id)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// ModifyCategory update the category's information.
func ModifyCategory(c *server.Context) error {
	var (
		err    error
		modify struct {
			Id       uint32 `json:"id"`
			Category string `json:"category" validate:"required,min=2,max=12"`
			ParentID uint32 `json:"parent_id"`
		}
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("Permission denied:", ErrNotAdmin)
		return core.WriteStatusAndDataJSON(c, constants.ErrPermission, ErrNotAdmin)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysqlErr", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&modify)
	if err != nil {
		logger.Error("jsonErr:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = category.Service.Modify(conn, modify.Id, modify.Category, modify.ParentID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
