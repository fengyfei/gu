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
 *     Initial: 2017/12/05        Jia Chenhui
 */

package config

import (
	"github.com/spf13/viper"
)

// rpcServerConfig represents the server config struct.
type rpcServerConfig struct {
	Address   string
	MysqlHost string
	MysqlPort string
	MysqlUser string
	MysqlPass string
	MysqlDb   string
}

var (
	Conf *rpcServerConfig
)

func init() {
	load()
}

// load read config file.
func load() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	Conf = &rpcServerConfig{
		Address:   viper.GetString("server.address"),
		MysqlHost: viper.GetString("mysql.host"),
		MysqlPort: viper.GetString("mysql.port"),
		MysqlUser: viper.GetString("mysql.user"),
		MysqlPass: viper.GetString("mysql.pass"),
		MysqlDb:   viper.GetString("mysql.db"),
	}
}
