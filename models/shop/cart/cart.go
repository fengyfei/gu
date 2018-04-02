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
 *     Initial: 2018/03/29        Shi Ruitao
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
		ID      uint32    `gorm:"column:id;primary_key;auto_increment" json:"id"`
		UserId  uint32    `gorm:"column:userid;not null"               json:"user_id"`
		WareId  uint32    `gorm:"column:wareid;not null"               json:"ware_id"`
		Count   uint8     `gorm:"column:count;not null;default:0"      json:"count"`
		Created time.Time `gorm:"column:created"                       json:"created"`
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

	err := db.Where("wareid = ? AND userid = ?", wareId, userId).Find(&cart).Error
	if err == nil {
		cart.Count += count
		return db.Save(&cart).Error
	}
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

// GetByUserID gets carts by userid.
func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint32) ([]Cart, error) {
	var (
		carts []Cart
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Where("userid = ?", userId).Find(&carts).Error

	return carts, err
}

// RemoveWhenOrder deletes wares by ids.
func (this *serviceProvider) RemoveWhenOrder(conn orm.Connection, wareIds []uint32) (err error) {
	tx := conn.(*gorm.DB).Begin().Exec("USE shop")
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Table("cart").Where("id in (?)", wareIds).Delete(&Cart{}).Error
	if err != nil {
		return err
	}

	return nil
}
