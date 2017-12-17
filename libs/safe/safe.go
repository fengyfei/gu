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
 *     Initial: 2017/12/17        Feng Yifei
 */

package safe

import (
	"sync"
)

// Safe contains a thread-safe value
type Safe struct {
	value interface{}
	lock  sync.RWMutex
}

// New create a new Safe instance given a value
func New(value interface{}) *Safe {
	return &Safe{
		value: value,
		lock:  sync.RWMutex{},
	}
}

// Get returns the value
func (s *Safe) Get() interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.value
}

// Set sets a new value
func (s *Safe) Set(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.value = value
}
