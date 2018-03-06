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
 *     Initial: 2018/02/02        Shi Ruitao
 *     Modify : 2018/03/06        Tong Yuehong
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

var (
	Service *serviceProvider
)

type serviceProvider struct{}

type (
	Address struct {
		ID        uint64    `gorm:"column:id;primary_key;auto_increment"`
		UserID    uint32    `gorm:"column:userid;not null"`
		Name      string    `gorm:"column:name;not null"`
		Phone     string    `gorm:"column:phone;not null"`
		Address   string    `gorm:"column:address;type:varchar(128)"`
		IsDefault bool      `gorm:"column:isdefault;not null;default:false"`
		Created   time.Time `gorm:"column:created"`
	}

	AddressData struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		IsDefault bool   `json:"isdefault"`
	}

	Modify struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required,len=11"`
		Address string `json:"address" validate:"required,max=128"`
	}
)

func (Address) TableName() string {
	return "address"
}

// Add the address
func (this *serviceProvider) Add(conn orm.Connection, userID uint32, add *AddressData) error {
	var (
		address Address
		err     error
	)

	addr := Address{
		UserID:    userID,
		Name:      add.Name,
		Phone:     add.Phone,
		Address:   add.Address,
		IsDefault: add.IsDefault,
		Created:   time.Now(),
	}

	db := conn.(*gorm.DB)

	if add.IsDefault {
		err = db.Where("userid = ? AND isdefault = ?", userID, DefaultAddress).Find(&address).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = db.Create(&addr).Error
				return err
			}
			return err
		}
		address.IsDefault = NotDefaultAddress
		err = db.Save(&address).Error
		if err != nil {
			return err
		}
	}

	err = db.Create(&addr).Error
	return err
}

// Set default address
func (this *serviceProvider) SetDefault(conn orm.Connection, userID uint32, id uint64) error {
	var (
		err     error
		add     Address
		address Address
	)

	tx := conn.(*gorm.DB).Begin().Exec("USE shop")
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("userid = ? AND isdefault = ?", userID, true).First(&add).Error
	if err != nil {
		return err
	}

	add.IsDefault = false
	err = tx.Save(&add).Error
	if err != nil {
		return err
	}

	err = tx.Table("address").Where("id = ?", id).First(&address).Error
	if err != nil {
		return err
	}

	address.IsDefault = true
	return tx.Save(&address).Error
}

// Modify address
func (this *serviceProvider) Modify(conn orm.Connection, modify *Modify) error {
	db := conn.(*gorm.DB)

	return db.Table("address").Where("id = ?", modify.ID).Update(map[string]interface{}{
		"name": modify.Name,
		"phone": modify.Phone,
		"address": modify.Address,
	}).Error
}

// Read the address
func (this *serviceProvider) Get(conn orm.Connection, userID uint32) (*[]Address, error) {
	var (
		address []Address
	)

	db := conn.(*gorm.DB)

	err := db.Model(&Address{}).Where("userid = ?", userID).Find(&address).Error

	return &address, err
}

func (this *serviceProvider) Delete(conn orm.Connection, id uint64) error {
	db := conn.(*gorm.DB)

	//address := Address{
	//	ID: id,
	//}
	//
	//return db.Where("userid = ?", userID).Delete(&address).Error

	return db.Table("address").Where("id = ?", id).Delete(&Address{}).Error
}
