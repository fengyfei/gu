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
 *     Initial: 2018/04/07        Tong Yuehong
 */

package workqueue

import (
	"sync"
	"testing"
	"time"
)

func TestParam(t *testing.T) {
	f := func(i int) {
	}

	err := Parallelize(-1, 3, f)
	if err == nil {
		t.Fatalf("error")
	}
	err = Parallelize(-1, -1, f)
	if err == nil {
		t.Fatalf("error")
	}

	err = Parallelize(1, -1, f)
	if err == nil {
		t.Fatalf("error")
	}
	err = Parallelize(1, 1, nil)
	if err == nil {
		t.Fatalf("error")
	}
}

func TestCompare(t *testing.T) {
	var (
		err    error
		num    = []int{1, 2, 3}
		numb   = []int{3, 2, 1}
		sum    = make([]int, 3)
		result = []int{4, 4, 4}
	)

	f := func(i int) {
		sum[i] = num[i] + numb[i]
	}

	err = Parallelize(3, 3, f)
	if err != nil {
		t.Fatal(err)
	}
	for i := range sum {
		if sum[i] != result[i] {
			t.Fatalf("error")
		}
	}

	err = Parallelize(2, 3, f)
	if err != nil {
		t.Fatal(err)
	}
	for i := range sum {
		if sum[i] != result[i] {
			t.Fatalf("error")
		}
	}

	err = Parallelize(5, 3, f)
	if err != nil {
		t.Fatal(err)
	}
	for i := range sum {
		if sum[i] != result[i] {
			t.Fatalf("error")
		}
	}
}

func TestMap(t *testing.T) {
	var (
		m = map[string]int{
			"first":  1,
			"second": 2,
			"third":  3,
		}
		err  error
		keys []string
		n    = make(map[string]int)
		mu   = sync.Mutex{}
	)

	for key := range m {
		keys = append(keys, key)
	}

	f := func(i int) {
		mu.Lock()
		defer mu.Unlock()

		n[keys[i]] = m[keys[i]]
	}

	err = Parallelize(3, 3, f)
	if err != nil {
		t.Fatal(err)
	}

	for k := range m {
		if m[k] != n[k] {
			t.Fatalf("error")
		}
	}
}

func TestHandleCrash(t *testing.T) {
	var (
		err    error
		num    = []int{1, 2, 3}
		numb   = []int{3, 2, 1}
		sum    = make([]int, 3)
		result = []int{4, 0, 4}
		e      interface{}
	)

	errHandler := func(err interface{}) {
		e = err
	}

	f := func(i int) {
		if i == 1 {
			panic("panic")
		}
		sum[i] = num[i] + numb[i]
	}

	err = Parallelize(3, 3, f, errHandler)

	if err != nil {
		t.Fatalf("error")
	}
	for i := range sum {
		if sum[i] != result[i] {
			t.Fatalf("error")
		}
	}

	time.Sleep(1e9)
}
