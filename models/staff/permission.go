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
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

const (
	permissionName = "permission"
)

// Permission represents permission of URL.
type Permission struct {
	URL     string `json:"url" gorm:"primary_key" sql:"type:varchar(255) NOT NULL DEFAULT ''"`
	RoleId  int16  `json:"roleid" gorm:"primary_key" sql:"type:SMALLINT NOT NULL DEFAULT 0"`
	Created *time.Time
}

// TableName returns table name in database.
func (Permission) TableName() string {
	return permissionName
}

// AddPermission create an associated record of the specified URL and role.
func (sp *serviceProvider) AddURLPermission(conn orm.Connection, url *string, rid *int16) error {
	now := time.Now()
	r := &Role{}
	permission := &Permission{}
	value := &Permission{
		URL:     *url,
		RoleId:  *rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(r).Where("id = ?", *rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errors.New("the role is not activated")
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
func (sp *serviceProvider) RemoveURLPermission(conn orm.Connection, url *string, rid *int16) error {
	r := &Role{}
	permission := &Permission{}
	condition := &Permission{
		URL:    *url,
		RoleId: *rid,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(r).Where("id = ?", *rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errors.New("the role is not activated")
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

// PermissionList lists all the roles of the specified URL.
func (sp *serviceProvider) URLPermissionList(conn orm.Connection, url *string) ([]Permission, error) {
	permission := &Permission{}
	result := []Permission{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(permission).Where("url = ?", *url).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}
