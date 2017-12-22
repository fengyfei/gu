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
 *     Initial: 2017/10/09        Jia Chenhui
 */

package account

import (
	"encoding/binary"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type AccountServiceProvider struct{}

var (
	AccountService *AccountServiceProvider = &AccountServiceProvider{}
	IDGen          *UniqueIDGenerator
)

func (asp *AccountServiceProvider) Create(id int, name string, info []byte) error {
	db, err := leveldb.OpenFile("account.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	idByte := intToByte(int(id))
	err = db.Put(idByte, info, nil)
	if err != nil {
		fmt.Printf("Put returned error: %v\n", err)
		return err
	}

	return nil
}

func (asp *AccountServiceProvider) Get(id int) ([]byte, error) {
	db, err := leveldb.OpenFile("account.db", nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	idByte := intToByte(id)

	data, err := db.Get(idByte, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func intToByte(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

type UniqueIDGenerator struct {
	id int64
}

func NewGenerator(start int64) *UniqueIDGenerator {
	return &UniqueIDGenerator{
		id: start,
	}
}

func (g *UniqueIDGenerator) Next() int64 {
	next := g.id
	g.id = g.id + 1

	return next
}
