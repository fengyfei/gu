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
 *     Initial: 2018/04/21        Li Zebang
 */

package bolt

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

var (
	// ErrEmptyPath would be returned when path is "".
	ErrEmptyPath = errors.New("boltdb path required")
	// ErrEmptyBucket would be returned when bucket is "".
	ErrEmptyBucket = errors.New("bucket name required")
)

// DB equals to bolt.DB
type DB struct {
	db *bolt.DB
}

// Open creates and opens a database at the given path. If the file does not exist
// then it will be created automatically. If path is "", it will return ErrEmptyPath.
func Open(path string) (*DB, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Close releases all database resources. All transactions must be closed before closing
// the database.
func (db *DB) Close() error {
	return db.db.Close()
}

// Set sets the value for a key in the bucket.If the key exist then its previous value
// will be overwritten. Supplied value must remain valid for the life of the transaction.
// If bucket is "", it will return ErrEmptyBucket. Returns an error if the bucket was
// created from a read-only transaction, if the key is blank, if the key is too large,
// or if the value is too large.
func (db *DB) Set(bucket, key, value []byte) error {
	if len(bucket) == 0 {
		return ErrEmptyBucket
	}

	return db.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

// Get retrieves the value for a key in the bucket. Returns a nil value if the key does not
// exist or if the key is a nested bucket. The returned value is only valid for the life of
// the transaction. If bucket is "", it will return ErrEmptyBucket.
func (db *DB) Get(bucket, key []byte) (value []byte, err error) {
	if len(bucket) == 0 {
		return nil, ErrEmptyBucket
	}

	return value, db.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		value = b.Get(key)
		return nil
	})
}
