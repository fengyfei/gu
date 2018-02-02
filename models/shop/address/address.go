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
 *     Initial: 2017/11/18        ShiChao
 */

package address

import (
	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
	"fmt"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type Address struct {
	ID        uint   `gorm:"primary_key;auto_increment"`
	UserId    uint   `gorm:"not null"`
	Address   string `gorm:"type:varchar(128)"`
	IsDefault bool   `gorm:"not null";default:"false"`
}

func (this *serviceProvider) Add(conn orm.Connection, userId uint, address string, isDefault bool) error {
	var (
		err error
	)
	addr := &Address{}
	addr.Address = address
	addr.UserId = userId
	addr.IsDefault = isDefault

	db := conn.(*gorm.DB).Exec("USE shop")

	if !isDefault {
		return db.Model(&Address{}).Create(addr).Error
	}

	another := &Address{}
	err = db.Find(&another, "user_id = ? AND is_default = ?", userId, true).Error
	if err == gorm.ErrRecordNotFound {
		return db.Model(&Address{}).Create(addr).Error
	}

	another.IsDefault = false
	err = db.Save(&another).Error
	if err != nil {
		return err
	}
	return db.Model(&Address{}).Create(addr).Error
}

func (this *serviceProvider) SetDefault(conn orm.Connection, userId uint, id int) error {
	var (
		err     error
		addr    Address
		another Address
	)
	db := conn.(*gorm.DB).Exec("USE shop")

	err = db.Find(&addr, "id = ?", id).Error
	if err != nil {
		return err
	}
	if addr.IsDefault {
		return nil
	}
	addr.IsDefault = true
	err = db.Save(&addr).Error
	if err != nil {
		return err
	}

	err = db.Find(&another, "user_id = ? AND id <> ? AND is_default = true", userId, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	fmt.Println(err)
	another.IsDefault = false
	return db.Save(&another).Error
}

func (this *serviceProvider) Modify(conn orm.Connection, id int, address string) error {
	var (
		err  error
		addr Address
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err = db.Find(&addr, "id = ?", id).Error
	if err != nil {
		return err
	}

	addr.Address = address
	return db.Save(&addr).Error
}
