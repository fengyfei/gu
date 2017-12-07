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
 *     Initial: 2017/12/07        Feng Yifei
 */

package godash

import (
	"testing"
)

// TestEachArray -
func TestEachArray(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	Each(arr, func(i, n int) {
		if n != arr[i] {
			t.Error("Each not working for array")
		}
	})
}

// TestEachMap -
func TestEachMap(t *testing.T) {
	data := map[string]string{
		"Feng": "Sinian",
		"Yang": "Chenglong",
		"Shi":  "Chao",
		"Wang": "Riyu",
	}

	Each(data, func(k, v string) {
		if v != data[k] {
			t.Error("Each not working for map")
		}
	})
}
