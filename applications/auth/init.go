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
 *     Initial: 2017/12/05        Jia Chenhui
 */

package main

import (
	"fmt"

	"github.com/fengyfei/gu/applications/auth/config"
	"github.com/fengyfei/gu/applications/auth/mysql"
	"github.com/fengyfei/gu/applications/auth/server"
)

func init() {
	initMysql()
}

// initMysql initializes the MySQL connection.
func initMysql() {
	user := config.ServerConfig.MysqlUser
	pass := config.ServerConfig.MysqlPass
	url := config.ServerConfig.MysqlHost
	port := config.ServerConfig.MysqlPort
	sqlName := config.ServerConfig.MysqlDb

	dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")

	mysql.InitPool(dataSource)
}

func run() {
	go server.InitServer()

	fmt.Println("RPC server started on:", server.RPCAddr)
	select {}
}
