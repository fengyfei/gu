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
	trueCategory = 1
)

var (
	Service *serviceProvider
)

type (
	Category struct {
		ID        uint64    `gorm:"column:id"`
		Name      string    `gorm:"column:name"`
		Desc      string    `gorm:"column:desc"`
		ParentID  uint64    `gorm:"column:parentid"`
		Status    int       `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created"`
	}

	CategoryAddReq struct {
		Name     string `json:"name" validate:"required,alphanumunicode,max=6"`
		Desc     string `json:"desc" validate:"required,max=50"`
		ParentID uint64 `json:"parentid"`
	}

	SubCategoryReq struct {
		PID uint64 `json:"pid"`
	}

	Modify struct {
		ID     uint64 `json:"id"`
		Status int    `json:"status"`
	}
)

// add new category - parentID = 0 -> main class, != 0 -> sub class
func (sp *serviceProvider) AddCategory(conn orm.Connection, name *string, desc *string, parentID *uint64) error {
	category := &Category{}
	category.Name = *name
	category.Desc = *desc
	category.Status = trueCategory
	category.ParentID = *parentID
	category.CreatedAt = time.Now()

	db := conn.(*gorm.DB).Exec("USE shop")

	return db.Model(&Category{}).Create(category).Error
}

// get categories by pid
func (sp *serviceProvider) GetCategory(conn orm.Connection, pid uint64) ([]Category, error) {
	var list []Category

	db := conn.(*gorm.DB).Exec("USE shop")
	res := db.Table("categories").Where("status > ? AND parentid = ?", 0, pid).Scan(&list)

	return list, res.Error
}

// Modify category status
func (sp *serviceProvider) ModifyCategory(conn orm.Connection, id uint64, status int) error {
	var category Category

	db := conn.(*gorm.DB)

	return db.Model(&category).Where("id = ?", id).Update("status", status).Error
}
