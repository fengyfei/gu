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

// User represents users information
type User struct {
	UserID    uint32    `gorm:"column:id;primary_key;auto_increment" json:"user_id"`
	UserName  string    `gorm:"column:username;type:varchar(128)" json:"user_name"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Sex       uint8     `gorm:"column:sex" json:"sex"` // 0 -> male, 1 -> female
	Password  string    `gorm:"column:password;type:varchar(128)" json:"password"`
	Phone     string    `gorm:"type:varchar(16)" json:"phone"`
	Type      int       `gorm:"column:type"` // 0 -> Wechat, 1 -> Mobile
	UnionID   string    `gorm:"column:unionid;type:varchar(128)" json:"union_id"`
	Created   time.Time `gorm:"column:created"`
	LastLogin time.Time `gorm:"column:lastlogin"`
	IsAdmin   bool      `gorm:"column:isadmin;not null;default:0"`
	IsActive  bool      `gorm:"column:isactive;not null;default:1"`
}

// TableName
func (u User) TableName() string {
	return "user"
}

// WeChatLogin login by wechat
func (this *UserServiceProvider) WeChatLogin(conn orm.Connection, UnionID string) (*User, error) {
	var (
		err  error
		user User
	)
	db := conn.(*gorm.DB)

	err = db.Where("unionID = ?", "1234").First(&user).Error
	if err == nil {
		lastLogin(conn, user.UserID)
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := fmt.Sprintf("%s%d", UnionID, r.Intn(10000))
	salt, err := security.SaltHashGenerate(&s)
	if err != nil {
		return nil, err
	}
	time := time.Now()
	user.UserName = "name"
	user.Type = WeChat
	user.Password = string(salt)
	user.IsActive = true
	user.IsAdmin = false
	user.UnionID = "1234"
	user.Avatar = ""
	user.Created = time
	user.LastLogin = time

	err = db.Model(&User{}).Create(&user).Error
	if err != nil {
		return nil, err
	}

	err = db.Where("unionID = ?", "1234").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// AddPhone wechat add a phone number
func (this *UserServiceProvider) AddPhone(conn orm.Connection, id uint32, phone string) error {
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

	user.Phone = phone
	return db.Save(&user).Error
}

// ChangeInfo change user information
func (this *UserServiceProvider) ChangeInfo(conn orm.Connection, id uint32, userName string, sex uint8) error { // todo
	var (
		user User
	)

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	user.UserName = userName
	user.Sex = sex

	return db.Save(&user).Error
}

// ChangeAvatar change avatar
func (this *UserServiceProvider) ChangeAvatar(conn orm.Connection, userID uint32, avatar string) error {
	updater := make(map[string]interface{})
	updater["avatar"] = avatar

	db := conn.(*gorm.DB)
	return db.Table("user").Where("id = ?", userID).Update(updater).Limit(1).Error
}

// PhoneRegister register by phone
func (this *UserServiceProvider) PhoneRegister(conn orm.Connection, userName, phone, password string) error {
	salt, err := security.SaltHashGenerate(&password)
	if err != nil {
		return err
	}

	user := User{
		UserName:  userName,
		Phone:     phone,
		Password:  string(salt),
		Type:      Mobile,
		IsAdmin:   false,
		Avatar:    "",
		Created:   time.Now(),
		LastLogin: time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&user).Error
}

// PhoneLogin login by phone
func (this *UserServiceProvider) PhoneLogin(conn orm.Connection, phone, password string) (*User, error) {
	var (
		user User
	)

	db := conn.(*gorm.DB)
	err := db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	if !security.SaltHashCompare([]byte(user.Password), &password) {
		return nil, ErrInvalidPass
	}

	err = lastLogin(conn, user.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePassword change password
func (this *UserServiceProvider) ChangePassword(conn orm.Connection, id uint32, oldPass, newPass string) (err error) {
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

	if !security.SaltHashCompare([]byte(user.Password), &oldPass) {
		err = ErrInvalidPass
		return err
	}

	salt, err := security.SaltHashGenerate(&newPass)
	if err != nil {
		return err
	}

	user.Password = string(salt)

	err = tx.Save(&user).Limit(1).Error
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

func lastLogin(conn orm.Connection, id uint32) error {
	updater := make(map[string]interface{})
	updater["lastlogin"] = time.Now()

	db := conn.(*gorm.DB)
	return db.Table("user").Where("id = ?", id).Update(updater).Limit(1).Error
}
