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
 *     Initial: 2017/11/16        Wang RiYu
 */

package ware

import (
  "time"
  "github.com/fengyfei/gu/libs/orm"
  "github.com/jinzhu/gorm"
)

type serviceProvider struct{}

var (
  Service *serviceProvider
)

type (
  Ware struct {
    ID         uint    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
    Name       string  `gorm:"type:varchar(50);not null"  json:"name" validate:"required"`
    Desc       string  `gorm:"type:varchar(100);not null" json:"desc"`
    Type       string  `gorm:"type:varchar(50);not null"  json:"type"`
    CategoryID uint    `gorm:"not null" json:"categoryId"`
    TotalSale  uint    `gorm:"not null" json:"total"`
    Inventory  uint    `gorm:"not null" json:"inventory"`
    Status     int8    `gorm:"type:TINYINT;default:1" json:"status"` // -1, hide or delete;1, common wares;2, promotion
    Price      float32 `gorm:"not null;type:float" json:"price"`
    SalePrice  float32 `gorm:"not null;type:float" json:"salePrice"` // promotion price
    Size       string  `gorm:"type:varchar(20)"    json:"size"`
    Color      string  `gorm:"type:varchar(20)"    json:"color"`
    Avatar     string  `gorm:"type:varchar(1000)"  json:"avatar"`
    Image      string  `gorm:"type:varchar(1000)"  json:"image"`
    Created    time.Time
  }

  BriefInfo struct {
    ID        uint    `json:"id"`
    Name      string  `json:"name"`
    TotalSale uint    `json:"total"`
    Inventory uint    `json:"inventory"`
    Status    int8    `json:"status"`
    Price     float32 `json:"price"`
    SalePrice float32 `json:"salePrice"`
    Avatar    string  `json:"avatar"`
  }

  UpdateReq struct {
    ID         uint   `json:"id" validate:"required"`
    Desc       string `json:"desc"`
    CategoryID uint   `json:"categoryId"`
    TotalSale  uint   `json:"total"`
    Size       string `json:"size"`
    Color      string `json:"color"`
    Avatar     string `json:"avatar"`
    Image      string `json:"image"`
    Inventory  uint   `json:"inventory"`
  }

  ModifyPriceReq struct {
    ID        uint    `json:"id" validate:"required"`
    Price     float32 `json:"price"`
    SalePrice float32 `json:"salePrice"`
  }
)

// add ware
func (sp *serviceProvider) CreateWare(conn orm.Connection, wareReq Ware) error {
  ware := &Ware{}
  ware.Name = wareReq.Name
  ware.Desc = wareReq.Desc
  ware.Type = wareReq.Type
  ware.CategoryID = wareReq.CategoryID
  ware.Price = wareReq.Price
  ware.SalePrice = wareReq.SalePrice
  ware.Avatar = wareReq.Avatar
  ware.Image = wareReq.Image
  ware.Inventory = wareReq.Inventory
  ware.Created = time.Now()

  db := conn.(*gorm.DB).Exec("USE shop")
  err := db.Model(&Ware{}).Create(ware).Error

  return err
}

// cid in arg is not nil ? get wares of category : get all wares from database
func (sp *serviceProvider) GetWareList(conn orm.Connection, arg ...uint) ([]Ware, error) {
  var (
    res  *gorm.DB
    list []Ware
  )

  db := conn.(*gorm.DB).Exec("USE shop")
  if len(arg) == 0 {
    res = db.Table("wares").Scan(&list)
  } else {
    res = db.Table("wares").Where("status > ? AND category_id = ?", 0, arg[0]).Scan(&list)
  }

  return list, res.Error
}

// get promotion wares (status = 2)
func (sp *serviceProvider) GetPromotionList(conn orm.Connection) ([]Ware, error) {
  var list []Ware

  db := conn.(*gorm.DB).Exec("USE shop")
  res := db.Table("wares").Where("status = ?", 2).Scan(&list)

  return list, res.Error
}

// update ware info
func (sp *serviceProvider) UpdateWare(conn orm.Connection, req UpdateReq) error {
  var res *gorm.DB

  db := conn.(*gorm.DB).Exec("USE shop")
  res = db.Table("wares").Where("id = ?", req.ID).Updates(req)

  return res.Error
}

// modify ware price
func (sp *serviceProvider) ModifyPrice(conn orm.Connection, req ModifyPriceReq) error {
  var res *gorm.DB

  db := conn.(*gorm.DB).Exec("USE shop")
  res = db.Table("wares").Where("id = ?", req.ID).Updates(map[string]interface{}{
    "price":      req.Price,
    "sale_price": req.SalePrice,
  })

  return res.Error
}

// get ware by id
func (sp *serviceProvider) GetByID(conn orm.Connection, id int32) (*Ware, error){
  db := conn.(*gorm.DB).Exec("USE shop")
  ware := &Ware{}
  err := db.Where("id = ?", id).First(&ware).Error

  return ware, err
}

// get homepage ware list, 10 wares per time
func (sp *serviceProvider) HomePageList(conn orm.Connection, id int) ([]BriefInfo, error) {
  var (
    count  uint
    list   []BriefInfo
    res    *gorm.DB
    fields = []string{"id", "name", "status", "price", "sale_price", "total_sale", "inventory", "avatar"}
  )

  db := conn.(*gorm.DB).Exec("USE shop")
  db.Table("wares").Where("status > ?", 0).Count(&count)

  if id == 0 || count <= 10 {
    if count <= 10 {
      res = db.Table("wares").Select(fields).Where("status > ?", 0).Scan(&list)
    } else {
      res = db.Table("wares").Select(fields).Order("id desc").Where("status > ?", 0).Limit(10).Scan(&list)
    }
  } else {
    res = db.Table("wares").Select(fields).Order("id desc").Where("status > ? AND id < ?", 0, id).Limit(10).Scan(&list)
  }

  return list, res.Error
}
