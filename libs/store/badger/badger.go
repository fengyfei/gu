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
 *     Initial: 2017/12/22        Yang Chenglong
 */

package badger

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

type BadgerDB struct {
	db *badger.DB
}

func NewBadgerDB(mode options.FileLoadingMode, dir string, compactFlag bool) (*BadgerDB, error) {
	opt := badger.DefaultOptions
	opt.TableLoadingMode = mode
	opt.Dir = dir
	opt.ValueDir = opt.Dir
	opt.DoNotCompact = compactFlag
	db, err := badger.Open(opt)
	return &BadgerDB{
		db: db,
	}, err
}

func (db *BadgerDB) Set(key, value []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

func (db *BadgerDB) Get(key []byte) ([]byte, error) {
	txn := db.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}
	v, err := item.Value()
	if err != nil {
		return nil, err
	}

	return v, err
}

func (db *BadgerDB) Delete(key []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
}
