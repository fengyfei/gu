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
 *     Initial: 2017/11/03        ShiChao
 */

package main

import (
  "fmt"
  "github.com/astaxie/beego"
  "github.com/fengyfei/gu/applications/beego/shop/mysql"
  "github.com/jinzhu/gorm"
  "github.com/fengyfei/gu/models/shop/user"
  "github.com/fengyfei/gu/models/shop/category"
  "github.com/fengyfei/gu/models/shop/ware"
  "github.com/fengyfei/gu/models/shop/address"
)

func init() {
  initMysql()
  initTable()
}

func initMysql() {
  user := beego.AppConfig.String("mysqluser")
  pass := beego.AppConfig.String("mysqlpass")
  url := beego.AppConfig.String("mysqlurl")
  port := beego.AppConfig.String("mysqlport")
  sqlName := beego.AppConfig.String("mysqlname")

  dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + ":" + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")

  mysql.InitPool(dataSource)
}

// initTable create the MySQL table. All MySQL tables need to be created here.
func initTable() {
  conn, err := mysql.Pool.Get()
  if err != nil {
    panic(err)
  }
  defer mysql.Pool.Release(conn)

  db := conn.(*gorm.DB).Set("gorm:table_options", "ENGINE=InnoDB").Set("gorm:table_options", "CHARSET=utf8")

  if !conn.(*gorm.DB).HasTable("users") {
    db.CreateTable(
      &user.User{},
    )
  }
  if !conn.(*gorm.DB).HasTable("categories") {
    db.CreateTable(
      &category.Category{},
    )
  }
  if !conn.(*gorm.DB).HasTable("wares") {
    db.CreateTable(
      &ware.Ware{},
    )
  }
  if !conn.(*gorm.DB).HasTable("addresses") {
    db.CreateTable(
      &address.Address{},
    )
  }

}
