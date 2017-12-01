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
 *     Initial: 2017/11/30        Wang RiYu
 */

package panel

import (
  "time"
  "github.com/fengyfei/gu/libs/orm"
  "github.com/jinzhu/gorm"
)

const type1 int8 = 1 // promotion && flash sale
const type2 int8 = 2 // recommends && advertising
const type3 int8 = 3 // second-hand && other things

type serviceProvider struct{}

var (
  Service *serviceProvider
)

type (
  Panel struct {
    ID       uint   `gorm:"primary_key;AUTO_INCREMENT"`
    Title    string `gorm:"type:varchar(50)"`
    Desc     string `gorm:"type:varchar(100)"`
    Type     int8   `gorm:"type:TINYINT;not null"`
    Status   int8   `gorm:"type:TINYINT;default:1"`
    Sequence int    `gorm:"unique_index;not null"`
    Related  string `gorm:"type:varchar(100)"` // wares id
    Created  time.Time
  }

  Detail struct {
    ID      uint   `gorm:"primary_key;AUTO_INCREMENT"`
    Belong  uint   `gorm:"unique_index;not null"`
    Picture string `gorm:"type:varchar(100)"`
    Content string `gorm:"type:LONGTEXT"`
    Created time.Time
  }

  PanelReq struct {
    Title    string `json:"title"`
    Desc     string `json:"desc"`
    Type     int8   `json:"type" validate:"required, eq=1|eq=2|eq=3"`
    Related  string `json:"related"`
    Sequence int    `json:"sequence"`
  }
)

// add panel
func (sp *serviceProvider) CreatePanel(conn orm.Connection, panelReq PanelReq) error {
  panel := &Panel{}
  panel.Title = panelReq.Title
  panel.Desc = panelReq.Desc
  panel.Type = panelReq.Type
  panel.Sequence = panelReq.Sequence
  panel.Related = panelReq.Related
  panel.Created = time.Now()

  db := conn.(*gorm.DB).Exec("USE shop")
  err := db.Model(&Panel{}).Create(panel).Error

  return err
}
