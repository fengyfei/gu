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
 *     Initial: 2018/01/21        Chen Yanchen
 */

package user

import (
	"time"
)

type UserServiceProvider struct{}

var UserServer *UserServiceProvider

type (
	User struct {
		Id       uint64     `orm:"column(id);pk"`
		UserName string     `orm:"column(name)";unique;	json:"user_name"`
		Password string     `orm:"column(password)" 	json:"password"`
		Phone    string     `orm:"column(phone)" 		json:"phone"`
		Created  *time.Time `orm:"column(created)"`
		Status   bool       `orm:"column(status)"`
		// Validate uint       `orm:"column(validate)" 	json:"validate"`
		ThemesNum  uint64 `orm:"column(themes)"`
		ArticleNum uint64 `orm:"column(articles)"`
		LastLogin  *time.Time
	}

	UserInfo struct {
		Id     uint64 `orm:"column(id);pk"`
		UserId uint64
	}

	R struct {
		Id      uint64
		UserId  uint64
		ThemeId uint64
		Type    uint8
		Status  bool
	}
)
