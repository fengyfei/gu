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
		UserID    uint64    `gorm:"column:userid;not null"`
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
)

type (
	Add struct {
		Name      string `json:"name" validate:"required"`
		Phone     string `json:"phone" validate:"required,len=11"`
		Address   string `json:"address" validate:"required,max=128"`
		IsDefault bool   `json:"isdefault"`
	}

	SetDefault struct {
		ID uint64 `json:"id"`
	}

	Modify struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required,len=11"`
		Address string `json:"address" validate:"required,max=128"`
	}

	Delete struct {
		ID uint64 `json:"id"`
	}
)

func (Address) TableName() string {
	return "address"
}

// Add the address
func (this *serviceProvider) Add(conn orm.Connection, userID uint64, add *Add) error {
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
		err = db.Where("userid = ? AND isdefault = ?", userID, DefaultAddress).First(&address).Error
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
func (this *serviceProvider) SetDefault(conn orm.Connection, userID, id uint64) error {
	var (
		address Address
		addr    Address
		err     error
	)

	tx := conn.(*gorm.DB).Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("userid = ? AND isdefault = ?", userID, DefaultAddress).First(&address).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	} else if err != nil {
		return err
	} else {
		if address.ID == id && address.UserID == userID {
			return nil
		}

		address.IsDefault = NotDefaultAddress
		err = tx.Save(&address).Error
		if err != nil {
			return err
		}
	}

	err = tx.Where("id = ? AND userid = ?", id, userID).Find(&addr).Error
	if err != nil {
		return err
	}

	addr.IsDefault = DefaultAddress

	err = tx.Save(&addr).Error
	return err
}

// Modify address
func (this *serviceProvider) Modify(conn orm.Connection, userID uint64, modify *Modify) error {
	var address Address

	db := conn.(*gorm.DB)

	err := db.Where("id = ? AND userid = ?", modify.ID, userID).First(&address).Error
	if err != nil {
		return err
	}

	address.Name = modify.Name
	address.Phone = modify.Phone
	address.Address = modify.Address
	return db.Save(&address).Error
}

// Read the address
func (this *serviceProvider) Get(conn orm.Connection, userID uint) (*[]Address, error) {
	var (
		address []Address
	)

	db := conn.(*gorm.DB)

	err := db.Model(&Address{}).Where("userid = ?", userID).Find(&address).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return &address, err
}

func (this *serviceProvider) Delete(conn orm.Connection, userID, id uint64) error {
	db := conn.(*gorm.DB)

	address := Address{
		ID: id,
	}

	return db.Where("userid = ?", userID).Delete(&address).Error
}
