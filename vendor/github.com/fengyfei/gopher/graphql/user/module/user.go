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
 *     Initial: 2017/09/27        Jia Chenhui
 */

package module

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gopher/graphql/user/mongo"
	"github.com/fengyfei/nuts/mgo/copy"
)

type UserServiceProvider struct {
}

var (
	UserService *UserServiceProvider
)

func init() {
	UserService = &UserServiceProvider{}
}

// user struct
type User struct {
	Login  string `json:"login"`
	Admin  string `json:"admin"`
	Active string `json:"active"`
}

// GetSingleInfo get single user information.
func (usp *UserServiceProvider) Get(login string) (User, error) {
	var (
		u User
	)

	query := bson.M{"login": login}
	err := copy.GetUniqueOne(mongo.MDInfo, query, &u)

	return u, err
}

// Create create single user.
func (usp *UserServiceProvider) Create(user *User) (bool, error) {
	u := User{
		Login:  user.Login,
		Admin:  user.Admin,
		Active: user.Active,
	}

	err := copy.Insert(mongo.MDInfo, &u)

	return err == nil, err
}
