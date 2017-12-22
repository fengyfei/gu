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
	"testing"
	"math/rand"
	"time"
	"strconv"
)

func TestDatabaseOpenClose(t *testing.T) {
	db, err := NewDatabase("tests/basic", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()
}

func TestDatabaseOperations(t *testing.T) {
	db, err := NewDatabase("tests/basic", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := []byte("level")
	if exists, _ := db.Has(key); exists {
		t.Fatal("Key exists")
	}

	if err = db.Put(key, []byte("first")); err != nil {
		t.Fatal(err)
	}

	if exists, _ := db.Has(key); !exists {
		t.Fatal("Key not exists")
	}

	if err = db.Put(key, []byte("second")); err != nil {
		t.Fatal(err)
	}

	if v, e := db.Get(key); e != nil {
		t.Fatal(e)
	} else if string(v) != "second" {
		t.Fatal("Invalid value")
	}

	if err = db.Delete(key); err != nil {
		t.Fatal(err)
	}

	if err = db.Delete(key); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkDatabase_Put(b *testing.B) {
	db, err := NewDatabase("tests/basic", 0, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		err := db.Put([]byte(strconv.Itoa(r.Intn(b.N))),  value())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDatabase_Get(b *testing.B) {
	db, err := NewDatabase("tests/basic", 0, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		db.Get([]byte(strconv.Itoa(r.Intn(b.N))))
	}
}

func value() []byte {
	return make([]byte, 100)
}

