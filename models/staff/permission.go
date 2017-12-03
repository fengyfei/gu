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
	permissionName = "permission"
)

// Permission represents permission of URL.
type Permission struct {
	URL     string `gorm:"primary_key" sql:"type:varchar(255) not null default ''"`
	RoleId  int16  `gorm:"column:roleid;primary_key" sql:"type:smallint(5) not null default 0"`
	Created *time.Time
}

// TableName returns table name in database.
func (Permission) TableName() string {
	return permissionName
}

// AddPermission create an associated record of the specified URL and role.
func (sp *serviceProvider) AddURLPermission(conn orm.Connection, url *string, rid int16) error {
	now := time.Now()
	r := &Role{}
	permission := &Permission{}
	value := &Permission{
		URL:     *url,
		RoleId:  rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errRoleInactive
		goto finish
	}

	err = txn.Model(permission).Create(value).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// RemovePermission remove the associated records of the specified URL and role.
func (sp *serviceProvider) RemoveURLPermission(conn orm.Connection, url *string, rid int16) error {
	r := &Role{}
	permission := &Permission{}
	condition := &Permission{
		URL:    *url,
		RoleId: rid,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errRoleInactive
		goto finish
	}

	err = txn.Model(permission).Delete(condition).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// URLPermissions lists all the roles of the specified URL.
func (sp *serviceProvider) URLPermissions(conn orm.Connection, url *string) (map[int16]bool, error) {
	permission := &Permission{}
	plist := []Permission{}
	result := make(map[int16]bool)

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(permission).Where("url = ?", *url).Find(&plist).Error

	if err != nil {
		return nil, err
	}

	for _, p := range plist {
		result[p.RoleId] = true
	}

	return result, nil
}

// Permissions lists all the roles.
func (sp *serviceProvider) Permissions(conn orm.Connection) ([]Permission, error) {
	permission := &Permission{}
	plist := []Permission{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(permission).Find(&plist).Error

	if err != nil {
		return nil, err
	}

	return plist, nil
}

func CreateAdminPermission(conn orm.Connection) error {
	url := "http://127.0.0.1:21000/api/v1/staff/create"
	rid := int16(1)
	now := time.Now()

	r := &Role{}
	permission := &Permission{}
	value := &Permission{
		URL:     url,
		RoleId:  rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errRoleInactive
		goto finish
	}

	err = txn.Model(permission).Create(value).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}
