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
 *     Initial: 2018/02/26        Shi Ruitao
 */

package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/fengyfei/gu/libs/network/proxy"
)

func main() {
	localAddr := flag.String("l", "", "local address")
	remoteAddr := flag.String("r", "", "remote address")
	flag.Parse()

	log.Printf("local address: %s, remote address: %s\n", *localAddr, *remoteAddr)

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatalln("local listen error:", err)
	}

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Fatalln("local accept error:", err)
		}

		remoteConn, err := net.Dial("tcp", *remoteAddr)
		if err != nil {
			log.Fatalln("remote dial error:", err)
		}

		go proxy.ProxyTCP(localConn, remoteConn, 10*time.Second, 10*time.Second)
	}
}
