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
 *     Initial: 2018/03/08        Shi Ruitao
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
	Order "github.com/fengyfei/gu/models/shop/order"
)

type (
	confirmOrderReq struct {
		ID uint64 `json:"id"`
	}
)

func CreateOrder(c *server.Context) error {
	var (
		req     Order.CreateReq
		err     error
		IP      string
		userId  uint32
		conn    orm.Connection
		signStr string
	)

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	IP = c.Request().RemoteAddr

	err = c.Validate(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	userId = uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	signStr, err = Order.Service.OrderByWechat(conn, userId, IP, &req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrWechatPay, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, signStr)
}

func ConfirmOrder(c *server.Context) error {
	var (
		req  confirmOrderReq
		conn orm.Connection
		err  error
	)

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

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		return core.WriteStatusAndDataJSON(c, constants.ErrPermission, nil)
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = Order.Service.ChangeStateByOne(conn, req.ID, Order.StatusConfirmed)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func GetUserOrder(c *server.Context) error {
	var (
		conn   orm.Connection
		orders *[]Order.Order
		err    error
	)

	userId := uint32(c.Request().Context().Value("user").(jwtgo.MapClaims)[util.UserID].(float64))

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	orders, err = Order.Service.GetUserOrder(conn, userId)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, orders)
}
