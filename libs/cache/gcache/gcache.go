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

package gcache

import (
	"sync"

	"github.com/fengyfei/gu/libs/cache"
	"github.com/golang/groupcache/lru"
)

type gcache struct {
	mu  sync.RWMutex
	lru *lru.Cache
}

// NewCache returns a empty cache with max size maxEntries.
func NewCache(maxEntries int) cache.Cache {
	if maxEntries <= 0 {
		maxEntries = cache.DefaultMaxEntries
	}

	return &gcache{
		lru: lru.New(maxEntries),
	}
}

func (c *gcache) Put(key, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lru.Add(key, value)

	return nil
}

func (c *gcache) Get(key interface{}) (interface{}, error) {
	// Read operation may cause re-order in the underlying lru cache.
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, exists := c.lru.Get(key); exists {
		return val, nil
	}

	return nil, cache.ErrNotExists
}

func (c *gcache) Invalidate(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lru.Remove(key)
}

func (c *gcache) Size() int {
	c.mu.RLock()
	c.mu.RUnlock()

	return c.lru.Len()
}

func (c *gcache) Close() {
}
