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
 *     Modify : 2018/02/27        Tong Yuehong
 */

package events

import (
	"sync"
)

// Subscription is a collection of subscriptions for all events.
type Subscription struct {
	mu          sync.RWMutex
	subscribers map[EventType]map[Subscriber]EventHandler
}

// NewSubscription - create a subscription.
func NewSubscription() *Subscription {
	return &Subscription{
		subscribers: make(map[EventType]map[Subscriber]EventHandler),
	}
}

// Subscribe will create a subscriber on the given eventType using the specified Handler.
func (s *Subscription) Subscribe(etype EventType, h EventHandler) Subscriber {
	s.mu.Lock()
	defer s.mu.Unlock()

	subscriber := make(chan interface{})

	m, ok := s.subscribers[etype]
	if ok {
		m[subscriber] = h
		return nil
	}

	n := make(map[Subscriber]EventHandler)
	n[subscriber] = h
	s.subscribers[etype] = n

	return subscriber
}

// Unsubscribe deletes the subscriber from the subscribers which eventType is etype.
func (s *Subscription) Unsubscribe(etype EventType, subscriber Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subscribers[etype], subscriber)
	return nil
}

// Notify handles the handler from the subscribers which eventType is etype.
func (s *Subscription) Notify(etype EventType, val interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := s.subscribers[etype]
	for _, handler := range sub {
		handler(val)
	}

	return nil
}

// NofityAll handles all the handlers.
func (s *Subscription) NotifyAll() []error {
	var (
		errs []error
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subscribers{
		for _, handler := range sub {
			handler(nil)
		}
	}

	return errs
}
