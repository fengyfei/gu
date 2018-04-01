/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/03/28      Shi Ruitao
 */

package collection

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

// Collection represents the wares which user collects.
type Collection struct {
	ID      uint32     `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserId  uint32     `gorm:"column:userid" json:"user_id"`
	WareId  uint32     `gorm:"column:wareid" json:"ware_id"`
	Created *time.Time `gorm:"column:created"`
}

// Add adds wares into someone's collection.
func (this *serviceProvider) Add(conn orm.Connection, userId uint32, wareId uint32) error {
	now := time.Now()

	item := &Collection{
		UserId:  userId,
		WareId:  wareId,
		Created: &now,
	}

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Table("collection").Create(&item).Error
}

// GetByUserID gets someone's collection by userid.
func (this *serviceProvider) GetByUserID(conn orm.Connection, userId uint32) ([]Collection, error) {
	items := []Collection{}

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Table("collection").Where("userid = ?", userId).Find(&items).Error

	return items, err
}

// Remove deletes wares from collection.
func (this *serviceProvider) Remove(conn orm.Connection, id []uint32) error {
	var (
		err error
	)

	tx := conn.(*gorm.DB).Begin().Exec("USE shop")
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Table("collection").Where("id in (?)", id).Delete(&Collection{}).Error
	if err != nil {
		return err
	}

	return nil
}
