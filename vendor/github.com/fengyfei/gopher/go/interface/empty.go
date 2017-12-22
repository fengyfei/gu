/*
 * MIT License
 *
 * Copyright (c) 2017 Feng Yifei.
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
 *     Initial: 2017/10/16        Feng Yifei
 */

package main

import (
	"fmt"
	"strconv"
)

// Empty - interface{}
type Empty interface{}

// Binary - uint64
type Binary uint64

func (b Binary) String() string {
	return strconv.FormatUint(b.get(), 2)
}

func (b Binary) get() uint64 {
	return uint64(b)
}

func testEmptyInterface(i Empty) bool {
	return i == nil
}

func main() {
	var bin *Binary

	fmt.Println("bin is empty:", testEmptyInterface(bin), bin)
}
