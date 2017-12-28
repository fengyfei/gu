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
 *     Initial: 2017/12/28        Feng Yifei
 */

package nsq

import (
	"fmt"
	"testing"
	"time"

	gnsq "github.com/nsqio/go-nsq"
)

func TestEvents(t *testing.T) {
	ev, err := NewEvents("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer ev.Close()

	handler := func(message *gnsq.Message) error {
		fmt.Printf("message: %+v", message)
		return nil
	}

	if err = ev.Publish([]byte("hello")); err != nil {
		t.Fatal("Publish error:", err)
	}

	if err = ev.CreateConsumer("channel", gnsq.HandlerFunc(handler)); err != nil {
		t.Fatal("Sub error:", err)
	}

	time.Sleep(2 * time.Second)
}
