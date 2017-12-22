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
 *     Initial: 2017/10/22        Jia Chenhui
 */

package pkg

import (
	"log"

	"github.com/nats-io/go-nats"
)

// NatsConn represents the nats connection.
type NatsConn struct {
	conn *nats.Conn
}

// NewConn create a nats connection.
// url can be empty string.
func NewConn(url string) (*NatsConn, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	c, err := nats.Connect(url)
	if err != nil {
		log.Printf("Connect to default url returned error: %v", err)
		return nil, err
	}

	nc := &NatsConn{
		conn: c,
	}

	log.Printf("Connection has been established.")
	return nc, nil
}
