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
 *     Initial: 2017/12/23        Feng Yifei
 */

package bloom

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

// Classic implements a classic Bloom Filter.
type Classic struct {
	bucket *Bucket
	h      hash.Hash64
	m      uint
	k      uint
}

// NewClassic returns a classic bloom filter.
func NewClassic(n uint, p float64) *Classic {
	if n == 0 {
		n = 1024
	}

	if p < 0.01 || p > 1.0 {
		p = 0.01
	}

	k := k(p)
	m := m(n, p)

	return &Classic{
		bucket: NewBucket(m),
		h:      fnv.New64(),
		m:      m,
		k:      k,
	}
}

// Contains - Filter interface implementation.
func (c *Classic) Contains(key []byte) bool {
	l, h := c.hash(key)

	for i := uint(0); i < c.k; i++ {
		if c.bucket.Get((uint(l)+uint(h)*i)%c.m) == 0 {
			return false
		}
	}

	return true
}

// Add - Filter interface implementation.
func (c *Classic) Add(key []byte) Filter {
	l, h := c.hash(key)

	for i := uint(0); i < c.k; i++ {
		c.bucket.Set((uint(l)+uint(h)*i)%c.m)
	}

	return c
}

func (c *Classic) hash(key []byte) (uint32, uint32) {
	c.h.Write(key)
	sum := c.h.Sum(nil)
	c.h.Reset()

	lower := binary.BigEndian.Uint32(sum[4:8])
	upper := binary.BigEndian.Uint32(sum[0:4])

	return lower, upper
}
