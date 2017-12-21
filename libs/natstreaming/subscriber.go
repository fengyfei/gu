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
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// Subscriber wraps a subscription to a channel.
type Subscriber struct {
	conn    *Connection // Original connection
	Subject string
	Group   string
	Handler stan.MsgHandler
	Sub     stan.Subscription
}

// Unsubscribe the subject.
func (sub *Subscriber) Unsubscribe() error {
	return sub.Sub.Unsubscribe()
}

// MessageHandler is the messages handler for the subscription.
func (sub *Subscriber) MessageHandler(msg *stan.Msg) {
	if sub.Handler != nil {
		sub.Handler(msg)
	}
}

func subscribeDefaultOption(opts *stan.SubscriptionOptions) error {
	opts.DurableName = ""
	opts.MaxInflight = 1
	opts.AckWait = 10 * time.Second
	opts.StartAt = 0 // StartPosition_NewOnly
	opts.ManualAcks = false

	return nil
}
