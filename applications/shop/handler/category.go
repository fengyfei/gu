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
 *     Initial: 2018/02/03        Shi Ruitao
 */

package handler

import (
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/shop/category"
)

type (
	categoryAddReq struct {
		Name     string `json:"name" validate:"required,alphanumunicode,max=6"`
		Desc     string `json:"desc" validate:"required,max=50"`
		ParentID uint64 `json:"parentid" validate:"required"`
	}

	subCategoryReq struct {
		PID uint64 `json:"pid" validate:"required"`
	}
)

// add new category
func AddCategory(c *server.Context) error {
	var (
		err    error
		addReq categoryAddReq
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&addReq)
	if err != nil {
		logger.Error("JSONBody():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(addReq)
	if err != nil {
		logger.Error("Validate():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = category.Service.AddCategory(conn, &addReq.Name, &addReq.Desc, &addReq.ParentID)
	if err != nil {
		logger.Error("category.Service.AddCategory():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// get all parent categories
func GetMainCategories(c *server.Context) error {
	var (
		pid uint64 = 0
		err error
		res []category.Category
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("mysql.Pool.Get()", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = category.Service.GetCategory(conn, pid)
	if err != nil {
		logger.Error("category.Service.GetCategory()", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get categories of the specified pid
func GetSubCategories(c *server.Context) error {
	var (
		err    error
		pidReq subCategoryReq
		res    []category.Category
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error("mysql.Pool.Get():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&pidReq)
	if err != nil {
		logger.Error("JSONBody():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(pidReq)
	if err != nil {
		logger.Error("Validate():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err = category.Service.GetCategory(conn, pidReq.PID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}
