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

var (
	Service *serviceProvider

	ErrPhoneNotEmpty = errors.New("phone number is not empty")
	ErrLoginFailed   = errors.New("invalid username or password.")
	ErrPassword      = errors.New("invalid password.")
)

type serviceProvider struct{}

type (
	User struct {
		ID       uint      `sql:"primary_key;auto_increment" gorm:"column:id"`
		OpenID   string    `gorm:"column:openid"`
		UnionID  string    `gorm:"column:unionid"`
		UserName string    `gorm:"column:username"`
		Phone    string    `gorm:"column:phone"`
		Password string    `gorm:"column:password"`
		Avatar   string    `gorm:"column:avatar"`
		Sex      uint8     `gorm:"column:sex"`
		IsAdmin  bool      `gorm:"column:isadmin"`
		Created  time.Time `gorm:"column:created"`
	}

	UserData struct {
		Token    string `json:"token"`
		UserName string `json:"username"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Sex      uint8  `json:"sex"`
	}
)

type (
	WechatCode struct {
		Code string `json:"code" validate:"required"`
	}

	WechatData struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
	}

	WechatLoginErr struct {
		Errcode string `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	WechatLogin struct {
		OpenID  string
		UnionID string
	}

	AddPhone struct {
		Phone string `json:"phone" validate:"required,len=11"`
	}

	ChangeInfo struct {
		Sex    uint8  `json:"sex"`
		Avatar string `json:"avatar"`
	}
)

type (
	PhoneRegister struct {
		UserName string `json:"username" validate:"required,alphaunicode,min=2,max=30"`
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
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
	return "user"
}

// Login by wechat
func (this *serviceProvider) WechatLogin(conn orm.Connection, login *WechatLogin) (*User, error) {
	var user User

	db := conn.(*gorm.DB)

	err := db.Where("openid = ? AND unionid = ?", login.OpenID, login.UnionID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = User{
				OpenID:  login.OpenID,
				UnionID: login.UnionID,
				IsAdmin: false,
				Created: time.Now(),
			}
			err = db.Create(&user).Error
			if err != nil {
				return nil, err
			}

			return &user, nil
		}

		return nil, err
	}

	return &user, nil
}

// Add a phone number
func (this *serviceProvider) AddPhone(conn orm.Connection, id uint, phone *AddPhone) error {
	var (
		user User
		err  error
	)

	tx := conn.(*gorm.DB).Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	if user.Phone != "" {
		err = ErrPhoneNotEmpty
		return err
	}

	user.Phone = phone.Phone

	err = tx.Save(&user).Error
	return err
}

// Change information
func (this *serviceProvider) ChangeInfo(conn orm.Connection, id uint, change *ChangeInfo) error {
	var user User

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	user.Avatar = change.Avatar
	user.Sex = change.Sex

	return db.Save(&user).Error
}

// Register by phone
func (this *serviceProvider) PhoneRegister(conn orm.Connection, register *PhoneRegister) error {
	salt, err := security.SaltHashGenerate(&register.Password)
	if err != nil {
		return err
	}

	user := User{
		OpenID:   "1234567890123456789012345678",
		UnionID:  "12345678901234567890123456789",
		UserName: register.UserName,
		Phone:    register.Phone,
		Password: string(salt),
		IsAdmin:  false,
		Created:  time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&user).Error
}

// Login by phone
func (this *serviceProvider) PhoneLogin(conn orm.Connection, login *PhoneLogin) (*User, error) {
	var user User

	db := conn.(*gorm.DB)

	err := db.Where("phone = ?", login.Phone).First(&user).Error
	if err != nil {
		return nil, err
	}

	if !security.SaltHashCompare([]byte(user.Password), &login.Password) {
		return nil, ErrLoginFailed
	}

	return &user, nil
}

// Change password
func (this *serviceProvider) ChangePassword(conn orm.Connection, id uint, change *ChangePass) error {
	var (
		user User
		err  error
	)

	tx := conn.(*gorm.DB).Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	if !security.SaltHashCompare([]byte(user.Password), &change.OldPass) {
		err = ErrPassword
		return err
	}

	salt, err := security.SaltHashGenerate(&change.NewPass)
	if err != nil {
		return err
	}

	user.Password = string(salt)

	err = tx.Save(&user).Error
	return err
}

// Get the user by ID
func (this *serviceProvider) GetUserByID(conn orm.Connection, id uint) (*User, error) {
	db := conn.(*gorm.DB).Exec("USE shop")
	user := &User{}

	err := db.Where("id = ?", id).First(&user).Error

	return user, err
}
