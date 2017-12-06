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

package server

import (
	"net"
	"net/rpc"

	"github.com/fengyfei/gu/applications/auth/config"
	"github.com/fengyfei/gu/libs/logger"
)

var (
	// address represents the address of the RPC server.
	address = config.Conf.Address
)

func InitServer(db string) {
	rpc.Register(newAuthRPC(db))
	rpc.HandleHTTP()

	go rpcListen()
}

func rpcListen() {
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	defer func() {
		logger.Info("listen rpc: \"%s\" close", address)
		if err := l.Close(); err != nil {
			logger.Error(err)
		}
	}()

	rpc.Accept(l)
}
