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
 *     Initial: 2018/01/19        Yang Chenglong
 */

package bloom

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewCount(t *testing.T) {
	for _, n := range []uint{1024, 2048, 4096, 8192} {
		falsePositive := 0
		c := NewCount(n, 0.03)

		for i := 1; i < int(n+1); i++ {
			key := bytes.Repeat([]byte("a"), i)

			if c.Contains(key) {
				falsePositive++
			}

			c.Add(key)

			if !c.Contains(key) {
				t.Fatalf("False Negative occured for key %s", string(key))
			}

			c.Delete(key)

			if c.Contains(key) {
				falsePositive++
			}
		}

		for i := int(n); i < 4*int(n)+1; i++ {
			key := bytes.Repeat([]byte("a"), i)

			if c.Contains(key) {
				falsePositive++
			}
		}

		fmt.Printf("n = %d, p := %f, falsePositive=%d\n", n, float64(falsePositive)/float64(1024), falsePositive)
	}
}
