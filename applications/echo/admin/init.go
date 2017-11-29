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
 *     Initial: 2017/11/01        Jia Chenhui
 */

package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/fengyfei/gu/applications/echo/admin/mysql"
	"github.com/fengyfei/gu/applications/echo/admin/routers"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/models/staff"
)

var (
	server *echo.Echo
)

func init() {
	readConfiguration()
	initMysql()
	initTable()
}

// initMysql  initializes the MySQL connection.
func initMysql() {
	user := configuration.mysqlUser
	pass := configuration.mysqlPass
	url := configuration.mysqlHost
	port := configuration.mysqlPort
	sqlName := configuration.mysqlDb

	dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")

	mysql.InitPool(dataSource)
}

// initTable create the MySQL table. All MySQL tables need to be created here.
func initTable() {
	conn, err := mysql.Pool.Get()
	if err != nil {
		panic(err)
	}
	defer mysql.Pool.Release(conn)

	db := conn.(*gorm.DB).Set("gorm:table_options", "ENGINE=InnoDB")

	switch {
	case !db.HasTable(&staff.Staff{}):
		db.CreateTable(&staff.Staff{})
		staff.CreateAdminStaff(conn)
		fallthrough
	case !db.HasTable(&staff.Role{}):
		db.CreateTable(&staff.Role{})
		staff.CreateAdminRole(conn)
		fallthrough
	case !db.HasTable(&staff.Relation{}):
		db.CreateTable(&staff.Relation{})
		staff.CreateAdminRelation(conn)
		fallthrough
	case !db.HasTable(&staff.Permission{}):
		db.CreateTable(&staff.Permission{})
		staff.CreateAdminPermission(conn)
		fallthrough
	default:
		fmt.Println("[MySQL] All tables have been created.")
	}
}

// startEchoServer starts an HTTP server.
func startEchoServer() {
	server = echo.New()

	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: configuration.corsHosts,
		AllowMethods: []string{echo.GET, echo.POST},
	}))
	server.Use(core.CustomJWT(configuration.tokenKey))
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${uri} ${method} ${status} ${remote_ip} ${latency_human} ${bytes_in} ${bytes_out}\n",
	}))

	server.HTTPErrorHandler = core.EchoRestfulErrorHandler
	server.Validator = core.NewEchoValidator()

	if configuration.isDebug {
		log.SetLevel(log.DEBUG)
	} else {
		log.SetLevel(log.INFO)
	}

	routers.InitRouter(server, configuration.tokenKey)

	server.Start(configuration.address)
}
