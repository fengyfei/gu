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
 *     Initial : 2018/03/06        Tong Yuehong
 */

package handler

import (
	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	models "github.com/fengyfei/gu/models/shop/address"
)

type (
	id struct {
		ID uint32 `json:"id"`
	}
)

// Add adds address.
func AddAddress(c *server.Context) error {
	var (
		add models.AddressData
	)

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

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	err = models.Service.Add(conn, userID, &add)
	if err != nil {
		logger.Error("Error in adding an address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// SetDefault sets default address.
func SetDefaultAddress(c *server.Context) error {
	var (
		id id
	)

	err := c.JSONBody(&id)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	userID := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(uint32)

	err = models.Service.SetDefault(conn, userID, id.ID)
	if err != nil {
		logger.Error("Error in setting the default address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Modify modify the address.
func ModifyAddress(c *server.Context) error {
	var (
		modify models.Modify
	)

	err := c.JSONBody(&modify)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&modify)
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

	err = models.Service.Modify(conn, &modify)
	if err != nil {
		logger.Error("Error in modifying the address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// Get gets address by userid.
func GetAddress(c *server.Context) error {
	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	userID := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(uint32)

	addr, err := models.Service.Get(conn, userID)
	if err != nil {
		logger.Error("Error in getting addresses:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, addr)
}

// Delete deletes address by id.
func DeleteAddress(c *server.Context) error {
	var (
		id id
	)

	err := c.JSONBody(&id)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.Delete(conn, id.ID)
	if err != nil {
		logger.Error("Error in deleting the address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
