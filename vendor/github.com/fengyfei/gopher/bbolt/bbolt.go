/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/09/14        Yang Chenglong
 */

package bbolt

import (
	"log"
	"os"

	"github.com/coreos/bbolt"
)

// Database represent a bolt instance.
type BboltDB struct {
	path    string
	db      *bolt.DB
	mode    os.FileMode
	options *bolt.Options
}

// NewDatabase creates a new bolt database.
func NewBboltDB(path string, mode os.FileMode, options *bolt.Options) (*BboltDB, error) {
	var (
		err error
	)

	db, err := bolt.Open(path, mode, options)
	if err != nil {
		log.Printf("[open] init db error: %v", err)
		return nil, err
	}

	return &BboltDB{
		db:      db,
		mode:    mode,
		options: options,
		path:    path,
	}, nil
}

// Put a key/vaue into the buck.
func (b *BboltDB) Put(buck string, key string, payload []byte) error {
	var (
		bucket *bolt.Bucket
		err    error
	)

	tx, err := b.db.Begin(true)
	if err != nil {
		log.Printf("[create] begin txn error: %v", err)
		return err
	}
	defer tx.Rollback()

	if bucket, err = tx.CreateBucketIfNotExists([]byte(buck)); err != nil {
		return err
	}

	if err := bucket.Put([]byte(key), payload); err != nil {
		log.Printf("put error: %v", err)
		return err
	}

	return tx.Commit()
}

// Path returns the location of the database.
func (b *BboltDB) Path() string {
	return b.db.Path()
}

// Close the internal database.
func (b *BboltDB) Close() error {
	return b.db.Close()
}

// Get a value via the key from the buck.
func (b *BboltDB) Get(buck string, key string) ([]byte, error) {
	var (
		err error
	)

	tx, err := b.db.Begin(false)
	if err != nil {
		log.Printf("[get] begin txn error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v := tx.Bucket([]byte(buck)).Get([]byte(key))

	return v, nil
}

// Delete a key/value from the buck..
func (b *BboltDB) Delete(buck string, key []byte) error {
	var (
		err error
	)

	tx, err := b.db.Begin(true)
	if err != nil {
		log.Printf("[get] begin txn error: %v", err)
		return err
	}
	defer tx.Rollback()

	return tx.Bucket([]byte(buck)).Delete([]byte(key))
}
