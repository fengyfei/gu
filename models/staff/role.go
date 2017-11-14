/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 *     Initial: 2017/11/08        Jia Chenhui
 */

package staff

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

const (
	roleName = "role"
)

type Role struct {
	Id      int16  `gorm:"primary_key;auto_increment"`
	Name    string `gorm:"type:varchar(30);not null;unique"`
	Intro   string `gorm:"type:varchar(255)"`
	Active  bool
	Created *time.Time
}

// TableName returns table name in database.
func (Role) TableName() string {
	return roleName
}

// CreateRole create role information.
func (sp *serviceProvider) CreateRole(conn orm.Connection, name, intro *string) error {
	now := time.Now()

	role := &Role{}
	value := &Role{
		Name:    *name,
		Intro:   *intro,
		Active:  true,
		Created: &now,
	}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(role).Create(value).Error
}

// ModifyRole modify role information.
func (sp *serviceProvider) ModifyRole(conn orm.Connection, id int16, name, intro *string) error {
	role := &Role{}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(role).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  *name,
		"intro": *intro,
	}).Limit(1).Error
}

// ModifyRoleActive modify role status.
func (sp *serviceProvider) ModifyRoleActive(conn orm.Connection, id int16, active bool) error {
	role := &Role{}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(role).Where("id = ?", id).Update("active", active).Limit(1).Error
}

// RoleList list all on the active role.
func (sp *serviceProvider) RoleList(conn orm.Connection) ([]Role, error) {
	r := &Role{}
	list := []Role{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(r).Where("active = true").Find(&list).Error

	if err != nil {
		return list, err
	}

	return list, nil
}

// GetRoleByID get one role detail information.
func (sp *serviceProvider) GetRoleByID(conn orm.Connection, id int16) (*Role, error) {
	role := &Role{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(role).Where("id = ? AND active = true", id).First(role).Error

	if err != nil {
		return nil, err
	}

	return role, nil
}
