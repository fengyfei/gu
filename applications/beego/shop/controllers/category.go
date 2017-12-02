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
 *     Initial: 2017/11/07        Wang RiYu
 */

package controllers

import (
  "encoding/json"
  "github.com/fengyfei/gu/applications/beego/base"
  "github.com/fengyfei/gu/libs/constants"
  "github.com/fengyfei/gu/libs/logger"
  "github.com/fengyfei/gu/applications/beego/shop/mysql"
  "github.com/fengyfei/gu/models/shop/category"
)

type (
  CategoryController struct {
    base.Controller
  }

  categoryAddReq struct {
    Name     string `json:"name"`
    Desc     string `json:"desc"`
    ParentID uint   `json:"parentId"`
  }

  subCategoryReq struct {
    PID uint `json:"pid"`
  }
)

// add new category
func (this *CategoryController) AddCategory() {
  var (
    err    error
    addReq categoryAddReq
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

  err = category.Service.AddCategory(conn, &addReq.Name, &addReq.Desc, &addReq.ParentID)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}
  logger.Info("create category", addReq.Name, "success")

finish:
  this.ServeJSON(true)
}

// get all parent categories
func (this *CategoryController) GetMainCategories() {
  var (
    pid uint = 0
    err error
    res []category.Category
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  res, err = category.Service.GetCategory(conn, pid)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = res

finish:
  this.ServeJSON(true)
}

// get categories of the specified pid
func (this *CategoryController) GetSubCategories() {
  var (
    err    error
    pidReq subCategoryReq
    res    []category.Category
  )

  conn, err := mysql.Pool.Get()
  defer mysql.Pool.Release(conn)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }

  err = json.Unmarshal(this.Ctx.Input.RequestBody, &pidReq)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

    goto finish
  }

  res, err = category.Service.GetCategory(conn, pidReq.PID)
  if err != nil {
    logger.Error(err)
    this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

    goto finish
  }
  this.Data["json"] = res

finish:
  this.ServeJSON(true)
}
