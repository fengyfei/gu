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
 *     Initial: 2017/12/14        Jia Chenhui
 */

package mysql

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"

	"github.com/fengyfei/gu/libs/orm/mysql"
	"github.com/fengyfei/gu/models/staff"
)

var (
	// The only MySQL connection pool for this application.
	Pool *mysql.Pool
)

func init() {
	initConnection()
	initTable()
}

// initConnection  initializes the MySQL connection.
func initConnection() {
	user := beego.AppConfig.String("mysql::User")
	pass := beego.AppConfig.String("mysql::Pass")
	url := beego.AppConfig.String("mysql::Host")
	port := beego.AppConfig.String("mysql::Port")
	sqlName := beego.AppConfig.String("mysql::Db")

	// You need to create a database manually.
	// SQL: CREATE DATABASE staff CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")

	Pool = mysql.InitPool(dataSource)
}

// initTable create the MySQL table. All MySQL tables need to be created here.
func initTable() {
	conn, err := Pool.Get()
	if err != nil {
		panic(err)
	}
	defer Pool.Release(conn)

	db := conn.(*gorm.DB).Set("gorm:table_options", "ENGINE=InnoDB").Set("gorm:table_options", "CHARSET=utf8mb4")

	switch {
	case !db.HasTable(&staff.Staff{}):
		db.CreateTable(&staff.Staff{})
		staff.CreateAdminStaff(conn)
		fallthrough
	case !db.HasTable(&staff.Role{}):
		db.CreateTable(&staff.Role{})
		fallthrough
	case !db.HasTable(&staff.Relation{}):
		db.CreateTable(&staff.Relation{})
		fallthrough
	case !db.HasTable(&staff.Permission{}):
		db.CreateTable(&staff.Permission{})
		fallthrough
	default:
		fmt.Println("[MySQL] All tables have been created.")
	}
}
