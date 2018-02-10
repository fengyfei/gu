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
 *     Initial: 2018/02/04        Shi Ruitao
 */

package category

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

const (
	Unused = 1
	InUse  = 2
)

var (
	Service *serviceProvider
)

type (
	Category struct {
		ID          uint64    `gorm:"column:id"`
		Category    string    `gorm:"column:category"`
		Description string    `gorm:"column:description"`
		ParentID    uint64    `gorm:"column:parentid"`
		Status      uint8     `gorm:"column:status"`
		Created     time.Time `gorm:"column:created"`
	}

	CategoryData struct {
		ID          uint64    `gorm:"column:id"`
		Category    string    `gorm:"column:category"`
		Description string    `gorm:"column:description"`
		ParentID    uint64    `gorm:"column:parentid"`
		Status      uint8     `gorm:"column:status"`
		Created     time.Time `gorm:"column:created"`
	}
)

type (
	Add struct {
		Category    string `json:"category" validate:"required,alphanumunicode,max=12"`
		Description string `json:"description" validate:"required,max=50"`
		Status      uint8  `json:"status"`
		ParentID    uint64 `json:"parentid"`
	}

	Modify struct {
		ID          uint64 `json:"id"`
		Category    string `json:"category" validate:"required,alphanumunicode,max=12"`
		Description string `json:"description" validate:"required,max=50"`
		Status      uint8  `json:"status"`
		ParentID    uint64 `json:"parentid"`
	}
)

func (Category) TableName() string {
	return "category"
}

// Add a category
func (sp *serviceProvider) Add(conn orm.Connection, add *Add) error {
	category := Category{
		Category:    add.Category,
		Description: add.Description,
		ParentID:    add.ParentID,
		Status:      add.Status,
		Created:     time.Now(),
	}

	db := conn.(*gorm.DB)

	return db.Create(&category).Error
}

func (sp *serviceProvider) GetCategory(conn orm.Connection, pid uint64) ([]Category, error) {
	var list []Category

	db := conn.(*gorm.DB).Exec("USE shop")
	res := db.Table("categories").Where("status > ? AND parentid = ?", 0, pid).Scan(&list)

	return list, res.Error
}

// Modify the category
func (sp *serviceProvider) Modify(conn orm.Connection, modify *Modify) error {
	var category Category

	db := conn.(*gorm.DB)

	err := db.Where("id = ?", modify.ID).First(&category).Error
	if err != nil {
		return err
	}

	category.Category = modify.Category
	category.Description = modify.Description
	category.ParentID = modify.ParentID
	category.Status = modify.Status

	return db.Save(&category).Error
}
