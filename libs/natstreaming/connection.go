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

package natstreaming

import (
	"errors"

	stan "github.com/nats-io/go-nats-streaming"
)

var (
	errEmptySubject = errors.New("Could't subscribe to a empty subject")
	errEmptyGroup   = errors.New("Could't subscribe to a subject with a empty group name")
)

// Connection is a connection to nats streaming server.
type Connection struct {
	ClusterID string // Should be same with the nats streaming server
	ClientID  string
	URL       string
	conn      stan.Conn
}

// NewConnection creates a Connection to a nats streaming server.
func NewConnection(clusterID, clientID, url string) (*Connection, error) {
	c := &Connection{
		ClusterID: clusterID,
		ClientID:  clientID,
		URL:       url,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) connect() error {
	conn, err := stan.Connect(c.ClusterID, c.ClientID, stan.NatsURL(c.URL))
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// Close tear down the connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Subscribe on subject.
func (c *Connection) Subscribe(subject string, handler stan.MsgHandler, opt stan.SubscriptionOption) (*Subscriber, error) {
	if subject == "" {
		return nil, errEmptySubject
	}

	return c.subscribe(subject, "", handler, opt)
}

// QueueSubscribe creates a queue subscription on a nats streaming server on a subject.
func (c *Connection) QueueSubscribe(subject, group string, handler stan.MsgHandler, opt stan.SubscriptionOption) (*Subscriber, error) {
	if subject == "" {
		return nil, errEmptySubject
	}

	if group == "" {
		return nil, errEmptyGroup
	}

	return c.subscribe(subject, group, handler, opt)
}

func (c *Connection) subscribe(subject, group string, handler stan.MsgHandler, opt stan.SubscriptionOption) (*Subscriber, error) {
	s := &Subscriber{
		conn:    c,
		Subject: subject,
		Group:   group,
		Handler: handler,
	}

	if opt == nil {
		opt = subscribeDefaultOption
	}

	sub, err := c.conn.QueueSubscribe(subject, group, handler, opt)
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
