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
 *     Initial: 2018/04/20        Li Zebang
 */

package increment

import (
	"fmt"
	"testing"
)

func TestBlotDB(t *testing.T) {
	db, err := open("test.db")
	if err != nil {
		panic(fmt.Errorf("open: %s", err.Error()))
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(fmt.Errorf("close: %s", err.Error()))
		}
	}()

	err = set(db, []byte("test"), []byte("test"), []byte("test"))
	if err != nil {
		panic(fmt.Errorf("set: %s", err.Error()))
	}

	value, err := get(db, []byte("test"), []byte("test"))
	if err != nil {
		panic(fmt.Errorf("get: %s", err.Error()))
	}
	if string(value) != "test" {
		panic(fmt.Errorf("get: incorrect value"))
	}
}
