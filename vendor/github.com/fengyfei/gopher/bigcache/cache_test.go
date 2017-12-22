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
 *     Initial: 2017/09/16        Jia Chenhui
 */

package cache_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/fengyfei/gopher/bigcache"
)

func BenchmarkCacheServiceProvider_SetOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.CacheServer.SetOne(key(rand.Intn(b.N)), value())
	}
}

func BenchmarkCacheServiceProvider_GetOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.CacheServer.GetOne(strconv.Itoa(rand.Intn(b.N)))
	}
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value() []byte {
	return make([]byte, 100)
}
