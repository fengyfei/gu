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
 *     Initial: 2018/01/18        Yang Chenglong
 */

package bloom

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

// Count - a implement of count bloom filter.
type Count struct {
	Bucket *CountBucket
	h      hash.Hash64
	m      uint
	k      uint
}

// NewCount - make a count bloom.
func NewCount(n uint, p float64) *Count {
	if n == 0 {
		n = 1024
	}

	if p < 0.001 || p > 1.0 {
		p = 0.001
	}

	k := k(p)
	m := m(n, p)

	return &Count{
		Bucket: NewCountBucket(m * 4),
		h:      fnv.New64(),
		m:      m,
		k:      k,
	}
}

// Delete - Delete a element from bucket.
func (c *Count) Delete(key []byte) {
	l, h := c.hash(key)

	for i := uint(0); i < c.k; i++ {
		c.Bucket.Remove((uint(l) + uint(h)*i) % c.m)
	}

	return
}

// Contains - Filter interface implementation.
func (c *Count) Contains(key []byte) bool {
	l, h := c.hash(key)

	for i := uint(0); i < c.k; i++ {
		if c.Bucket.Get((uint(l)+uint(h)*i)%c.m) <= 0 {
			return false
		}
	}

	return true
}

// Add - Filter interface implementation.
func (c *Count) Add(key []byte) Filter {
	l, h := c.hash(key)

	for i := uint(0); i < c.k; i++ {
		c.Bucket.Set((uint(l) + uint(h)*i) % c.m)
	}

	return c
}

func (c *Count) hash(key []byte) (uint32, uint32) {
	c.h.Write(key)
	sum := c.h.Sum(nil)
	c.h.Reset()

	lower := binary.BigEndian.Uint32(sum[4:8])
	upper := binary.BigEndian.Uint32(sum[0:4])

	return lower, upper
}

// CountBucket is a bit wise container used in count bloom filters.
type CountBucket struct {
	bits []byte
	m    uint
}

// NewCountBucket creates a empty Bucket with at least m bits.
func NewCountBucket(m uint) *CountBucket {
	return &CountBucket{
		bits: make([]byte, (m+7)/8),
		m:    m,
	}
}

// Set the bit at position pos.
func (b *CountBucket) Set(pos uint) {
	ind, off := b.position(pos)
	if off < 4 {
		b.bits[ind] += 1
	}

	b.bits[ind] += 16
}

// Get the bit at position.
func (b *CountBucket) Get(pos uint) uint {
	ind, off := b.position(pos)
	if off < 4 {
		return uint(b.bits[ind] & 15)
	}
	return uint(b.bits[ind] >> off)
}

// Remove - decrease the count of one element.
func (b *CountBucket) Remove(pos uint) {
	ind, off := b.position(pos)

	if b.Get(pos) > 0 {
		if off < 4 {
			b.bits[ind] -= 1
		}

		b.bits[ind] -= 16
	}
}

func (b *CountBucket) position(pos uint) (uint, uint) {
	return pos * 4 / 8, pos * 4 % 8
}
