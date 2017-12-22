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
 *     Initial: 2017/10/29        Feng Yifei
 */

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := stan.Connect(
		"test-cluster",
		"test-client",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		return err
	}
	defer logCloser(conn)

	wg := &sync.WaitGroup{}

	var sub stan.Subscription

	sub, err = conn.Subscribe("counter", func(msg *stan.Msg) {
		// Add jitter..
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Simulate failed acks to show redelivery behavior. 80% of the time
		// the ack will "succeed"
		if rand.Float32() > 0.2 {
			if err := msg.Ack(); err != nil {
				log.Printf("failed to ack")
				sub.Close()
				return
			}

			fmt.Printf("seq = %d [redelivered = %v, acked = true]\n", msg.Sequence, msg.Redelivered)
			wg.Done()
		} else {
			fmt.Printf("seq = %d [redelivered = %v, acked = false]\n", msg.Sequence, msg.Redelivered)
		}
	}, stan.SetManualAckMode(), stan.AckWait(time.Second))
	if err != nil {
		return err
	}
	defer logCloser(sub)

	// Publish up to 10.
	for i := 0; i < 10; i++ {
		wg.Add(1)

		err := conn.Publish("counter", nil)
		if err != nil {
			return err
		}
	}

	// Wait until all messages have been processed.
	wg.Wait()

	return nil
}
