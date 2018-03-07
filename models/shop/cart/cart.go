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
 *     Initial: 2017/11/30        ShiChao
 *	   Modify : 2018/02/02        Shi Ruitao
 *	   Modify : 2018/03/06        Tong Yuehong
 */

package cart

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type (
	Cart struct {
		ID      uint64    `gorm:"primary_key;auto_increment"  json:"id"`
		UserId  uint32    `gorm:"not null"                    json:"user_id"`
		WareId  uint64    `gorm:"not null"                    json:"ware_id"`
		Count   uint32    `gorm:"not null;default:0"          json:"count"`
		Created time.Time `gorm:"column:created"              json:"created"`
	}

	AddCartReq struct {
		WareId uint64 `json:"wareId"`
		Count  uint32 `json:"count"`
	}

	RemoveCartReq struct {
		Id uint64 `json:"id" validate:"required"`
	}
)

func (Cart) TableName() string {
	return "cart"
}

func (this *serviceProvider) Add(conn orm.Connection, userId uint32, wareId uint64, count uint32) error {
	var (
		cart Cart
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Where("ware_id = ?", wareId).Find(&cart).Error
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

func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint32) ([]Cart, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	var carts []Cart

	err := db.Where("user_id = ?", userId).Find(&carts).Error

	return carts, err
}

func (this *serviceProvider) RemoveById(conn orm.Connection, id uint64) error {
	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Table("cart").Where("id = ?", id).Delete(&Cart{}).Error
}

func (this *serviceProvider) RemoveWhenOrder(conn orm.Connection, wareIdList []uint64) error {
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
