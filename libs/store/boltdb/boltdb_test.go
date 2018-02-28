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
 *     Initial: 2018/02/26        Tong Yuehong
 */

package boltdb

import (
	"fmt"
	"testing"
)

var (
	store *Store
)

func TestNewStore(t *testing.T) {
	var err error
	store, err = NewStore("./user.db")

	if err != nil {
		t.Fatal(err)
	}
}

func TestWriter_Put(t *testing.T) {
	var err error
	writer, err := store.Writer()
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Put("user", []byte(fmt.Sprintf("%d", 10)), []byte("test"))
	if err != nil {
		t.Fatal(err)
		writer.Rollback()
	}

	err = writer.Put("user", []byte(""), []byte("test"))
	if err == nil {
		t.Fatal(err)
	}

	err = writer.Put("test", []byte("test"), []byte("test"))
	if err != nil {
		t.Fatal(err)
		writer.Rollback()
	} else {
		writer.Commit()
	}
}

func TestReader_Get(t *testing.T) {
	var err error
	reader, err := store.Reader()
	if err != nil {
		t.Fatal(err)
	}

	v, err := reader.Switch("user").Get([]byte(fmt.Sprintf("%d", 10)))
	if err != nil || string(v) != "test" {
		t.Fatal(err)
	}

	a, err := reader.Switch("test").Get([]byte("test"))
	if err != nil || string(a) != "test" {
		t.Fatal(err)
	}

	err = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
}
