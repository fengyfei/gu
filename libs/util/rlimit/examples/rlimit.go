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
 *     Initial: 2017/12/26        Yang Chenglong
 */

package main

import (
	"fmt"
	"os"

	"github.com/fengyfei/gu/libs/util/rlimit"
)

func main() {
	cur, max, err := rlimit.GetMaxOpens()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cur = %d, Max = %d \n", cur, max)

	for i := uint64(0); i < cur; i++ {
		fd, err := os.Open("../rlimit.go")
		if err != nil {
			fmt.Printf("Open file failed %+v, cur = %d\n", err, i+1)
			break
		}
		defer fd.Close()
	}

	if err = rlimit.SetNumFiles(1000000); err != nil {
		fmt.Printf("Set max opened files %+v", err)
		return
	}

	cur, max, err = rlimit.GetMaxOpens()
	if err != nil {
		fmt.Printf("Get Max open failed %+v", err)
	}
	fmt.Printf("After set, Cur = %d, Max = %d\n", cur, max)

	for i := uint64(0); i < 100000; i++ {
		fd, err := os.Open("../rlimit.go")
		if err != nil {
			fmt.Printf("Open file failed %+v, cur = %d\n", err, i+1)
			return
		}
		defer fd.Close()
	}

	fmt.Println("Succeed")
}
