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
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm"
)

const (
	urlName = "url"
)

type URL struct {
	Id      int16  `json:"id" gorm:"primary_key;auto_increment"`
	Content string `json:"content" gorm:"type:varchar(30);not null;unique"`
	Active  bool   `json:"active"`
	Created *time.Time
}

// TableName returns table name in database.
func (URL) TableName() string {
	return urlName
}

type serviceProvider struct{}

var (
	// Service handles operations on model URL.
	Service *serviceProvider
)

func init() {
	Service = &serviceProvider{}
}

// Create create a URL information.
func (sp *serviceProvider) Create(conn orm.Connection, content *string) error {
	now := time.Now()
	u := &URL{}
	value := &URL{
		Content: *content,
		Active:  true,
		Created: &now,
	}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(u).Create(value).Error
}

// Modify modify the specified URL content.
func (sp *serviceProvider) Modify(conn orm.Connection, uid *int16, content *string) error {
	u := &URL{}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(u).Where("id = ?", *uid).Update("content", *content).Error
}

// ModifyActive modify the specified URL status.
func (sp *serviceProvider) ModifyActive(conn orm.Connection, uid *int16, active *bool) error {
	u := &URL{}

	db := conn.(*gorm.DB).Exec("USE staff")

	return db.Model(u).Where("id = ?", *uid).Update("active", *active).Error
}

// List list all active URL.
func (sp *serviceProvider) List(conn orm.Connection) ([]URL, error) {
	u := &URL{}
	result := []URL{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(u).Where("active = true").Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// GetByID get the specified URL information.
func (sp *serviceProvider) GetByID(conn orm.Connection, uid *int32) (*URL, error) {
	u := &URL{}

	db := conn.(*gorm.DB).Exec("USE staff")
	err := db.Model(u).Where("id = ? AND active = true", *uid).First(u).Error

	if err != nil {
		return nil, err
	}

	return u, nil
}
