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
	"errors"
	"fmt"

	"github.com/nats-io/go-nats"
)

// Sub subscribe to the connection based on the provided subject.
// nc can be nil.
// msgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
func AsyncSub(nc *NatsConn, subj string, msgHandler nats.MsgHandler) (*nats.Subscription, error) {
	var (
		err error
	)

	if nc == nil {
		nc, err = NewConn("")
		if err != nil {
			return nil, err
		}
	}

	defer nc.conn.Flush()

	if subj == "" || msgHandler == nil {
		err := errors.New("empty subject or msgHandler")
		return nil, err
	}

	s, err := nc.conn.Subscribe(subj, msgHandler)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func DisplayMsg(msg *nats.Msg) {
	fmt.Printf("Received on [%s]: '%s'\n", msg.Subject, string(msg.Data))
}
