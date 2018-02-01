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
 *     Initial: 2017/12/01      Lin Hao
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

type CollectionItem struct {
	ID        int32 `gorm:"primary_key;auto_increment"`
	UserId    int32
	WareId    int32
	CreatedAt *time.Time
}

func (this *serviceProvider) Add(conn orm.Connection, userId, wareId int32) error {
	now := time.Now()

	item := &CollectionItem{}
	item.CreatedAt = &now
	item.UserId = userId
	item.WareId = wareId

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Model(&CollectionItem{}).Create(&item).Error
}

func (this *serviceProvider) GetByUserID(conn orm.Connection, userId int32) ([]CollectionItem, error) {
	db := conn.(*gorm.DB).Exec("USE shop")

	items := []CollectionItem{}

	err := db.Where("user_id = ?", userId).Find(&items).Error

	return items, err
}

func (this *serviceProvider) Remove(conn orm.Connection, id int32) error {
	db := conn.(*gorm.DB).Exec("USE shop")

	item := &CollectionItem{}
	item.ID = id

	return db.Delete(&item).Error
}
