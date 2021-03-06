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

func manualAckOption(opts *stan.SubscriptionOptions) error {
	opts.ManualAcks = true

	return nil
}

func main() {
	const (
		subject = "manual-ack"
	)

	var (
		subscriber *ns.Subscriber
		wg         = &sync.WaitGroup{}
	)

	messageHandler := func(msg *stan.Msg) {
		logger.Info(time.Now().UnixNano(), msg.Timestamp, msg.Sequence, msg.Subject, string(msg.Data))
		if err := msg.Ack(); err != nil {
			logger.Error("ACK error:", err)
		}
		wg.Done()
	}

	conn, err := ns.NewConnection("test-cluster", "tester", "nats://localhost:4222")
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()

	// subscribe
	subscriber, err = conn.Subscribe(subject, messageHandler, manualAckOption)
	if err != nil {
		logger.Error(err)
		return
	}
	defer subscriber.Unsubscribe()

	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			logger.Info("Send message sequence:", i, time.Now().UnixNano())
			if err := conn.Publish(subject, []byte(fmt.Sprintf("Message %d", i))); err != nil {
				logger.Error("Publish error:", err)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	logger.Debug("Closing")
}
