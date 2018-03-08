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
 *     Initial: 2018/03/06        Tong Yuehong
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
	Cart "github.com/fengyfei/gu/models/shop/cart"
)

type reqRemove struct {
	IDs []uint32 `json:"ids"`
}

// Add adds wares to cart.
func AddCart(c *server.Context) error {
	var (
		req  Cart.AddCartReq
		err  error
		conn orm.Connection
	)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysql.Pool.Get()", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error("JSONBody():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(&req)
	if err != nil {
		logger.Error("Validate():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userID := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(uint32)

	err = Cart.Service.Add(conn, userID, req.WareId, req.Count)
	if err != nil {
		logger.Error("Cart.Service.Add():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// GetByUser gets carts by userid.
func GetByUser(c *server.Context) error {
	var (
		items []Cart.Cart
		err   error
		conn  orm.Connection
	)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysql.Pool.Get():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	userID := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(uint32)

	items, err = Cart.Service.GetByUserID(conn, userID)
	if err != nil {
		logger.Error("Cart.Service.GetByUserID():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, items)
}

// Remove deletes a ware by id.
func Remove(c *server.Context) error {
	var (
		req  Cart.RemoveCartReq
		err  error
		conn orm.Connection
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

	err = c.Validate(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = Cart.Service.RemoveById(conn, req.Id)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// RemoveWhenOrder deletes wares by ids.
func RemoveWhenOrder(c *server.Context) error {
	var (
		req  reqRemove
		err  error
		conn orm.Connection
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

	err = Cart.Service.RemoveWhenOrder(conn, req.IDs)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
