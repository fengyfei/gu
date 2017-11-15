package user

import (
	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
	"time"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type User struct {
	ID        int32  `gorm:"primary_key;auto_increment"`
	UserName  string `gorm:"unique;type:varchar(128)"`
	NickName  string `gorm:"type:varchar(30)"`
	Phone     string `gorm:"unique"`
	Type      string `gorm:"type:varchar(30)"`
	CreatedAt *time.Time
}

func (this *serviceProvider) WechatLogin(conn orm.Connection, unionId *string) (string, error) {

	user := &User{}
	user.UserName = *unionId

	db := conn.(*gorm.DB).Exec("USE user")

	err := db.Where("id = ?", *unionId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user.UserName, db.Model(user).Create(&User{}).Error
	}

	return user.UserName, nil
}
