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

package url

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/role"
)

const (
	filterName = "filter"
)

type Filter struct {
	Id      int16 `json:"id" gorm:"primary_key;auto_increment"`
	URLId   int16 `json:"urlid"`
	RoleId  int16 `json:"roleid"`
	Created *time.Time
}

// TableName returns table name in database.
func (Filter) TableName() string {
	return filterName
}

// AddFilter create an associated record of the specified URL and role.
func (sp *serviceProvider) AddFilter(conn orm.Connection, uid, rid *int16) error {
	now := time.Now()
	u := &URL{}
	r := &role.Role{}
	f := &Filter{}
	value := &Filter{
		URLId:   *uid,
		RoleId:  *rid,
		Created: &now,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(u).Where("id = ?", *uid).First(u).Error
	if err != nil {
		goto finish
	}

	if !u.Active {
		err = errors.New("the url is not activated")
	}

	err = txn.Model(r).Where("id = ?", *rid).First(r).Error
	if err != nil {
		goto finish
	}

	if !r.Active {
		err = errors.New("the role is not activated")
		goto finish
	}

	err = txn.Model(f).Create(value).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// RemoveFilter remove the associated records of the specified URL and role.
func (sp *serviceProvider) RemoveFilter(conn orm.Connection, uid, rid *int16) error {
	u := &URL{}
	r := &role.Role{}
	f := &Filter{}
	condition := &Filter{
		URLId:  *uid,
		RoleId: *rid,
	}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(u).Where("id = ?", *uid).First(u).Error
	if err != nil {
		goto finish
	}

	if !u.Active {
		err = errors.New("the url is not activated")
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

	err = txn.Model(f).Delete(condition).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return err
}

// FilterList lists all the roles of the specified URL.
func (sp *serviceProvider) FilterList(conn orm.Connection, uid *int16) ([]Filter, error) {
	u := &URL{}
	f := &Filter{}
	result := []Filter{}

	db := conn.(*gorm.DB)
	txn := db.Begin().Exec("USE staff")

	err := txn.Model(u).Where("id = ?", *uid).First(u).Error
	if err != nil {
		goto finish
	}

	if !u.Active {
		err = errors.New("the url is not activated")
		goto finish
	}

	err = txn.Model(f).Where("urlid = ?", *uid).Find(&result).Error

finish:
	if err == nil {
		err = txn.Commit().Error
	}

	if err != nil {
		txn.Rollback()
	}

	return result, err
}
