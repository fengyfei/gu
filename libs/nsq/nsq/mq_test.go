/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/05/14        Tong Yuehong
 */

package nsq

import (
	"github.com/nsqio/go-nsq"
	"testing"
	"fmt"
)

func TestMq_NewConf(t *testing.T) {
	var mq Mq
	f := func(config *nsq.Config) {
		config.MaxInFlight = 1000
	}

	mq.NewConf(f)
	pro, err := mq.CreateProducer("127.0.0.1:4150")
	if err != nil {
		t.Fatal(err)
	}

	err = pro.Pub("test", "first message")
	if err != nil {
		t.Fatal(err)
	}

	h := func (message *nsq.Message) error {
		msg := string(message.Body)
		fmt.Println(msg)
		return nil
	}

	_, err = mq.CreateConsumer("test", "ch", "127.0.0.1:4150", nsq.HandlerFunc(h))
	if err != nil {
		t.Fatal(err)
	}
}
