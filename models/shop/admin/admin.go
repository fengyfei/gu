/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 *     Initial: 2017/11/26        ShiChao
 */

package admin

import (
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/fengyfei/gu/libs/orm"
	"time"
	"github.com/fengyfei/gu/libs/security"
)

type serviceProvider struct{}

var (
	Service        *serviceProvider
	errLoginFailed = errors.New("invalid username or password.")
	errPassword    = errors.New("invalid password.")
)

type Admin struct {
	ID        int32  `gorm:"primary_key;auto_increment"`
	Name  string `gorm:"unique;type:varchar(128)"`
	Pass      string `gorm:"type:varchar(128)"`
	CreatedAt *time.Time
}

func (this *serviceProvider) AddAdmin(conn orm.Connection, name, password *string) error {
	salt, err := security.SaltHashGenerate(password)
	if err != nil {
		return err
	}

	now := time.Now()

	admin := &Admin{}
	admin.Name = *name
	admin.Pass = string(salt)
	admin.CreatedAt = &now

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Model(&Admin{}).Create(&admin).Error
}

func (this *serviceProvider) Login(conn orm.Connection, name, password *string) (int32, error) {

	db := conn.(*gorm.DB).Exec("USE shop")
	admin := &Admin{}

	err := db.Where("name = ?", *name).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return 0, err
	}

	if !security.SaltHashCompare([]byte(admin.Pass), password) {
		return 0, errLoginFailed
	}

	return admin.ID, err
}

func (this *serviceProvider) ChangePassword(conn orm.Connection, id int32, oldPass, newPass *string) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	admin := &Admin{}

	err := db.Where("id = ?", id).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	if !security.SaltHashCompare([]byte(admin.Pass), oldPass) {
		return errPassword
	}

	salt, err := security.SaltHashGenerate(newPass)
	if err != nil {
		return err
	}
	admin.Pass = string(salt)
	return db.Save(&admin).Error
}

func (this *serviceProvider) GetAdminByID(conn orm.Connection, ID int32) (*Admin, error){
	db := conn.(*gorm.DB).Exec("USE shop")
	admin := &Admin{}

	err := db.Where("id = ?", ID).First(&admin).Error

	return admin, err
}
