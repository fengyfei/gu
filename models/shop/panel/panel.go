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

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
)

type (
	Panel struct {
		ID        uint32    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
		Title     string    `gorm:"type:varchar(50)" json:"title"`
		Desc      string    `gorm:"type:varchar(100)" json:"desc"`
		Type      uint8     `gorm:"type:TINYINT;not null" json:"type"` // 1 promotion && flash sale;2 recommends && advertising;3 second-hand && other things
		Status    uint8     `gorm:"type:TINYINT;default:1" json:"status"`
		Sequence  uint8     `gorm:"unique_index;not null" json:"sequence"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	}

	PanelsPage struct {
		ID        uint          `json:"id"`
		Title     string        `json:"title"`
		Desc      string        `json:"desc"`
		Type      uint8         `json:"type"`
		Status    uint8         `json:"status"`
		Sequence  int           `json:"sequence"`
		CreatedAt time.Time     `json:"createdAt"`
		Content   []interface{} `json:"content"`
	}

	Detail struct {
		ID        uint      `gorm:"primary_key;AUTO_INCREMENT"`
		Belong    uint      `gorm:"unique_index;not null"`
		Picture   string    `gorm:"type:varchar(100)"`
		Content   string    `gorm:"type:LONGTEXT"`
		UpdatedAt time.Time `json:"updatedAt"`
		CreatedAt time.Time `json:"createdAt"`
	}

	PanelReq struct {
		Title    string `json:"title" validate:"required"`
		Desc     string `json:"desc"`
		Type     uint8  `json:"type" validate:"eq=1|eq=2|eq=3"`
		Sequence uint8  `json:"sequence"`
	}

	PromotionReq struct {
		Belong  uint   `json:"belong" validate:"required"`
		Content string `json:"content" validate:"required"`
	}

	RecommendReq struct {
		Belong  uint   `json:"belong" validate:"required"`
		Picture string `json:"picture" validate:"required"`
		Content string `json:"content"`
	}
)

// add panel
func (sp *serviceProvider) CreatePanel(conn orm.Connection, panelReq PanelReq) error {
	panel := &Panel{}
	panel.Title = panelReq.Title
	panel.Desc = panelReq.Desc
	panel.Type = panelReq.Type
	panel.Sequence = panelReq.Sequence

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Model(&Panel{}).Create(panel).Error

	return err
}

// add promotion list
func (sp *serviceProvider) AddPromotionList(conn orm.Connection, promotionReq PromotionReq) error {
	promotion := &Detail{}
	promotion.Belong = promotionReq.Belong
	promotion.Content = promotionReq.Content

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Model(&Detail{}).Create(promotion).Error

	return err
}

// add recommend
func (sp *serviceProvider) AddRecommend(conn orm.Connection, recommendReq RecommendReq) error {
	recommend := &Detail{}
	recommend.Belong = recommendReq.Belong
	recommend.Picture = recommendReq.Picture
	recommend.Content = recommendReq.Content

	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Model(&Detail{}).Create(recommend).Error

	return err
}

// get panels
func (sp *serviceProvider) GetPanels(conn orm.Connection) ([]PanelsPage, error) {
	var (
		list []PanelsPage
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res := db.Table("panels").Where("status > ?", 0).Scan(&list)

	return list, res.Error
}

// get detail of panel
func (sp *serviceProvider) GetDetail(conn orm.Connection, id uint) (Detail, error) {
	var (
		detail Detail
	)

	db := conn.(*gorm.DB).Exec("USE shop")
	res := db.Table("details").Where("belong = ?", id).First(&detail)

	return detail, res.Error
}
