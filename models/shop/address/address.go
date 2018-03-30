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
 *     Initial: 2018/03/27        Shi Ruitao
 */

package address

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

const (
	DefaultAddress    = true
	NotDefaultAddress = false
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider.
	Service *serviceProvider
)

type (
	// Address represents the address of users.
	Address struct {
		ID        uint32    `gorm:"column:id;primary_key;auto_increment" json:"id"`
		UserID    uint32    `gorm:"column:userid;not null"`
		UserName  string    `gorm:"column:name;not null" json:"user_name"`
		Phone     string    `gorm:"column:phone;not null" json:"phone"`
		Address   string    `gorm:"column:address;type:varchar(128)" json:"address"`
		IsDefault bool      `gorm:"column:isdefault;not null;default:false" json:"is_default"`
		Created   time.Time `gorm:"column:created"`
	}
)

// TableName returns table name in database.
func (Address) TableName() string {
	return "address"
}

// Add adds address.
func (this *serviceProvider) Add(conn orm.Connection, userID uint32, add *Address) (err error) {
	var (
		address Address
	)

	addr := &Address{
		UserID:    userID,
		UserName:  add.UserName,
		Phone:     add.Phone,
		Address:   add.Address,
		IsDefault: add.IsDefault,
		Created:   time.Now(),
	}

	tx := conn.(*gorm.DB).Begin().Exec("USE shop")
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	if add.IsDefault {
		err = tx.Where("userid = ? AND isdefault = ?", userID, DefaultAddress).Find(&address).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = tx.Create(&addr).Error
				return err
			}
			return err
		}
		address.IsDefault = NotDefaultAddress
		err = tx.Save(&address).Error
		if err != nil {
			return err
		}
	}

	err = tx.Create(&addr).Error
	return err
}

// SetDefault sets default address.
func (this *serviceProvider) SetDefault(conn orm.Connection, userID uint32, id uint32) (err error) {
	tx := conn.(*gorm.DB).Begin().Exec("USE shop")
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Model(&Address{}).Where("userid = ? AND isdefault = ?", userID, true).Update("isdefault", NotDefaultAddress).Limit(1).Error
	if err != nil {
		return err
	}

	return tx.Model(&Address{}).Where("id = ? AND userid = ?", id, userID).Update("isdefault", DefaultAddress).Limit(1).Error
}

// Modify modify the address.
func (this *serviceProvider) Modify(conn orm.Connection, userID uint32, add *Address) error {
	db := conn.(*gorm.DB)
	return db.Model(&Address{}).Where("id = ? && userid = ?", add.ID, userID).Updates(map[string]interface{}{"name": add.UserName, "phone": add.Phone, "address": add.Address, "isdefault": add.IsDefault}).Limit(1).Error
}

// Get gets address by userid.
func (this *serviceProvider) Get(conn orm.Connection, userID uint32) (*[]Address, error) {
	var (
		address []Address
	)

	db := conn.(*gorm.DB)

	err := db.Model(&Address{}).Where("userid = ?", userID).Find(&address).Error

	return &address, err
}

// Delete deletes address by id.
func (this *serviceProvider) Delete(conn orm.Connection, id uint32) error {
	db := conn.(*gorm.DB)

	return db.Table("address").Where("id = ?", id).Delete(&Address{}).Limit(1).Error
}
