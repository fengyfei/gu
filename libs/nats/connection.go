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
 *     Initial: 2017/12/21        Feng Yifei
 */

package nats

import (
	"errors"

	nc "github.com/nats-io/go-nats"
)

var (
	errEmptySubject = errors.New("Could't subscribe to a empty subject")
	errEmptyGroup   = errors.New("Could't subscribe to a subject with a empty group name")
)

// Connection is a connection to nats server.
type Connection struct {
	URL  string
	conn *nc.Conn
}

// NewConnection creates a Connection to a nats server.
func NewConnection(url string, opt ...nc.Option) (*Connection, error) {
	c := &Connection{
		URL: url,
	}

	if err := c.connect(opt...); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) connect(opt ...nc.Option) error {
	conn, err := nc.Connect(c.URL, opt...)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// Subscribe on subject.
func (c *Connection) Subscribe(subject string, handler nc.MsgHandler) (*Subscriber, error) {
	if subject == "" {
		return nil, errEmptySubject
	}

	return c.subscribe(subject, "", handler)
}

func (c *Connection) subscribe(subject, group string, handler nc.MsgHandler) (*Subscriber, error) {
	var (
		sub *nc.Subscription
		err error
	)
	s := &Subscriber{
		conn:    c,
		Subject: subject,
		Group:   group,
	}

	if group != "" {
		sub, err = c.conn.QueueSubscribe(subject, group, handler)
	} else {
		sub, err = c.conn.Subscribe(subject, handler)
	}
	if err != nil {
		return nil, err
	}
	s.Sub = sub

	return s, nil
}

// Publish a message on a subject.
func (c *Connection) Publish(subject string, message []byte) error {
	return c.conn.Publish(subject, message)
}

// Close tears down the connection to a nats server.
func (c *Connection) Close() {
	c.conn.Close()
}
