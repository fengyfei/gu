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

package handler

//import (
//  "encoding/json"
//  "github.com/fengyfei/gu/applications/beego/base"
//  "github.com/fengyfei/gu/libs/constants"
//  "github.com/fengyfei/gu/libs/logger"
//  "github.com/fengyfei/gu/applications/beego/shop/mysql"
//  "github.com/fengyfei/gu/models/shop/ware"
//  "github.com/fengyfei/gu/libs/orm"
//  "github.com/fengyfei/gu/applications/beego/shop/util"
//  "errors"
//)
//
//type (
//  WareController struct {
//    base.Controller
//  }
//
//  categoryReq struct {
//    CID uint `json:"cid" validate:"required"`
//  }
//
//  homeReq struct {
//    LastID int `json:"lastId"`
//  }
//
//  detailReq struct {
//    ID uint `json:"id" validate:"required"`
//  }
//
//  recommendReq struct {
//    UserID uint `json:"id" validate:"required"`
//  }
//
//  ChangeStatusReq struct {
//    IDs    []uint `json:"ids" validate:"required,min=1"`
//    Status int8   `json:"status" validate:"required,eq=-1|eq=1|eq=2|eq=3"`
//  }
//)
//
//// add new ware
//func (this *WareController) CreateWare() {
//  var (
//    err    error
//    addReq ware.Ware
//    conn   orm.Connection
//  )
//
//  conn, err = mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &addReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(addReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  if len(addReq.Avatar) > 0 {
//    addReq.Avatar, err = util.SavePicture(addReq.Avatar, "ware/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//  if len(addReq.Image) > 0 {
//    addReq.Image, err = util.SavePicture(addReq.Image, "ware/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//  if len(addReq.DetailPic) > 0 {
//    addReq.DetailPic, err = util.SavePicture(addReq.DetailPic, "wareIntro/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//
//  err = ware.Service.CreateWare(conn, addReq)
//  if err != nil {
//    logger.Error(err)
//    if (len(addReq.Avatar) > 0 && !util.DeletePicture(addReq.Avatar)) ||
//       (len(addReq.Image) > 0 && !util.DeletePicture(addReq.Image)) ||
//       (len(addReq.DetailPic) > 0 && !util.DeletePicture(addReq.DetailPic)) {
//      logger.Error(errors.New("create ware failed and delete it's pictures go wrong, please delete picture manually"))
//    }
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  logger.Info("create ware", addReq.Name, "success")
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
//}
//
//// get all wares
//func (this *WareController) GetAllWare() {
//  var (
//    err error
//    res []ware.Ware
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  res, err = ware.Service.GetAllWare(conn)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// get ware by categoryID
//func (this *WareController) GetWareByCategory() {
//  var (
//    err    error
//    cidReq categoryReq
//    res    []ware.BriefInfo
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &cidReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(cidReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  res, err = ware.Service.GetByCID(conn, cidReq.CID)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// get promotion wares
//func (this *WareController) GetPromotion() {
//  var (
//    err error
//    res []ware.BriefInfo
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  res, err = ware.Service.GetPromotionList(conn)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// update ware info
//func (this *WareController) UpdateWithID() {
//  var (
//    err error
//    req ware.UpdateReq
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  if len(req.Avatar) > 0 {
//    req.Avatar, err = util.SavePicture(req.Avatar, "ware/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//  if len(req.Image) > 0 {
//    req.Image, err = util.SavePicture(req.Image, "ware/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//  if len(req.DetailPic) > 0 {
//    req.DetailPic, err = util.SavePicture(req.DetailPic, "wareIntro/")
//    if err != nil {
//      logger.Error(err)
//      this.WriteStatusAndDataJSON(constants.ErrInternalServerError, nil)
//    }
//  }
//
//  err = ware.Service.UpdateWare(conn, req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//    if (len(req.Avatar) > 0 && !util.DeletePicture(req.Avatar)) ||
//       (len(req.Image) > 0 && !util.DeletePicture(req.Image)) ||
//       (len(req.DetailPic) > 0 && !util.DeletePicture(req.DetailPic)) {
//      logger.Error(errors.New("update ware failed and delete it's pictures go wrong, please delete picture manually"))
//    }
//  }
//
//  logger.Info("update ware info of id:", req.ID)
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
//}
//
//// modify price
//func (this *WareController) ModifyPrice() {
//  var (
//    err error
//    req ware.ModifyPriceReq
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = ware.Service.ModifyPrice(conn, req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  logger.Info("modify price of", req.ID)
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
//}
//
//// get homepage list with last wareID
//func (this *WareController) HomePageList() {
//  var (
//    err   error
//    idReq homeReq
//    res   []ware.BriefInfo
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &idReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(idReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  res, err = ware.Service.HomePageList(conn, idReq.LastID)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// get recommend list
//func (this *WareController) RecommendList() {
//  var (
//    err   error
//    idReq recommendReq
//    res   []ware.BriefInfo
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &idReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(idReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  res, err = ware.Service.GetRecommendList(conn, idReq.UserID)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// get detail info of ware
//func (this *WareController) GetDetail() {
//  var (
//    err error
//    req detailReq
//    res *ware.Ware
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(req)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  res, err = ware.Service.GetByID(conn, req.ID)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, res)
//}
//
//// change status of wares
//func (this *WareController) ChangeStatus() {
//  var (
//    err       error
//    changeReq ChangeStatusReq
//  )
//
//  conn, err := mysql.Pool.Get()
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  err = json.Unmarshal(this.Ctx.Input.RequestBody, &changeReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  err = this.Validate(changeReq)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrInvalidParam, nil)
//  }
//
//  // status: -1 -> delete or sold out
//  // 1 -> common
//  // 2 -> promotion
//  // 3 -> today new wares
//  err = ware.Service.ChangeStatus(conn, changeReq.IDs, changeReq.Status)
//  if err != nil {
//    logger.Error(err)
//    this.WriteStatusAndDataJSON(constants.ErrMysql, nil)
//  }
//
//  logger.Info("change ware status success")
//  this.WriteStatusAndDataJSON(constants.ErrSucceed, nil)
//}
