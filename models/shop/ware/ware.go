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
 *     Initial: 2018/03/06        Shi Ruitao
 */

package ware

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

const (
	// hide or delete wares
	invalidation = -1
	//judge the effectiveness of the wares
	judgement = 0
	// common wares
	common = 1
	// promotion wares
	promotion = 2
	// new wares
	newWare = 3
	//recommend wares
	recommend = 4
)

type (
	// Ware represents the ware information
	Ware struct {
		ID               uint32    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
		Name             string    `gorm:"type:varchar(50);not null"  json:"name" validate:"required,alphanumunicode,max=12"`
		Desc             string    `gorm:"type:varchar(100);not null" json:"desc" validate:"alphanumunicode,max=50"`
		ParentCategoryID uint32    `gorm:"not null" json:"parent_category_id"`
		CategoryID       uint32    `gorm:"not null" json:"category_id"`
		TotalSale        uint32    `gorm:"not null" json:"totalale"`
		Inventory        uint32    `gorm:"not null" json:"inventory"`
		Status           int8      `gorm:"not null;type:TINYINT;default:1" json:"status"` // -1, hide or delete;1, common wares;2, promotion;3, new wares;4, recommend wares
		Price            float32   `gorm:"not null;type:float" json:"price"`
		SalePrice        float32   `gorm:"not null;type:float" json:"sale_price"` // promotion price
		Avatar           string    `gorm:"not null;type:varchar(100)"   json:"avatar"`
		Image            string    `gorm:"type:varchar(100)"   json:"image"`
		DetailPic        string    `gorm:"type:varchar(100)"   json:"detail_pic"`
		CreatedAt        time.Time `json:"createdAt"`
	}

	// BriefInfo represents the wares display
	BriefInfo struct {
		ID        uint32  `json:"id"`
		Name      string  `json:"name"`
		TotalSale uint32  `json:"total_sale"`
		Inventory uint32  `json:"inventory"`
		Status    int8    `json:"status"`
		Price     float32 `json:"price"`
		SalePrice float32 `json:"sale_price"`
		Avatar    string  `json:"avatar"`
	}

	//UpdateReq represents information to be updated.
	UpdateReq struct {
		ID               uint32 `json:"id" validate:"required"`
		Name             string `json:"name"`
		Desc             string `json:"desc" validate:"max=50"`
		ParentCategoryID uint32 `json:"parent_category_id"`
		CategoryID       uint32 `json:"category_id"`
		TotalSale        uint32 `json:"total_sale"`
		Avatar           string `json:"avatar"`
		Image            string `json:"image"`
		DetailPic        string `json:"detail_pic"`
		Inventory        uint32 `json:"inventory"`
	}
	// ModifyPriceReq represents promotion price and original price to be updated
	ModifyPriceReq struct {
		ID        uint32  `json:"id" validate:"required"`
		Price     float32 `json:"price"`
		SalePrice float32 `json:"sale_price"`
	}
)

// add ware
func (sp *serviceProvider) CreateWare(conn orm.Connection, wareReq *Ware) error {
	ware := &Ware{}
	ware.Name = wareReq.Name
	ware.Desc = wareReq.Desc
	ware.ParentCategoryID = wareReq.ParentCategoryID
	ware.CategoryID = wareReq.CategoryID
	ware.Price = wareReq.Price
	ware.SalePrice = wareReq.SalePrice
	ware.Avatar = wareReq.Avatar
	ware.Image = wareReq.Image
	ware.DetailPic = wareReq.DetailPic
	ware.Inventory = wareReq.Inventory

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Model(&Ware{}).Create(ware).Error

	return err
}

