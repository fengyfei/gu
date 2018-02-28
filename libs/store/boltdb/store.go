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
	"time"

	bolt "github.com/coreos/bbolt"
)

var (
	ErrEmptyPath = errors.New("boltdb: empty path")
	ErrInvalidParam = errors.New("Invalid Param")
)

// Store represents a general bolt db.
type Store struct {
	path string
	db   *bolt.DB
}

// NewStore create/open a bolt database.
func NewStore(path string) (*Store, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return nil, err
	}

	return &Store{
		path: path,
		db:   db,
	}, nil
}

// Reader returns a reader.
func (st *Store) Reader() (*Reader, error) {
	tx, err := st.db.Begin(false)
	if err != nil {
		return nil, err
	}

	return &Reader{
		store: st,
		tx:    tx,
	}, nil
}

// Writer returns a writer.
func (st *Store) Writer() *Writer {
	return &Writer{
		store: st,
	}
}
