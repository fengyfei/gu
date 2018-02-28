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
)

// Writer handles write operations on a bolt DB.
type Writer struct {
	store *Store
}

type Param struct {
	Bucket string
	Key    [][]byte
	Value  [][]byte
}

//CommPut put key-value into db, can't be used in multiple goroutines.
func (wr *Writer) CommPut(p ...Param) error {
	return wr.store.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range p {
			if len(bucket.Key) != len(bucket.Value) {
				return ErrInvalidParam
			}

			b, err := tx.CreateBucketIfNotExists([]byte(bucket.Bucket))
			if err != nil {
				return err
			}

			if len(bucket.Key) == 0 {
				continue
			}

			for i := 0; i < len(bucket.Key); i++ {
				if err := b.Put(bucket.Key[i], bucket.Value[i]); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

//MultiplePut only useful when there are multiple goroutines putting Key-Value to db.
func (wr *Writer) MultiplePut(p ...Param) error {
		return wr.store.db.Batch(func(tx *bolt.Tx) error {
			for _, bucket := range p {
				if len(bucket.Key) != len(bucket.Value) {
					return ErrInvalidParam
				}

				b, err := tx.CreateBucketIfNotExists([]byte(bucket.Bucket))
				if err != nil {
					return err
				}

				if len(bucket.Key) == 0 {
					continue
				}

				for i := 0; i < len(bucket.Key); i++ {
					if err := b.Put(bucket.Key[i], bucket.Value[i]); err != nil {
						return err
					}
				}
			}
			return nil
		})
}
