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
 *	   Modify: 2018/02/02         Shi Ruitao
 */

package cart

import (
	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
	"time"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type (
	CartItem struct {
		ID        uint `gorm:"primary_key;auto_increment"`
		UserId    uint `gorm:"not null"`
		WareId    uint `gorm:"not null"`
		Count     uint `gorm:"not null";default:0`
		CreatedAt *time.Time
	}
	AddCartReq struct {
		WareId uint `json:"wareId"`
		Count  uint `json:"count"`
	}

	RemoveCartReq struct {
		Id uint `json:"id"`
	}
)

func (this *serviceProvider) Add(conn orm.Connection, userId, wareId, count uint) error {
	db := conn.(*gorm.DB).Exec("USE shop")

	item := &CartItem{}
	item.UserId = userId
	item.WareId = wareId
	err := db.First(&item).Error
	if err == nil {
		item.Count += count
		return db.Save(&item).Error
	}

	if err == gorm.ErrRecordNotFound {
		now := time.Now()
		item.CreatedAt = &now
		item.Count = count
	} else {
		return err
	}

	return db.Model(&CartItem{}).Create(&item).Error
}

func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint) ([]CartItem, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	items := []CartItem{}

	err := db.Where("user_id = ?", userId).Find(&items).Error

	return items, err
}

func (this *serviceProvider) RemoveById(conn orm.Connection, id uint) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	item := &CartItem{}
	item.ID = id

	return db.Delete(&item).Error
}

func (this *serviceProvider) RemoveWhenOrder(tx *gorm.DB, userId uint, wareIdList []uint) error {
	for i := 0; i < len(wareIdList); i++ {
		item := &CartItem{}
		item.ID = wareIdList[i]
		err := tx.Delete(&item).Error
		if err != nil {
			return err
		}
	}
	return nil
}
