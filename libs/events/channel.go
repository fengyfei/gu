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
 *     Initial: 2018/02/24        Tong Yuehong
 */

package events

import (
	"sync"
)

// Channel represents a serial events queue.
type Channel struct {
	Ch     chan Event
	closed chan struct{}
	once   sync.Once
}

// NewChannel - new a channel.
func NewChannel(size int) *Channel {
	if size <= 0 {
		size = 10
	}

	return &Channel{
		Ch:     make(chan Event, size),
		closed: make(chan struct{}),
		once:   sync.Once{},
	}
}

// Done - a channel about getting the state of the Ch channel.
func (ch *Channel) Done() chan struct{} {
	return ch.closed
}

// Send - send a event.
func (ch *Channel) Send(ev Event) error {
	select {
	case ch.Ch <- ev:
		return nil
	case <-ch.closed:
		return ErrClosed
	}
}

// Close - close channel.
func (ch *Channel) Close() error {
	ch.once.Do(func() {
		close(ch.closed)
	})

	return nil
}
