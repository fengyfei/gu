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

// Bucket is a bit wise container used in bloom filters.
type Bucket struct {
	bits []byte
	m    uint
}

// NewBucket creates a empty bucket with at least m bits.
func NewBucket(m uint) *Bucket {
	return &Bucket{
		bits: make([]byte, (m+7)/8),
		m:    m,
	}
}

// Set the bit at position pos.
func (b *Bucket) Set(pos uint) {
	ind, off := b.position(pos)

	b.bits[ind] |= (1 << off)
}

// Get the bit at position.
func (b *Bucket) Get(pos uint) uint {
	ind, off := b.position(pos)

	return uint((b.bits[ind] & (1 << off)) >> off)
}

func (b *Bucket) position(pos uint) (uint, uint) {
	return pos / 8, pos % 8
}
