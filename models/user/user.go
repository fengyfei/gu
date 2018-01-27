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
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/libs/orm"
)

type UserServiceProvider struct{}

const (
	Wechat = iota
	Github
	Mobile
)

var (
	UserServer *UserServiceProvider
	session    *mongo.Connection
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
	UserId     uint64     `gorm:"primary_key;auto_increment"`
	WechatCode string     `gorm:"unique;type:varchar(128)"`
	UnionId    string     `gorm:"unique;type:varchar(128)"`
	Username   string     `gorm:"unique;size:16"`
	AvatarId   string     `gorm:"type:varchar(32)"`
	Created    *time.Time `gorm:"column:created"`
	LastLogin  *time.Time `gorm:"column:lastlogin"`
	Type       int        `grom:"column:type"`
	Status     bool       `gorm:"not null;default:1"`
	ThemeNum   int64      `gorm:"not null;default:0"`
	ArticleNum int64      `gorm:"not null;default:0"`
}

// Avatar
type Avatar struct {
	AvatarId bson.ObjectId `bson:"_id,omitempty"`
	UserId   uint64        `bson:"UserId"`
	Avatar   string        `bson:"Avatar"`
}

// WechatLogin
func (this *UserServiceProvider) WechatLogin(conn orm.Connection, username, wechatCode, unionId *string) (uint64, error) {

	user := &User{}
	res := &User{}
	user.Username = *username
	user.WechatCode = *wechatCode
	user.UnionId = *unionId
	user.Type = Wechat

	db := conn.(*gorm.DB).Exec("USE user")

	err := db.Where("union_id = ?", *unionId).First(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// not found, create new user
			err = db.Model(&User{}).Create(&user).Error
			// insert user's avatar
			UserServer.ChangeAvatar(&res.UserId, &res.AvatarId)
			if err != nil {
				return 0, err
			}
			return user.UserId, nil
		}
		return 0, err
	}
	return res.UserId, nil
}

// ChangeUsername
func (this *UserServiceProvider) ChangeUsername(conn orm.Connection, UserId uint64, newname *string) error {
	db := conn.(*gorm.DB).Exec("USE user")
	user := &User{}
	return db.Where("id = ?", UserId).First(&user).Error
}

// ChangeAvatar
func (this *UserServiceProvider) ChangeAvatar(userId *uint64, avatar *string) (string, error) {
	var res Avatar
	updater := bson.M{"$set": bson.M{
		"avatar": *avatar,
	}}
	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Update(bson.M{"_id": bson.ObjectId(*userId)}, updater)

	if err != nil {
		return "", nil
	}

	query := bson.M{"userId": userId}
	err = conn.GetUniqueOne(query, &res)
	return res.AvatarId.Hex(), err
}
