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
 *     Initial: 2017/12/23        Feng Yifei
 */

package bloom

import (
	"fmt"
	"testing"
)

func TestOptimalK(t *testing.T) {
	for _, p := range []float64{0.2, 0.15, 0.1, 0.08, 0.05, 0.01, 0.005} {
		fmt.Printf("k(%f)=%d\n", p, k(p))
	}
}

func TestOptimalM(t *testing.T) {
	for _, n := range []uint{100, 1024, 10000, 100000} {
		for _, p := range []float64{0.2, 0.15, 0.1, 0.08, 0.05, 0.01, 0.005} {
			fmt.Printf("m(%d,%f)=%d\n", n, p, m(n, p))
		}
	}
}
