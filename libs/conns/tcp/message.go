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
 *     Initial: 2018/01/08        Feng Yifei
 */

package tcp

import (
	"net"
)

// Message is a general tcp message format.
type Message struct {
	len     int32
	id      int32
	payload []byte
}

// NewMessage create a Message with capacity.
func NewMessage(cap int) *Message {
	return &Message{
		len:     0,
		id:      0,
		payload: make([]byte, cap),
	}
}

// Reset clean the Message.
func (m *Message) Reset() {
	m.len = 0
	m.id = 0
}

// Read read from conn.
func (m *Message) Read(conn net.Conn) error {
	size, err := conn.Read(m.payload)
	if err != nil {
		return err
	}

	m.len = int32(size)

	return nil
}

// Write write to conn.
func (m *Message) Write(conn net.Conn) error {
	_, err := conn.Write(m.payload[:m.len])
	if err != nil {
		return err
	}

	return nil
}
