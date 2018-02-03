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
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type (
	Address struct {
		ID        uint   `gorm:"column:id;primary_key;auto_increment"`
		UserId    uint   `gorm:"column:userid;not null"`
		Address   string `gorm:"column:address;type:varchar(128)"`
		IsDefault bool   `gorm:"column:isdefault;not null";default:"false"`
		Phone     string `gorm:"column:phone;not null"`
	}

	AddReq struct {
		Address   string `json:"address" validate:"max=128"`
		Phone     string `json:"phone" validate:"len=11"`
		IsDefault bool   `json:"isdefault"`
	}

	SetDefaultReq struct {
		Id int `json:"id" validate:"required"`
	}

	ModifyReq struct {
		Id        int    `json:"id" validate:"required"`
		Address   string `json:"address" validate:"max=128"`
		IsDefault bool   `json:"isdefault"`
		Phone     string `json:"phone" validate:"len=11"`
	}
)

func (Address) TableName() string {
	return "address"
}

// Add the address
func (this *serviceProvider) Add(conn orm.Connection, userId uint, req *AddReq) error {
	var (
		addr       Address
		repetition Address
		err        error
	)

	addr.Address = req.Address
	addr.UserId = userId
	addr.IsDefault = req.IsDefault
	addr.Phone = req.Phone

	db := conn.(*gorm.DB)

	err = db.Where("address = ? and phone = ?", req.Address, req.Phone).First(&repetition).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if !req.IsDefault {
				return db.Create(&addr).Error
			}

			another := Address{}
			err = db.Find(&another, "userid = ? AND isdefault = ?", userId, true).Error
			if err == gorm.ErrRecordNotFound {
				return db.Create(&addr).Error
			}

			another.IsDefault = false
			err = db.Save(&another).Error
			if err != nil {
				return err
			}
			return db.Create(&addr).Error
		}
	}

	err = errors.New("Repeat the address")
	return err
}

// Read the address
func (this *serviceProvider) AddressRead(conn orm.Connection, userId uint) (*[]Address, error) {
	var (
		err     error
		address []Address
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err = db.Find(&address, "userid = ?", userId).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return &address, err
}

// Set default address
func (this *serviceProvider) SetDefault(conn orm.Connection, userId uint, id int) error {
	var (
		err     error
		addr    Address
		another Address
	)

	db := conn.(*gorm.DB).Exec("USE shop")

	err = db.Find(&addr, "userid = ? AND isdefault = true", userId).Error
	if err == gorm.ErrRecordNotFound {
		logger.Info("IsDefault is false")
	} else {
		logger.Error("db.Find():", err)
		return err
	}

	addr.IsDefault = false

	err = db.Find(&another, "id = ?", id).Error
	if err != nil {
		return err
	}
	if addr.IsDefault {
		return nil
	}

	another.IsDefault = true
	return db.Save(&another).Error
}

// Modify address
func (this *serviceProvider) Modify(conn orm.Connection, req *ModifyReq) error {
	var (
		err        error
		addr       Address
		repetition Address
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err = db.Find(&addr, "id = ?", req.Id).Error
	if err != nil {
		return err
	}

	addr.Address = req.Address
	addr.Phone = req.Phone
	err = db.Where("address = ? and phone = ?", req.Address, req.Phone).First(&repetition).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.Save(&addr).Error
		}
	}
	err = errors.New("Repeat the address")
	return err
}
