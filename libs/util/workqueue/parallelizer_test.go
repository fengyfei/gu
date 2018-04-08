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
	"fmt"
	"testing"
)

func TestParallelize(t *testing.T) {
	f := func(i int) {
		fmt.Println("a")
	}

	Parallelize(-1, 3, f)

	Parallelize(-1, -1, f)

	Parallelize(3, -1, f)
}

func TestParallelize2(t *testing.T) {
	var num = []int{1, 2, 3}
	var numb = []int{3, 2, 1}
	var sum = make([]int, 3)

	f := func(i int) {
		sum[i] = num[i] + numb[i]
	}

	Parallelize(3, 3, f)
	fmt.Println(sum)

	Parallelize(2, 3, f)
	fmt.Println(sum)

	Parallelize(5, 3, f)
	fmt.Println(sum)
}

func TestParallelize3(t *testing.T) {
	var num = []int{1, 2, 3}
	var numb = []int{3, 2, 1}
	var sum = make([]int, 3)

	f := func(i int) {
		if i == 2 {
			panic("panic")
		}
		sum[i] = num[i] + numb[i]
	}

	Parallelize(3, 3, f)
	fmt.Println(sum)
}
