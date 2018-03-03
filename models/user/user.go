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
 *     Initial: 2018/01/21        Chen Yanchen
 *	   Modify:  2018/02/01        Shi Ruitao
 */

package user

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/libs/mongo"
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
	// UserServer
	UserServer = &UserServiceProvider{}

	session    *mongo.Connection

	// ErrInvalidPass
	ErrInvalidPass = errors.New("the password error.")
	// Mobile phone registration cannot add phone number.
	ErrAddPhone = errors.New("Mobile phone registration cannot add phone number.")
)

type (
	PhoneRegister struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Password string `json:"password" validate:"required,min=6,max=30"`
		UserName string `json:"username" validate:"required,alphaunicode,min=2,max=30"`
	}

	PhoneLogin struct {
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
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
		Phone    string `json:"phone" validate:"required,alphanum,len=11"`
		Code     string `json:"code" validate:"required"`
	}

	WechatData struct {
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
	}

	WechatPhone struct {
		Phone string `json:"phone" validate:"required,alphanum,len=11"`
	}

	UserData struct {
		Token    string `json:"token"`
		UserName string `json:"username"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Sex      uint8  `json:"sex"`
	}

	ChangeInfo struct {
		UserName string `json:"username"`
		Sex      uint8  `json:"sex"`
		Avatar   string `json:"avatar"`
	}
	Avatar struct {
		UserID uint32 `bson:"userid,omitempty" json:"id"`
		Avatar string `bson:"avatar" json:"avatar"`
	}
)

// User represents users information
type User struct {
	UserID    uint32    `gorm:"column:id;primary_key;auto_increment" json:"userID"`
	UserName  string    `gorm:"column:username;size:16"`
	Sex       uint8     `gorm:"column:sex"`
	Password  string    `gorm:"column:password;type:varchar(128)" json:"password" validate:"required,alphanum,min=6,max=30"`
	Phone     string    `gorm:"type:varchar(16)" json:"phone" validate:"required,numeric,len=11"`
	UnionID   string    `gorm:"column:unionid;type:varchar(128)"`
	Created   time.Time `gorm:"column:created"`
	LastLogin time.Time `gorm:"column:lastlogin"`
	IsAdmin   bool      `gorm:"column:isadmin"`
	Type      int       `grom:"column:type"`
	IsActive  bool      `gorm:"column:isactive;not null;default:1"`
}

func InitMongo(dbName, tableName, url string) {
	var collection = dbName
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)

	session = mongo.NewConnection(s, tableName, collection)

}

// TableName
func (u User) TableName() string {
	return "users"
}

// WeChatLogin
func (this *UserServiceProvider) WeChatLogin(conn orm.Connection, info *WechatLogin) (*User, *Avatar, error) {
	var err error
	var user User
	var avatar Avatar
	db := conn.(*gorm.DB)

	err = db.Where("unionID = ?", info.UnionID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			s := fmt.Sprintf("%s%d", info.UnionID, r.Intn(10000))
			salt, err := security.SaltHashGenerate(&s)
			if err != nil {
				return nil, nil, err
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
				return nil, nil, err
			}

			err = db.Where("unionID = ?", info.UnionID).First(&user).Error
			if err != nil {
				return nil, nil, err
			}
			c := session.Connect()
			avatar = Avatar{
				UserID: user.UserID,
				Avatar: "null",
			}
			err = c.Insert(avatar)
			if err != nil {
				return nil, nil, err
			}
			return &user, &avatar, err
		}
		return nil, nil, err
	}

	c := session.Connect()
	err = c.GetUniqueOne(bson.M{"userid": user.UserID}, &avatar)
	if err != nil {
		return nil, nil, err
	}
	return &user, &avatar, nil
}

// wechat add a phone number
func (this *UserServiceProvider) AddPhone(conn orm.Connection, id uint32, phone *WechatPhone) error {
	var user User
	db := conn.(*gorm.DB)
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}
	if user.Type == 1 {
		return ErrAddPhone
	}
	user.Phone = phone.Phone
	return db.Save(&user).Error
}

// Change user information
func (this *UserServiceProvider) ChangeInfo(conn orm.Connection, id uint32, change *ChangeInfo) error {
	var user User
	var avatar = &Avatar{
		UserID: id,
		Avatar: change.Avatar,
	}

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	c := session.Connect()
	err = c.Update(bson.M{"userid": id}, avatar)
	if err != nil {
		return err
	}

	user.UserName = change.UserName
	user.Sex = change.Sex

	return db.Save(&user).Error
}

// Register by phone
func (this *UserServiceProvider) PhoneRegister(conn orm.Connection, register *PhoneRegister) error {
	var avatar Avatar

	salt, err := security.SaltHashGenerate(&register.Password)
	if err != nil {
		return err
	}

	user := User{
		UserName:  register.UserName,
		Phone:     register.Phone,
		UnionID:   register.Phone,
		Password:  string(salt),
		Type:      Mobile,
		IsAdmin:   false,
		Created:   time.Now(),
		LastLogin: time.Now(),
	}

	db := conn.(*gorm.DB)
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	c := session.Connect()
	avatar = Avatar{
		UserID: user.UserID,
		Avatar: "null",
	}
	err = c.Insert(avatar)
	if err != nil {
		return err
	}
	return err
}

func (this *UserServiceProvider) PhoneLogin(conn orm.Connection, login *PhoneLogin) (*User, *Avatar, error) {
	var user User
	var avatar Avatar
	var updater = make(map[string]interface{})
	updater["lastlogin"] = time.Now()

	db := conn.(*gorm.DB)
	err := db.Where("phone = ?", login.Phone).First(&user).Error
	if err != nil {
		return nil, nil, err
	}
	if !security.SaltHashCompare([]byte(user.Password), &login.Password) {
		return nil, nil, ErrInvalidPass
	}

	err = db.Model(&user).Where("id = ?", user.UserID).Update(updater).Limit(1).Error
	if err != nil {
		return nil, nil, err
	}

	c := session.Connect()
	err = c.GetUniqueOne(bson.M{"userid": user.UserID}, &avatar)
	if err != nil {
		return nil, nil, err
	}
	return &user, &avatar, nil
}

func (this *UserServiceProvider) ChangePassword(conn orm.Connection, id uint32, change *ChangePass) error {
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
func (this *UserServiceProvider) GetUserByID(conn orm.Connection, userId uint32) (*User, error) {
	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}

	err := db.Where("id = ?", userId).First(&user).Error

	return user, err
}
