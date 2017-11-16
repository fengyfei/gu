/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/11/13        Wang RiYu
 */

package category

import (
  "time"
  "github.com/fengyfei/gu/libs/orm"
  "github.com/jinzhu/gorm"
)

type serviceProvider struct{}

var (
  Service *serviceProvider
)

type Category struct {
  ID        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
  Name      string    `gorm:"type:varchar(50);not null" json:"name"`
  Desc      string    `gorm:"type:varchar(100);not null" json:"desc"`
  ParentID  uint      `gorm:"not null" json:"parentId"`
  CreatedAt time.Time `json:"createdTime"`
}

func (sp *serviceProvider) AddCategory(conn orm.Connection, name *string, desc *string, parentID *uint) error {
  category := &Category{}
  category.Name = *name
  category.Desc = *desc
  category.ParentID = *parentID

  db := conn.(*gorm.DB).Exec("USE shop")
  err := db.Model(&Category{}).Create(category).Error

  return err
}

func (sp *serviceProvider) GetCategory(conn orm.Connection) ([]Category, error) {
  var list []Category

  db := conn.(*gorm.DB).Exec("USE shop")
  res := db.Table("categories").Where("parent_id = ?", 0).Scan(&list)

  return list, res.Error
}
