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
 *     Modify : 2017/12/22        Yang Chenglong
 */

package bbolt

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var db *BboltDB

func init() {
	db,_ = NewBboltDB("user.db",0666,nil)
}

func BenchmarkUserServiceProvider_Put(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("test%s", strconv.Itoa(r.Intn(b.N)))
		err := db.Put("user", name, value())
		if err != nil {
			log.Printf("create testing: %v", err)
		}
	}
}

func BenchmarkUserServiceProvider_Get(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("test%s", strconv.Itoa(r.Intn(b.N)))
		db.Get("user", name)
	}
}

func value() []byte {
	return make([]byte, 100)
}
