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
 *     Initial: 2018/03/07        Tong Yuehong
 */

package category

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider.
	Service *serviceProvider
)

type (
	// Category represents the class of wares.
	Category struct {
		ID       uint32    `gorm:"column:id"`
		Category string    `gorm:"column:category"`
		ParentID uint32    `gorm:"column:parent_id"`
		IsActive bool      `gorm:"column:is_active"`
		Created  time.Time `gorm:"column:created"`
	}

	// Info represents the category's information.
	Info struct {
		Category string `json:"category" validate:"required,min=2,max=12"`
		ParentID uint32 `json:"parent_id"`
	}
)

// TableName returns table name in database.
func (Category) TableName() string {
	return "category"
}

// Add adds a new category.
func (sp *serviceProvider) Add(conn orm.Connection, add *Info) error {
	category := Category{
		Category: add.Category,
		ParentID: add.ParentID,
		IsActive: true,
		Created:  time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&category).Error
}

// GetMainCategory gets all the parentid.
func (sp *serviceProvider) GetMainCategory(conn orm.Connection) ([]Category, error) {
	var (
		list []Category
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Table("category").Where("is_active = ? AND parent_id = ?", true, 0).Find(&list).Error

	return list, err
}

// GetSubCategory gets subcategory which parentid is pid.
func (sp *serviceProvider) GetSubCategory(conn orm.Connection, pid uint32) ([]Category, error) {
	var (
		list []Category
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Table("category").Where("is_active = ? AND parent_id = ?", true, pid).Find(&list).Error

	return list, err
}

// Delete deletes the category.
func (sp *serviceProvider) Delete(conn orm.Connection, id uint32) error {
	db := conn.(*gorm.DB)

	return db.Table("category").Where("id = ?", id).Update("is_active", false).Error
}

// Modify modify category's information.
func (sp *serviceProvider) Modify(conn orm.Connection, id uint32, add *Info) error {
	db := conn.(*gorm.DB)

	return db.Table("category").Where("id = ?", id).Update(map[string]interface{}{
		"category":  add.Category,
		"parent_id": add.ParentID,
	}).Error
}
