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
	"testing"
)

var (
	store *Store
)

func TestNewStore(t *testing.T) {
	var (
		err error
	)

	store, err = NewStore("./user.db")

	if err != nil {
		t.Fatal(err)
	}
}

func TestRollback(t *testing.T) {
	var (
		err error
	)

	wr, err := store.Writer()
	if err != nil {
		t.Fatal(err)
	}

	err = wr.Put("user", []byte("rollback"), []byte("a"))
	if err != nil {
		wr.Rollback()
		t.Fatal(err)
	}
	wr.Commit()

	writer, err := store.Writer()
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Put("user", []byte("rollback"), []byte("test"))
	if err != nil {
		writer.Rollback()
		t.Fatal(err)
	}

	err = writer.Put("user", []byte(""), []byte("test"))
	if err != nil {
		writer.Rollback()
	} else {
		writer.Commit()
	}

	reader, err := store.Reader()
	if err != nil {
		t.Fatal(err)
	}

	m, err := reader.Switch("user").Get([]byte("rollback"))
	if err != nil || string(m) != "a" {
		t.Fatal(err)
	}

	err = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommit(t *testing.T) {
	var (
		err error
	)

	w, err := store.Writer()
	if err != nil {
		t.Fatal(err)
	}

	err = w.Put("user", []byte("commit"), []byte("a"))
	if err != nil {
		w.Rollback()
		t.Fatal(err)
	}

	err = w.Put("user", []byte("commit"), []byte("test"))
	if err != nil {
		w.Rollback()
		t.Fatal(err)
	} else {
		w.Commit()
	}

	reader, err := store.Reader()
	if err != nil {
		t.Fatal(err)
	}
	n, err := reader.Switch("user").Get([]byte("commit"))
	if err != nil || string(n) != "test" {
		t.Fatal(err)
	}

	err = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReader_Get(t *testing.T) {
	var (
		err error
	)

	wr, err := store.Writer()
	if err != nil {
		t.Fatal(err)
	}

	err = wr.Put("user", []byte("get"), []byte("get"))
	if err != nil {
		wr.Rollback()
		t.Fatal(err)
	}
	wr.Commit()


	reader, err := store.Reader()
	if err != nil {
		t.Fatal(err)
	}

	m, err := reader.Switch("user").Get([]byte("get"))
	if err != nil || string(m) != "get" {
		t.Fatal(err)
	}

	err = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
}
