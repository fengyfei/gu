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
 */

package user

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/libs/security"
)

// UserServiceProvider
type UserServiceProvider struct{}

const (
	// WeChat
	WeChat = iota
	// Mobile
	Mobile
)

var (
	// UserService
	UserService = &UserServiceProvider{}

	// ErrInvalidPass
	ErrInvalidPass = errors.New("the password error.")
	// Mobile phone registration cannot add phone number.
	ErrAddPhone = errors.New("Mobile phone registration cannot add phone number.")
)

type (
	PhoneRegister struct {
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
		UserName string `json:"username" validate:"required,alphaunicode,min=2,max=30"`
	}

	PhoneLogin struct {
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
	}

	ChangePass struct {
		OldPass string `json:"oldpass" validate:"required,min=6,max=30"`
		NewPass string `json:"newpass" validate:"required,min=6,max=30"`
	}

	WechatLogin struct {
		UnionID    string `json:"unionid"`
		SessionKey string `json:"session_key"`
	}

	WechatCode struct {
		UserName string `json:"username"`
		Sex      uint8  `json:"sex"`
		Phone    string `json:"phone" validate:"required,numeric,len=11"`
		Code     string `json:"code" validate:"required"`
	}

	WechatData struct {
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
	}

	WechatPhone struct {
		Phone string `json:"phone" validate:"required,numeric,len=11"`
	}

	UserData struct {
		Token    string `json:"token"`
		UserName string `json:"username"` // todo
		Phone    string `json:"phone"`    // todo
		Avatar   string `json:"avatar"`
		Sex      uint8  `json:"sex"`
	}

	ChangeInfo struct {
		UserName string `json:"username"`
		Sex      uint8  `json:"sex"`
		Avatar   string `json:"avatar"`
	}
)

// User represents users information
type User struct {
	UserID    uint32    `gorm:"column:id;primary_key;auto_increment" json:"user_id"`
	UserName  string    `gorm:"column:username;type:varchar(128)" json:"user_name"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Sex       uint8     `gorm:"column:sex"`
	Password  string    `gorm:"column:password;type:varchar(128)" json:"password""`
	Phone     string    `gorm:"type:varchar(16)" json:"phone"`
	Type      int       `gorm:"column:type"`
	UnionID   string    `gorm:"column:unionid;type:varchar(128)"`
	Created   time.Time `gorm:"column:created"`
	LastLogin time.Time `gorm:"column:lastlogin"`
	IsAdmin   bool      `gorm:"column:isadmin;not null;default:0"`
	IsActive  bool      `gorm:"column:isactive;not null;default:1"`
}

// TableName
func (u User) TableName() string {
	return "user"
}

// WeChatLogin
func (this *UserServiceProvider) WeChatLogin(conn orm.Connection, info *WechatLogin) (*User, error) { // todo
	var (
		err  error
		user User
	)
	db := conn.(*gorm.DB)

	err = db.Where("unionID = ?", info.UnionID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := fmt.Sprintf("%s%d", info.UnionID, r.Intn(10000))
	salt, err := security.SaltHashGenerate(&s)
	if err != nil {
		return nil, err
	}

	user.Type = WeChat
	user.Password = string(salt)
	user.IsActive = true
	user.IsAdmin = false
	user.UnionID = info.UnionID
	user.Created = time.Now()
	user.LastLogin = time.Now()

	err = db.Model(&User{}).Create(&user).Error
	if err != nil {
		return nil, err
	}

	err = db.Where("unionID = ?", info.UnionID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// wechat add a phone number
func (this *UserServiceProvider) AddPhone(conn orm.Connection, id uint32, phone *WechatPhone) error {
	var (
		user User
	)

	db := conn.(*gorm.DB)
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}
	if user.Type == Mobile {
		return ErrAddPhone
	}
	user.Phone = phone.Phone
	return db.Save(&user).Error
}

// Change user information
func (this *UserServiceProvider) ChangeInfo(conn orm.Connection, id uint32, change *ChangeInfo) error {
	var (
		user User
	)

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	if len(user.Avatar) > 0 { // todo
		DeletePicture(user.Avatar)
	}

	user.UserName = change.UserName
	user.Sex = change.Sex
	user.Avatar = change.Avatar

	return db.Save(&user).Error
}

// Register by phone
func (this *UserServiceProvider) PhoneRegister(conn orm.Connection, register *PhoneRegister) error { // todo
	salt, err := security.SaltHashGenerate(&register.Password)
	if err != nil {
		return err
	}

	user := User{
		UserName:  register.UserName,
		Phone:     register.Phone,
		UnionID:   register.Phone, // todo
		Password:  string(salt),
		Type:      Mobile,
		IsAdmin:   false,
		Created:   time.Now(),
		LastLogin: time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&user).Error
}

func (this *UserServiceProvider) PhoneLogin(conn orm.Connection, login *PhoneLogin) (*User, error) {
	var (
		user User
	)
	var updater = make(map[string]interface{})
	updater["lastlogin"] = time.Now()

	db := conn.(*gorm.DB)
	err := db.Where("phone = ?", login.Phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	if !security.SaltHashCompare([]byte(user.Password), &login.Password) {
		return nil, ErrInvalidPass
	}

	err = db.Model(&user).Where("id = ?", user.UserID).Update(updater).Limit(1).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (this *UserServiceProvider) ChangePassword(conn orm.Connection, id uint32, change *ChangePass) (err error) {
	var (
		user User
	)

	tx := conn.(*gorm.DB).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	if !security.SaltHashCompare([]byte(user.Password), &change.OldPass) {
		err = ErrInvalidPass
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

// GetUserByID gets user's information by userId.
func (this *UserServiceProvider) GetUserByID(conn orm.Connection, userID uint32) (*User, error) {
	db := conn.(*gorm.DB)
	user := &User{}

	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
