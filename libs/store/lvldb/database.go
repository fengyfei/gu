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
 *     Initial: 2017/12/22        Feng Yifei
 */

package lvldb

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	leveldberr "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	errInvalidPath = errors.New("path is invalid")
)

// Database is a goleveldb instance wrapper.
type Database struct {
	path string
	db   *leveldb.DB
}

// NewDatabase creates a new leveldb database.
func NewDatabase(path string, cache, handles int) (*Database, error) {
	if path == "" {
		return nil, errInvalidPath
	}

	if cache < 16 {
		cache = 16
	}

	if handles < 16 {
		handles = 16
	}

	db, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache,
		WriteBuffer:            cache / 4 * opt.MiB,
		Filter:                 filter.NewBloomFilter(10),
	})

	if _, corrupted := err.(*leveldberr.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}

	if err != nil {
		return nil, err
	}

	return &Database{
		path: path,
		db:   db,
	}, nil
}

// Path returns the location of the database.
func (db *Database) Path() string {
	return db.path
}

// Close the internal database.
func (db *Database) Close() error {
	return db.db.Close()
}
