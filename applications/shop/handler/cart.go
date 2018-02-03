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
	"github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	Cart "github.com/fengyfei/gu/models/shop/cart"
)

// Add commodity
func Add(c *server.Context) error {
	var (
		req   Cart.AddCartReq
		err   error
		token *jwt.Token
		conn  orm.Connection
	)

	token, err = util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint64(claims[util.UserID].(float64))

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

	err = Cart.Service.Add(conn, userId, req.WareId, req.Count)
	if err != nil {
		logger.Error("Cart.Service.Add():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func GetByUser(c *server.Context) error {
	var (
		items []Cart.CartItem
		err   error
		conn  orm.Connection
		token *jwt.Token
	)

	token, err = util.Parse(c)
	if err != nil {
		logger.Error("Error in parsing token:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims[util.UserID].(float64))

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("mysql.Pool.Get():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	items, err = Cart.Service.GetByUserID(conn, userId)
	if err != nil {
		logger.Error("Cart.Service.GetByUserID():", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, items)
}

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
