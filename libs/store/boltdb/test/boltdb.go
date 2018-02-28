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
 *     Initial: 2018/02/26        Tong Yuehong
 */

package main

import (
	"fmt"
	"sync"
	"time"

	bolt "github.com/fengyfei/gu/libs/store/boltdb"
)

func main() {
	kv, err := bolt.NewStore("./user.db")

	if err != nil {
		fmt.Println(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(i int) {
			err := kv.Writer().MultiplePut(bolt.Param{
				Bucket:"user",
				Key: [][]byte{[]byte("one"), []byte("two")},
				Value: [][]byte{[]byte("one"), []byte("two")},
			})

			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(i)
		time.Sleep(500 * time.Millisecond)
	}

	reader, err := kv.Reader()
	if err != nil {
		fmt.Println(err)
	}

	v, err := reader.Switch("user").Get([]byte("one"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v", string(v))

	reader.ForEach(func(k, v []byte) error {
		fmt.Printf("key=%s, value=%s\n", k, v)
		return nil
	})

	wg.Wait()
}
