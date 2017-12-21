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

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fengyfei/gu/libs/logger"
	ns "github.com/fengyfei/gu/libs/natstreaming"
	stan "github.com/nats-io/go-nats-streaming"
)

type messageHandler struct {
	id int
	wg *sync.WaitGroup
}

func (m *messageHandler) Handler(msg *stan.Msg) {
	logger.Info(m.id, time.Now().UnixNano(), msg.Timestamp, msg.Sequence, msg.Subject, string(msg.Data))
	m.wg.Done()
}

func main() {
	const (
		subject = "group"
		group   = "worker"
	)

	var (
		subscriber *ns.Subscriber
		wg         = &sync.WaitGroup{}
	)

	conn, err := ns.NewConnection("test-cluster", "tester", "nats://localhost:4222")
	if err != nil {
		logger.Error(err)
		return
	}

	defer conn.Close()

	for i := 0; i < 4; i++ {
		m := &messageHandler{
			id: i,
			wg: wg,
		}
		subscriber, err = conn.QueueSubscribe(subject, group, m.Handler, nil)
		if err != nil {
			logger.Error(err)
			return
		}
		defer subscriber.Unsubscribe()
	}

	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			if err := conn.Publish(subject, []byte(fmt.Sprintf("Message %d", i))); err != nil {
				logger.Error("Publish error:", err)
			}
		}
		wg.Done()
	}()

	wg.Wait()
	logger.Debug("Closing")
}
