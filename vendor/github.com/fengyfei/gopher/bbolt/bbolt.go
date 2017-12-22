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

	"github.com/coreos/bbolt"
)

type UserServiceProvider struct{}

var (
	UserService *UserServiceProvider = &UserServiceProvider{}
	UserDB      *bolt.DB
)

func init() {
	Open()
}

func Open() error {
	db, err := bolt.Open("user.db", 0666, nil)
	if err != nil {
		log.Printf("[open] init db error: %v", err)
		return err
	}
	UserDB = db

	tx, err := UserDB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.CreateBucketIfNotExists([]byte("user")); err != nil {
		return err
	}

	return tx.Commit()
}

func (usp *UserServiceProvider) Put(buc string, key string, payload []byte) error {
	tx, err := UserDB.Begin(true)
	if err != nil {
		log.Printf("[create] begin txn error: %v", err)
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(buc))

	if err := bucket.Put([]byte(key), payload); err != nil {
		log.Printf("put error: %v", err)
		return err
	}

	return tx.Commit()
}

func (usp *UserServiceProvider) Get(bucket string, key string) ([]byte, error) {
	tx, err := UserDB.Begin(false)
	if err != nil {
		log.Printf("[get] begin txn error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v := tx.Bucket([]byte(bucket)).Get([]byte(key))

	return v, nil
}
