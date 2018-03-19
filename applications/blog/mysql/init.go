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
 *     Initial: 2018/02/06        Tong Yuehong
 */

package mysql

import (
	"fmt"

	"github.com/fengyfei/gu/applications/blog/conf"
	"github.com/fengyfei/gu/libs/orm/mysql"
)

var Pool *mysql.Pool

func InitMysql() {
	user := conf.Config.MysqlUser
	pass := conf.Config.MysqlPass
	url := conf.Config.MysqlHost
	port := conf.Config.MysqlPort
	sqlName := conf.Config.MysqlDb

	dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")
	go InitPool(dataSource)
}

// InitPool initialize the connection pool.
func InitPool(db string) {
	Pool = mysql.NewPool(db, 20)

	if Pool == nil {
		panic("MySQL DB connection error.")
	}
}
