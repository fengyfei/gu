package user

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type User struct {
	ID        string `gorm:"primary_key;unique;type:varchar(128)"`
	UserName  string `gorm:"type:varchar(30)"`
	Address   string `gorm:"type:varchar(128)"`
	Phone     string `gorm:"unique"`
	Type      string `gorm:"type:varchar(30)"`
	CreatedAt *time.Time
}

func (this *serviceProvider) WechatLogin(conn orm.Connection, unionId *string) (string, error) {

	user := &User{}
	user.ID = *unionId

	db := conn.(*gorm.DB).Exec("USE user")

	err := db.Where("id = ?", *unionId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user.ID, db.Model(user).Create(&User{}).Error
	}

	return user.ID, nil
}
