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
	"github.com/fengyfei/gu/models/role"
)

const (
	roleManagerName = "rolemanager"
)

// Staff role management table.
type RoleManager struct {
	Id      int32 `json:"id" gorm:"primary_key;auto_increment"`
	StaffID int32 `json:"staffid"`
	RoleID  int16 `json:"roleid"`
	Created *time.Time
}

// TableName returns table name in database.
func (RoleManager) TableName() string {
	return roleManagerName
}

// AddRole add a role to staff.
func (sp *serviceProvider) AddRole(conn orm.Connection, sid *int32, rid *int16) error {
	now := time.Now()
	s := &Staff{}
	r := &role.Role{}
	rm := &RoleManager{
		StaffID: *sid,
		RoleID:  *rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", *sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errors.New("the staff is not activated")
		goto finish
	}

	err = txn.Model(r).Where("id = ?", *rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errors.New("the role is not activated")
		goto finish
	}

	err = txn.Create(rm).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// RemoveRole remove role from staff.
func (sp *serviceProvider) RemoveRole(conn orm.Connection, sid *int32, rid *int16) error {
	s := &Staff{}
	r := &role.Role{}
	rm := &RoleManager{}
	condition := &RoleManager{
		StaffID: *sid,
		RoleID:  *rid,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", *sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errors.New("the staff is not activated")
		goto finish
	}

	err = txn.Model(r).Where("id = ?", *rid).First(r).Error
	if err != nil {
		goto finish
	}

	err = txn.Model(rm).Delete(condition).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// RoleList lists all the roles of the specified staff.
func (sp *serviceProvider) RoleList(conn orm.Connection, sid *int32) ([]role.Role, error) {
	s := &Staff{}
	rm := &RoleManager{}
	result := []RoleManager{}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", *sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errors.New("the staff is not activated")
		goto finish
	}

	err = txn.Model(rm).Where("staffid = ?", *sid).Find(&result).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return result, err
}