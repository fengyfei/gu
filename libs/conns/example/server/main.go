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
	"net"

	"github.com/fengyfei/gu/libs/conns/tcp"
	"github.com/fengyfei/gu/libs/logger"
)

var (
	session = tcp.NewSession("sessionOne")
)

func main() {
	server, err := tcp.NewServer(":18600")
	if err != nil {
		panic(err)
	}
	defer server.Listener().Close()

	logger.Debug("Listener on:", server.Listener().Addr())

	for {
		conn, err := server.Listener().Accept()
		if err != nil {
			continue
		}

		logger.Debug("A client connected:", conn.RemoteAddr().String())

		// add conn to server and session
		server.Add(conn.RemoteAddr().String())
		go tcpPipe(conn)
	}
}

// tcpPipe read from conn and broadcast.
func tcpPipe(conn net.Conn) {
	ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	session.Add(ipStr, conn)
	defer func() {
		logger.Debug("Disconnected:", ipStr)
		conn.Close()
	}()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		logger.Debug("Read message from "+ipStr+":", string(msg))
		logger.Debug("session number is:", session.Counts())

		if session.Counts() >= int32(2) {
			err := session.Broadcast([]byte(msg))
			if err != nil {
				logger.Debug("Broadcast error:", err)
			}
			logger.Debug("Broadcast success, connections number is:", session.Counts())
		}
	}
}
