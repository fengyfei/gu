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
 *     Initial: 2018/02/27        Feng Yifei
 */

package events

import (
	"sync"
)

// Subscription -
type Subscription struct {
	mu          sync.RWMutex
	subscribers map[EventType]map[Subscriber]EventHandler
}

// NewSubscription -
func NewSubscription() *Subscription {
	return &Subscription{
		subscribers: make(map[EventType]map[Subscriber]EventHandler),
	}
}

// Subscribe -
func (s *Subscription) Subscribe(etype EventType, h EventHandler) Subscriber {
	return nil
}

// Unsubscribe -
func (s *Subscription) Unsubscribe(etype EventType, subscriber Subscriber) error {
	return nil
}

// Notify -
func (s *Subscription) Notify(etype EventType, val interface{}) error {
	return nil
}

// NofityAll -
func (s *Subscription) NofityAll() []error {
	var (
		errs []error
	)

	return errs
}
