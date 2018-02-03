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
 *     Initial: 2018/02/01        Shi Ruitao
 */

package handler

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	models "github.com/fengyfei/gu/models/shop/address"
)

// Add address
func AddAddress(c *server.Context) error {
	var add models.Add

	err := c.JSONBody(&add)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&add)
	if err != nil || !util.IsValidPhone(add.Phone) {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint64(claims[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.Add(conn, userID, &add)
	if err != nil {
		logger.Error("Error in adding an address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func SetDefaultAddress(c *server.Context) error {
	var set models.SetDefault

	err := c.JSONBody(&set)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&set)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.SetDefault(conn, userID, set.ID)
	if err != nil {
		logger.Error("Error in setting the default address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func ModifyAddress(c *server.Context) error {
	var modify models.Modify

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

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.Modify(conn, userID, &modify)
	if err != nil {
		logger.Error("Error in modifying the address:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func GetAddress(c *server.Context) error {
	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	addr, err := models.Service.Get(conn, userID)
	if err != nil {
		logger.Error("Error in getting addresses:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, addr)
}

func DeleteAddress(c *server.Context) error {
	var delete models.Delete

	err := c.JSONBody(&delete)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&delete)
	if err != nil {
		logger.Error("Invalid parameters:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	token, err := util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims[util.UserID].(float64))

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = models.Service.Delete(conn, userID, delete.ID)
	if err != nil {
		logger.Error("Error in deleting:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
