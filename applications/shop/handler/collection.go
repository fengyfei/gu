/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the 'Software'), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/03/28      Shi Ruitao
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
	"github.com/fengyfei/gu/libs/orm"
	Collection "github.com/fengyfei/gu/models/shop/collection"
)

// AddColl adds wares into someone's collection.
func AddColl(c *server.Context) error {
	var (
		err  error
		conn orm.Connection
		add  struct {
			WareID uint32 `json:"ware_id" validate:"required"`
		}
	)

	err = c.JSONBody(&add)
	if err != nil {
		logger.Error("Error in JSONBody:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&add)
	if err != nil {
		logger.Error("Validate():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	err = Collection.Service.Add(conn, userID, add.WareID)
	if err != nil {
		logger.Error("Error in modifying the collection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// GetByUserID gets collections by userID.
func GetByUserID(c *server.Context) error {
	var (
		err  error
		conn orm.Connection
	)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	userID := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	collections, err := Collection.Service.GetByUserID(conn, userID)
	if err != nil {
		logger.Error("Error in modifying the collection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, collections)
}

// RemoveMany deletes collectios by ids.
func RemoveColl(c *server.Context) error {
	var (
		err  error
		conn orm.Connection
		req  struct {
			IDs []uint32 `json:"ids"`
		}
	)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = Collection.Service.Remove(conn, req.IDs)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