// get all wares
func (sp *serviceProvider) GetAllWare(conn orm.Connection) ([]Ware, error) {
	var (
		res  *gorm.DB
		list []Ware
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Scan(&list)

	return list, res.Error
}

// get wares by parent categoryID
func (sp *serviceProvider) GetByParentCID(conn orm.Connection, cid uint32) ([]BriefInfo, error) {
	var (
		res  *gorm.DB
		list []BriefInfo
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Where("status > ? AND parent_category_id = ?", 0, cid).Scan(&list)
	return list, res.Error
}

// get wares by categoryID
func (sp *serviceProvider) GetByCID(conn orm.Connection, cid uint32) ([]BriefInfo, error) {
	var (
		res  *gorm.DB
		list []BriefInfo
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Where("status > ? AND category_id = ?", 0, cid).Scan(&list)
	return list, res.Error
}

// get promotion wares (status = 2)
func (sp *serviceProvider) GetPromotionList(conn orm.Connection) ([]BriefInfo, error) {
	var list []BriefInfo

	db := conn.(*gorm.DB).Exec("USE shop")
	res := db.Table("wares").Where("status = ?", promotion).Scan(&list)

	return list, res.Error
}

// get new ware
func (sp *serviceProvider) GetNewWares(conn orm.Connection) ([]BriefInfo, error) {
	var (
		res  *gorm.DB
		list []BriefInfo
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Where("status = ?", newWare).Scan(&list)

	return list, res.Error
}

// update ware info
func (sp *serviceProvider) UpdateWare(conn orm.Connection, req *UpdateReq) error {
	var imgs Ware

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Table("wares").Where("id = ?", req.ID).Updates(req).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	if len(req.Avatar) > 0 {
		util.DeletePicture(imgs.Avatar)
	}
	if len(req.Image) > 0 {
		util.DeletePicture(imgs.Image)
	}
	if len(req.DetailPic) > 0 {
		util.DeletePicture(imgs.DetailPic)
	}
	return nil
}

// modify ware price
func (sp *serviceProvider) ModifyPrice(conn orm.Connection, req *ModifyPriceReq) error {
	var res *gorm.DB

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Where("id = ?", req.ID).Updates(map[string]interface{}{
		"price":      req.Price,
		"sale_price": req.SalePrice,
	})

	return res.Error
}

// get detail ware by id
func (sp *serviceProvider) GetByID(conn orm.Connection, id uint32) (*Ware, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	ware := &Ware{}
	err := db.Table("wares").Where("id = ?", id).First(ware).Error

	return ware, err
}

// get homepage ware list, 10 wares per time
func (sp *serviceProvider) HomePageList(conn orm.Connection, id uint32) ([]BriefInfo, error) {
	var (
		count  uint32
		list   []BriefInfo
		res    *gorm.DB
		fields = []string{"id", "name", "status", "price", "sale_price", "total_sale", "inventory", "avatar"}
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	db.Table("wares").Where("status > ?", judgement).Count(&count)

	if id == 0 || count <= 10 {
		if count <= 10 {
			res = db.Table("wares").Select(fields).Where("status > ?", judgement).Scan(&list)
		} else {
			res = db.Table("wares").Select(fields).Order("id desc").Where("status > ?", judgement).Limit(10).Scan(&list)
		}
	} else {
		res = db.Table("wares").Select(fields).Order("id desc").Where("status > ? AND id < ?", judgement, id).Limit(10).Scan(&list)
	}

	return list, res.Error
}

// RecommendList
func (sp *serviceProvider) GetRecommendList(conn orm.Connection) ([]BriefInfo, error) {
	var (
		list   []BriefInfo
		res    *gorm.DB
		fields = []string{"id", "name", "status", "price", "sale_price", "total_sale", "inventory", "avatar"}
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Select(fields).Order("total_sale desc").Where("status = ?", recommend).Limit(10).Scan(&list)

	return list, res.Error
}

// change status
func (sp *serviceProvider) ChangeStatus(conn orm.Connection, reqList []uint32, status int8) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Table("wares").Where("id in (?)", reqList).Update("status", status).Error

	return err
}

// get wares by id slice
func (sp *serviceProvider) GetByIDs(conn orm.Connection, ids []string) ([]BriefInfo, error) {
	var (
		res  *gorm.DB
		list []BriefInfo
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res = db.Table("wares").Where("status > ? AND id in (?) ", judgement, ids).Scan(&list)

	return list, res.Error
}
