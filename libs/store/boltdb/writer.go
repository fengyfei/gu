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
 *     Modify : 2018/02/26        Tong Yuehong
 */

package boltdb

import (
	"github.com/coreos/bbolt"
	"bytes"
)

// Writer handles write operations on a bolt DB.
type Writer struct {
	store *Store
	tx    *bolt.Tx
}

// Bucket create a bucket.
func (w *Writer) Bucket(bucket string) (*bolt.Bucket, error) {
	b, err := w.tx.CreateBucketIfNotExists([]byte(bucket))
	if err != nil {
		return nil, err
	}

	return b, nil
}

//Put put key-value into db, can't be used in multiple goroutines.
func (w *Writer) Put(bucket string, key []byte, value []byte) error {
	b, err := w.Bucket(bucket)
	if err != nil {
		return err
	}

	if bytes.Equal(key, nil){
		return ErrInvalidKey
	}

	err = b.Put(key, value)
	if err != nil {
		return err
	}

	return nil
}

// Commit closes the internal transaction.
func (w *Writer) Commit() error {
	return w.tx.Commit()
}

// Rollback closes the transaction and ignores all previous updates.
func (w *Writer) Rollback() error {
	return w.tx.Rollback()
}
