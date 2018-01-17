/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd.
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
 *     Initial: 2018/01/16        Jia Chenhui
 */

package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/fengyfei/gu/libs/conns/tcp"
	"github.com/fengyfei/gu/libs/logger"
)

func main() {
	var (
		err  error
		conn *tcp.Conn
	)

	conn, err = tcp.NewConn(":18600")
	if err != nil {
		panic(err)
	}
	defer conn.Conn().Close()
	logger.Debug("Client connected to:", conn.Conn().RemoteAddr().String())

	go onMsgReceived(conn.Conn())

	for {
		var msg string
		fmt.Scanln(&msg)
		if msg == "quit" {
			break
		}
		b := []byte(msg + "\n")
		conn.Conn().Write(b)
	}
}

func onMsgReceived(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		logger.Debug("Received msg:", msg)
		if err != nil {
			break
		}
	}
}
