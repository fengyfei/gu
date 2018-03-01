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
 *     Initial: 2018/02/26        Feng Yifei
 */

package boltdb

import (
	"errors"
	"bytes"

	bolt "github.com/coreos/bbolt"
)

var (
	errNoBucket = errors.New("boltdb: no bucket choosed")
)

// Reader handles read-only operations on a bolt DB.
type Reader struct {
	store  *Store
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

// Switch to a bucket.
func (r *Reader) Switch(bucket string) *Reader {
	r.bucket = r.tx.Bucket([]byte(bucket))

	return r
}

// Get read the given key from the bucket.
func (r *Reader) Get(key []byte) ([]byte, error) {
	var (
		rv []byte
	)

	if r.bucket == nil {
		return nil, errNoBucket
	}

	if bytes.Equal(key, nil){
		return nil, ErrInvalidKey
	}

	v := r.bucket.Get(key)
	if v != nil {
		// Value returned from Get() is only valid while the transaction is open, so a copy is required.
		rv = make([]byte, len(v))
		copy(rv, v)
	}

	return rv, nil
}

// Close the internal transaction.
func (r *Reader) Close() error {
	return r.tx.Rollback()
}

// ForEach
func (r *Reader) ForEach(fn func(k, v []byte) error) error {
	if r.bucket == nil {
		return errNoBucket
	}

	return r.bucket.ForEach(fn)
}
