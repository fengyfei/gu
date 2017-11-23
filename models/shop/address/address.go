package address

import (
	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
	"fmt"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type Address struct {
	ID        int    `gorm:"primary_key;auto_increment"`
	UserName  string `gorm:"not null;type:varchar(128)"`
	Address   string `gorm:"type:varchar(128)"`
	IsDefault bool
}

func (this *serviceProvider) Add(conn orm.Connection, userName, address string, isDefault bool) error {
	var (
		err error
	)
	addr := &Address{}
	addr.Address = address
	addr.UserName = userName
	addr.IsDefault = isDefault

	db := conn.(*gorm.DB).Exec("USE shop")

	if !isDefault {
		return db.Model(&Address{}).Create(addr).Error
	}

	another := &Address{}
	err = db.Find(&another, "user_name = ? AND is_default = ?", userName, true).Error
	if err == gorm.ErrRecordNotFound {
		return db.Model(&Address{}).Create(addr).Error
	}

	another.IsDefault = false
	err = db.Save(&another).Error
	if err != nil {
		return err
	}
	return db.Model(&Address{}).Create(addr).Error
}

func (this *serviceProvider) SetDefault(conn orm.Connection, userName string, id int) error {
	var (
		err     error
		addr    Address
		another Address
	)
	db := conn.(*gorm.DB).Exec("USE shop")

	err = db.Find(&addr, "id = ?", id).Error
	if err != nil {
		return err
	}
	if addr.IsDefault {
		return nil
	}
	addr.IsDefault = true
	err = db.Save(&addr).Error
	if err != nil {
		return err
	}

	err = db.Find(&another, "user_name = ? AND id <> ? AND is_default = true", userName, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	fmt.Println(err)
	another.IsDefault = false
	return db.Save(&another).Error
}

func (this *serviceProvider) Modify(conn orm.Connection, id int, address string) error {
	var (
		err  error
		addr Address
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err = db.Find(&addr, "id = ?", id).Error
	if err != nil {
		return err
	}

	addr.Address = address
	return db.Save(&addr).Error
}
