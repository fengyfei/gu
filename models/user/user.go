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

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/libs/mongo"
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

// PhoneRegister
type PhoneRegister struct {
	Phone    string `json:"phone" validate:"required,alphanum,len=11"`
	Password string `json:"password" validate:"required,min=6,max=30"`
	Name     string `json:"name" validate:"required,alphaunicode,min=2,max=30"`
}

// PhoneLogin
type PhoneLogin struct {
	Phone    string `json:"phone" validate:"required,alphanum,len=11"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

// ChangePass
type ChangePass struct {
	OldPass string `json:"oldPass" validate:"required,min=6,max=30"`
	NewPass string `json:"newPass" validate:"required,min=6,max=30"`
}

var (
	// UserServer
	UserServer     *UserServiceProvider
	session        *mongo.Connection
	// DefaultImage
	DefaultImage   = ""
	// ErrInvalidUser
	ErrInvalidUser = errors.New("User doesn't exists.")
	// ErrInvalidPass
	ErrInvalidPass = errors.New("the password error.")
	// ErrUserExists
	ErrUserExists  = errors.New("User already exists.")
)

func init() {
	const collection = "avatar"
	url := conf.BBSConfig.MongoURL + "/" + "bbs"
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)

	session = mongo.NewConnection(s, "bbs", collection)
	UserServer = &UserServiceProvider{}
}

// User represents users information
type User struct {
	UserID     uint64    `gorm:"column:id;primary_key;auto_increment" json:"userID"`
	UserName   string    `gorm:"column:UserName;size:16"`
	Password   string    `gorm:"column:Password;type:varchar(128)" json:"password" validate:"required,alphanum,min=6,max=30"`
	Phone      string    `gorm:"type:varchar(16)" json:"phone" validate:"required,numeric,len=11"`
	UnionID    string    `gorm:"column:UnionID;type:varchar(128)"`
	AvatarID   string    `gorm:"column:AvatarID;type:varchar(128)"`
	Created    time.Time `gorm:"column:Created"`
	LastLogin  time.Time `gorm:"column:LastLogin"`
	Type       int       `grom:"column:Type"`
	IsActive   bool      `gorm:"column:IsActive;not null;default:1"`
	ArticleNum int64     `gorm:"column:ArticleNum;not null;default:0"`
}

// TableName
func (u User) TableName() string {
	return "users"
}

// Avatar
type Avatar struct {
	AvatarID bson.ObjectId `bson:"_id,omitempty"`
	UserID   uint64        `bson:"UserID"`
	Avatar   string        `bson:"Avatar"`
}

// WeChatLogin
func (this *UserServiceProvider) WeChatLogin(username, unionID *string) (uint64, error) {
	conn, err := initialize.Pool.Get()
	if err != nil {
		return 0, err
	}

	user := &User{}
	db := conn.(*gorm.DB).Exec("USE user")
	err = db.Where("unionID = ?", *unionID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			s := fmt.Sprintf("%s%s", unionID, r.Intn(10000))
			salt, err := security.SaltHashGenerate(&s)
			if err != nil {
				return 0, err
			}

			user.UserName = *username
			user.Type = WeChat
			user.Password = string(salt)
			user.AvatarID = DefaultImage
			user.IsActive = true
			user.UnionID = *unionID
			user.Created = time.Now()
			user.LastLogin = time.Now()

			err = db.Model(&User{}).Create(&user).Error
			if err != nil {
				return 0, err
			}

			err = db.Where("unionID = ?", *unionID).First(&user).Error
			return user.UserID, nil
		}
		return 0, err
	}

	return user.UserID, nil
}

// ChangeName
func (this *UserServiceProvider) ChangeName(userID uint64, newname *string) error {
	var (
		err  error
		user User
	)

	conn, err := initialize.Pool.Get()
	if err != nil {
		return err
	}

	db := conn.(*gorm.DB).Exec("USE user")

	err = db.Where("id = ?", 1001).Find(&user).Error
	if err != nil {
		return err
	}
	user.UserName = *newname

	return db.Model(&user).Save(&user).Error
}

// ChangeAvatar
func (this *UserServiceProvider) ChangeAvatar(userID uint64, avatar *string) (string, error) {
	var res Avatar
	updater := bson.M{"$set": bson.M{
		"avatar": *avatar,
	}}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Update(bson.M{"_id": bson.ObjectId(userID)}, updater)
	if err != nil {
		return "", nil
	}

	query := bson.M{"userID": userID}
	err = conn.GetUniqueOne(query, &res)

	return res.AvatarID.Hex(), err
}

// PhoneRegister
func (this *UserServiceProvider) PhoneRegister(p *PhoneRegister) error {
	salt, err := security.SaltHashGenerate(&p.Password)
	if err != nil {
		return err
	}

	conn, err := initialize.Pool.Get()
	if err != nil {
		return err
	}

	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}

	err = db.Where("Phone = ?", p.Phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user.UserName = p.Name
		user.Type = Mobile
		user.Phone = p.Phone
		user.Password = string(salt)
		user.AvatarID = DefaultImage
		user.IsActive = true
		user.Created = time.Now()
		user.LastLogin = time.Now()

		return db.Model(&User{}).Create(&user).Error
	}

	return ErrUserExists
}

// PhoneLogin
func (this *UserServiceProvider) PhoneLogin(p *PhoneLogin) (uint64, error) {
	conn, err := initialize.Pool.Get()
	if err != nil {
		return 0, err
	}

	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}

	err = db.Where("Phone = ?", p.Phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return 0, ErrInvalidUser
	}

	if !security.SaltHashCompare([]byte(user.Password), &p.Password) {
		return 0, ErrInvalidPass
	}

	return user.UserID, err
}

// ChangePassword
func (this *UserServiceProvider) ChangePassword(id uint, oldPass, newPass *string) error {
	conn, err := initialize.Pool.Get()
	if err != nil {
		return err
	}

	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}

	err = db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}

	if !security.SaltHashCompare([]byte(user.Password), oldPass) {
		return ErrInvalidPass
	}

	salt, err := security.SaltHashGenerate(newPass)
	if err != nil {
		return err
	}
	user.Password = string(salt)
	return db.Save(&user).Error
}

// GetUserByID gets user's information by userId.
func (this *UserServiceProvider) GetUserByID(userId uint64) (*User, error) {
	conn, err := initialize.Pool.Get()
	if err != nil {
		return nil, err
	}

	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}

	err = db.Where("id = ?", userId).First(&user).Error

	return user, err
}
