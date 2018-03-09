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
 *     Initial: 2017/11/30        Tong Yuehong
 */

package cart

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider.
	Service *serviceProvider
)

type (
	Cart struct {
		ID      uint32    `gorm:"primary_key;auto_increment"  json:"id"`
		UserId  uint32    `gorm:"not null"                    json:"user_id"`
		WareId  uint32    `gorm:"not null"                    json:"ware_id"`
		Count   uint8     `gorm:"not null;default:0"          json:"count"`
		Created time.Time `gorm:"column:created"              json:"created"`
	}

	AddCartReq struct {
		WareId uint32 `json:"wareId" validate:"required"`
		Count  uint8 `json:"count"  validate:"required"`
	}

	RemoveCartReq struct {
		Id uint32 `json:"id" validate:"required"`
	}
)

// TableName returns table name in database.
func (Cart) TableName() string {
	return "cart"
}

// Add adds wares to cart.
func (this *serviceProvider) Add(conn orm.Connection, userId uint32, wareId uint32, count uint8) error {
	var (
		cart Cart
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Where("ware_id = ? AND user_id = ?", wareId, userId).Find(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cart := &Cart{
				UserId:  userId,
				WareId:  wareId,
				Count:   count,
				Created: time.Now(),
			}
			return db.Table("cart").Create(&cart).Error
		} else {
			return err
		}
	}

	cart.Count += count
	return db.Save(&cart).Error
}

// GetByUserID gets carts by userid.
func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint32) ([]Cart, error) {
	var (
		carts []Cart
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Where("user_id = ?", userId).Find(&carts).Error

	return carts, err
}

// RemoveById deletes a ware by id.
func (this *serviceProvider) RemoveById(conn orm.Connection, id uint32) error {
	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Table("cart").Where("id = ?", id).Delete(&Cart{}).Error
}

// RemoveWhenOrder deletes wares by ids.
func (this *serviceProvider) RemoveWhenOrder(conn orm.Connection, wareIdList []uint32) error {
	var (
		err error
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	for i := 0; i < len(wareIdList); i++ {
		err = db.Table("cart").Where("id = ?", wareIdList[i]).Delete(&Cart{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
