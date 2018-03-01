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
 *     Initial: 2018/03/01        Tong Yuehong
 */

package events

import (
	"fmt"
	"testing"
)

func TestSubscription_Subscribe(t *testing.T) {
	var (
		subscription = NewSubscription()
	)

	m := func(v interface{}) {}
	subscription.Subscribe(1, m)

	_, ok := subscription.subscribers[1]

	if !ok {
		t.Fatalf("error")
	}
}

func TestSubscription_Unsubscribe(t *testing.T) {
	var subscription = NewSubscription()
	m := func(v interface{}) {}
	subscriber := subscription.Subscribe(3, m)

	subscription.Unsubscribe(3, subscriber)

	s := subscription.subscribers[3]
	_, ok := s[subscriber]
	if ok {
		t.Fatalf("error")
	}
}

func TestSubscription_Notify(t *testing.T) {
	var (
		subscription = NewSubscription()
	)
	m := func(v interface{}) {
		fmt.Println(v)
	}

	subscription.Subscribe(3, m)

	subscription.Notify(3, "test")
}

func TestSubscription_NotifyAll(t *testing.T) {
	var (
		subscription = NewSubscription()
	)

	m := func(v interface{}) {
		if v == nil {
			fmt.Println("allone")
		} else {
			fmt.Println(v)
		}
	}

	n := func(v interface{}) {
		if v == nil {
			fmt.Println("alltwo")
		} else {
			fmt.Println(v)
		}
	}

	subscription.Subscribe(6, m)
	subscription.Subscribe(7, n)

	subscription.NotifyAll()
}
