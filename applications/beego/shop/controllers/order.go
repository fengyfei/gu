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
 *     Initial: 2017/11/25        ShiChao
 */

package controllers

import (
	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/libs/constants"
	"encoding/json"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/libs/orm"
	Order "github.com/fengyfei/gu/models/shop/order"
)

type (
	OrderController struct {
		base.Controller
	}

	createReq struct {
		Orders     []Order.OrderItem `json:"orders" validate:"required"`
		ReceiveWay int8              `json:"receiveWay" validate:"required"`
	}

	confirmOrderReq struct {
		ID     uint `json:"id"`
	}
)

func (this *OrderController) CreateOrder() {
	var (
		req     createReq
		err     error
		IP      string
		userId  uint
		conn    orm.Connection
		signStr string
	)

	userId = this.Ctx.Request.Context().Value("userId").(uint)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	IP = this.Ctx.Input.IP()

	err = this.Validate(&req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	signStr, err = Order.Service.OrderByWechat(conn, userId, IP, req.ReceiveWay, req.Orders)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatPay}

		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespPaySign: signStr}

finish:
	this.ServeJSON(true)
}

func (this *OrderController) ConfirmOrder() {

	var (
		req  confirmOrderReq
		conn orm.Connection
		err  error
	)

	isAdmin := this.Ctx.Request.Context().Value("isAdmin").(bool)
	if !isAdmin {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrPermission}
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = this.Validate(&req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = Order.Service.ChangeStateByOne(conn, req.ID, Order.StatusConfirmed)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}

func (this *OrderController) GetUserOrder() {
	var (
		conn   orm.Connection
		orders *[]Order.Order
		err    error
	)

	userId := this.Ctx.Request.Context().Value("userId").(uint)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	orders, err = Order.Service.GetUserOrder(conn, userId)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespKeyData: orders}

finish:
	this.ServeJSON(true)
}

