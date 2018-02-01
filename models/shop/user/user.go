/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd.
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
 *     Initial: 2018/02/01        Shi Ruitao
 */

package user

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/libs/security"
)

type serviceProvider struct{}

var (
	Service        *serviceProvider
	typeWechat     = "wechat"
	typePhone      = "phone"
	errLoginFailed = errors.New("invalid username or password.")
	errPassword    = errors.New("invalid password.")
)

type User struct {
	ID        uint   `sql:"primary_key;auto_increment"`
	UserName  string `gorm:"column:username"`
	NickName  string `gorm:"column:nickname"`
	Phone     string
	Type      string
	Pass      string
	CreatedAt *time.Time
}

func (this *serviceProvider) WechatLogin(conn orm.Connection, nickName, unionId *string) (uint, error) {

	user := &User{}
	res := &User{}
	user.UserName = *unionId
	user.NickName = *nickName
	user.Type = typeWechat

	db := conn.(*gorm.DB).Exec("USE shop")

	err := db.Where("username = ?", *unionId).First(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// not found, create new user
			err = db.Model(&User{}).Create(&user).Error
			if err != nil {
				return 0, err
			}

			return user.ID, nil
		}
		return 0, err
	}

	return res.ID, nil
}

// register by phoneNumber
func (this *serviceProvider) PhoneRegister(conn orm.Connection, phone, password, nickName *string) error {
	salt, err := security.SaltHashGenerate(password)
	if err != nil {
		return err
	}

	now := time.Now()

	user := &User{}
	user.UserName = *phone
	user.Phone = *phone
	user.Type = typePhone
	user.NickName = *nickName
	user.Pass = string(salt)
	user.CreatedAt = &now

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Create(&user).Error
}

func (this *serviceProvider) PhoneLogin(conn orm.Connection, phone, password *string) (uint, error) {

	db := conn.(*gorm.DB).Exec("USE shop")
	user := &User{}

	err := db.Where("phone = ?", *phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return 0, err
	}

	if !security.SaltHashCompare([]byte(user.Pass), password) {
		return 0, errLoginFailed
	}

	return user.ID, err
}

func (this *serviceProvider) ChangePassword(conn orm.Connection, id uint, oldPass, newPass *string) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	user := &User{}

	err := db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	if !security.SaltHashCompare([]byte(user.Pass), oldPass) {
		return errPassword
	}

	salt, err := security.SaltHashGenerate(newPass)
	if err != nil {
		return err
	}
	user.Pass = string(salt)
	return db.Save(&user).Error
}

func (this *serviceProvider) GetUserByID(conn orm.Connection, ID uint) (*User, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	user := &User{}

	err := db.Where("id = ?", ID).First(&user).Error

	return user, err
}
