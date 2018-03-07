/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the 'Software'), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/03/07      Tong Yuehong
 */

package collection

import (
	"time"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type Collection struct {
	ID      int32  `gorm:"primary_key;auto_increment" json:"id"`
	UserId  uint32 `json:"user_id"`
	WareId  uint64 `json:"ware_id"`
	Created *time.Time
}

func (this *serviceProvider) Add(conn orm.Connection, userId uint32, wareId uint64) error {
	now := time.Now()

	item := &Collection{
		UserId:  userId,
		WareId:  wareId,
		Created: &now,
	}

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Table("Collect").Create(&item).Error
}

func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint32) ([]Collection, error) {
	items := []Collection{}

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Table("collect").Where("user_id = ?", userId).Find(&items).Error

	return items, err
}

func (this *serviceProvider) Remove(conn orm.Connection, id uint64) error {
	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Table("collect").Where("id = ?", id).Delete(&Collection{}).Error
}
