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
 *     Initial: 2017/11/16        Wang RiYu
 */

package ware

import (
  "time"
  _ "github.com/fengyfei/gu/libs/orm"
  _ "github.com/jinzhu/gorm"
)

type serviceProvider struct{}

var (
  Service *serviceProvider
)

type Ware struct {
  ID         uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
  Name       string    `gorm:"type:varchar(50);not null"  json:"name"`
  Desc       string    `gorm:"type:varchar(100);not null" json:"desc"`
  CategoryID uint      `gorm:"not null" json:"categoryId"`
  TotalSale  uint      `gorm:"not null" json:"totalSale"`
  Price      float32   `gorm:"not null" json:"price"`
  SalePrice  float32   `gorm:"not null" json:"salePrice"`
  Status     int       `gorm:"not null" json:"status"`
  Size       string    `gorm:"type:varchar(20)"   json:"size"`
  Color      string    `gorm:"type:varchar(20)"   json:"color"`
  Image      string    `gorm:"type:varchar(1000)" json:"image"`
  Inventory  uint      `json:"inventory"`
  Created    time.Time `json:"createdTime"`
}
