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

var chConnL = make(chan *net.Conn)
var chConnR = make(chan *net.Conn)

func main() {
	conn1 := flag.String("l", "", "Listen on addr")
	conn2 := flag.String("r", "", "Agency addr")
	flag.Parse()

	if *conn1 == "" || *conn2 == "" {
		log.Println("-l or -r can't be empty")
		return
	}
	go func() {
		dst := <-chConnL
		src := <-chConnR
		proxy.ProxyTCP(*dst, *src, 20*time.Second, 2*time.Second)
	}()
	go client(*conn1)
	go server(*conn2)
	time.Sleep(1e9 * 200)
}

func client(conn string) {
	listener, err := net.Listen("tcp", conn)
	if err != nil {
		log.Println("Listen Error: ", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept Error: ", err)
			return
		}

		chConnL <- &conn
	}
}

func server(connection string) {
	conn, err := net.Dial("tcp", connection)
	if err != nil {
		log.Println("Dial Error: ", err)
		return
	}

	chConnR <- &conn
}
