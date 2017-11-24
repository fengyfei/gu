/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/11/23        Wang RiYu
 */

package controllers

import (
  "encoding/json"
  "github.com/fengyfei/gu/applications/beego/base"
  "github.com/fengyfei/gu/libs/constants"
  "github.com/fengyfei/gu/libs/logger"
  "github.com/fengyfei/gu/applications/beego/shop/mysql"
  "github.com/fengyfei/gu/models/shop/ware"
)

type (
  WareController struct {
    base.Controller
  }

  categoryReq struct {
    CID uint `json:"cid"`
  }
)

// add new ware
func (this *WareController) CreateWare() {
  var (
    err    error
    addReq ware.Ware
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  err = json.Unmarshal(this.Ctx.Input.RequestBody, &addReq)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

    goto finish
  }

  err = ware.Service.CreateWare(conn, addReq)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
  this.ServeJSON(true)
}

// get all wares
func (this *WareController) GetWareList() {
  var (
    err error
    res []ware.Ware
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  res, err = ware.Service.GetWareList(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = res

finish:
  this.ServeJSON(true)
}

// get ware by categoryID
func (this *WareController) GetWareByCategory() {
  var (
    err    error
    cidReq categoryReq
    res    []ware.Ware
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  err = json.Unmarshal(this.Ctx.Input.RequestBody, &cidReq)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

    goto finish
  }

  res, err = ware.Service.GetWareList(conn, cidReq.CID)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = res

finish:
  this.ServeJSON(true)
}

// get promotion wares
func (this *WareController) GetPromotion() {
  var (
    err error
    res []ware.Ware
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  res, err = ware.Service.GetPromotionList(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = res

finish:
  this.ServeJSON(true)
}