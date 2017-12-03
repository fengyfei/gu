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
	relationName = "relation"
)

var (
	errStaffInactive = errors.New("the staff is not activated")
	errRoleInactive  = errors.New("the role is not activated")
)

// Staff role relation table.
type Relation struct {
	StaffId int32 `gorm:"column:staffid;primary_key" sql:"type:int(11) not null default 0"`
	RoleId  int16 `gorm:"column:roleid;primary_key" sql:"type:smallint(5) not null default 0"`
	Created *time.Time
}

// TableName returns table name in database.
func (Relation) TableName() string {
	return relationName
}

// AddRole add a role to staff.
func (sp *serviceProvider) AddRole(conn orm.Connection, sid int32, rid int16) error {
	now := time.Now()
	s := &Staff{}
	r := &Role{}
	relation := &Relation{}
	value := &Relation{
		StaffId: sid,
		RoleId:  rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errStaffInactive
		goto finish
	}

	err = txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errRoleInactive
		goto finish
	}

	err = txn.Model(relation).Create(value).Error

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
func (sp *serviceProvider) RemoveRole(conn orm.Connection, sid int32, rid int16) error {
	s := &Staff{}
	r := &Role{}
	relation := &Relation{}
	condition := &Relation{
		StaffId: sid,
		RoleId:  rid,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errStaffInactive
		goto finish
	}

	err = txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	err = txn.Model(relation).Delete(condition).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// AssociatedRoles list all the roles of the specified staff.
func (sp *serviceProvider) AssociatedRoles(conn orm.Connection, sid int32) (map[int16]bool, error) {
	s := &Staff{}
	relation := &Relation{}
	rlist := []Relation{}
	result := make(map[int16]bool)

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errStaffInactive
		goto finish
	}

	err = txn.Model(relation).Where("staffid = ?", sid).Find(&rlist).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	for _, r := range rlist {
		result[r.RoleId] = true
	}

	return result, err
}

func CreateAdminRelation(conn orm.Connection) error {
	sid := int32(1)
	rid := int16(1)
	now := time.Now()

	s := &Staff{}
	r := &Role{}
	relation := &Relation{}
	value := &Relation{
		StaffId: sid,
		RoleId:  rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(s).Where("id = ?", sid).First(s).Error
	if err != nil {
		goto finish
	}

	if !s.Active || s.Resigned {
		err = errStaffInactive
		goto finish
	}

	err = txn.Model(r).Where("id = ?", rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errRoleInactive
		goto finish
	}

	err = txn.Model(relation).Create(value).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}
