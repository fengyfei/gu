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
 *     Initial: 2018/02/01        Shi Ruitao
 *     Modify:  2018/02/01        Li Zebang
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

const (
	typeWechat = "wechat"
	typePhone  = "phone"
)

var (
	Service *serviceProvider

	errLoginFailed = errors.New("invalid username or password.")
	errPassword    = errors.New("invalid password.")
)

type (
	User struct {
		ID       uint      `sql:"primary_key;auto_increment"`
		UserName string    `gorm:"column:username"`
		NickName string    `gorm:"column:nickname"`
		Phone    string    `gorm:"column:phone"`
		Type     string    `gorm:"column:type"`
		Password string    `gorm:"column:password"`
		Created  time.Time `gorm:"column:created"`
	}

	WechatLoginReq struct {
		UserName   string `json:"userName" validate:"required,alphanum,min=6,max=30"`
		WechatCode string `json:"wechatCode" validate:"required"`
	}

	WechatLogin struct {
		Data WechatLoginData
	}

	WechatLoginData struct {
		Errmsg  string
		Unionid string
	}

	PhoneRegister struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
		NickName string `json:"name" validate:"required,alphaunicode,min=2,max=30"`
	}

	PhoneLogin struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
	}

	ChangePass struct {
		OldPass string `json:"oldPass" validate:"required,min=6,max=30"`
		NewPass string `json:"newPass" validate:"required,min=6,max=30"`
	}
)

func (User) TableName() string {
	return "users"
}

// Login by wechat
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

// Register by phoneNumber
func (this *serviceProvider) PhoneRegister(conn orm.Connection, req *PhoneRegister) error {
	salt, err := security.SaltHashGenerate(&req.Password)
	if err != nil {
		return err
	}

	user := User{
		UserName: req.Phone,
		Phone:    req.Phone,
		Type:     typePhone,
		NickName: req.NickName,
		Password: string(salt),
		Created:  time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&user).Error
}

// Login by phone
func (this *serviceProvider) PhoneLogin(conn orm.Connection, req *PhoneLogin) (uint, error) {
	var (
		user User
	)

	db := conn.(*gorm.DB)

	err := db.Where("phone = ?", req.Phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return 0, err
	}

	if !security.SaltHashCompare([]byte(user.Password), &req.Password) {
		return 0, errLoginFailed
	}

	return user.ID, err
}

// Change password
func (this *serviceProvider) ChangePassword(conn orm.Connection, id uint, req *ChangePass) error {
	var (
		user User
	)

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	if !security.SaltHashCompare([]byte(user.Password), &req.OldPass) {
		return errPassword
	}

	salt, err := security.SaltHashGenerate(&req.NewPass)
	if err != nil {
		return err
	}

	user.Password = string(salt)

	return db.Save(&user).Error
}

// Get the user by ID
func (this *serviceProvider) GetUserByID(conn orm.Connection, ID uint) (*User, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	user := &User{}

	err := db.Where("id = ?", ID).First(&user).Error

	return user, err
}
