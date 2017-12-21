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
	errInvalidSubscription = errors.New("Try to close a invalid subscription")
)

// Subscriber wraps a subscription to a nats server.
type Subscriber struct {
	conn    *Connection
	Subject string
	Group   string
	Handler nc.MsgHandler
	Sub     *nc.Subscription
}

// Unsubscribe the subject.
func (s *Subscriber) Unsubscribe() error {
	if s.Sub == nil || !s.Sub.IsValid() {
		return errInvalidSubscription
	}
	return s.Sub.Unsubscribe()
}

// MessageHandler is the messages handler for the subscription.
func (s *Subscriber) MessageHandler(msg *nc.Msg) {
	if s.Handler != nil {
		s.Handler(msg)
	}
}
