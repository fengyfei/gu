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
 *     Initial: 2018/02/25        Tong Yuehong
 */

package main

import (
	"fmt"
	"github.com/fengyfei/gu/libs/events"
	"sync"
	"time"
)

func main() {
	a := events.NewChannel(9)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		t1 := time.NewTimer(time.Second * 5)
		flag := false
		for {
			select {
			case c := <-a.Ch:
				fmt.Println(c)
				if !flag {
					t1.Reset(time.Second * 5)
				}
			case <-a.Done():
				flag = true
				break
			case <-t1.C:
				wg.Done()
				fmt.Println("finish")
				return
			}
		}
	}()

	go func() {
		for i := 0; i <= 20; i++ {
			time.Sleep(500 * time.Millisecond)
			var event interface{} = i
			err := a.Send(event)
			if err == events.ErrClosed {
				break
			}
		}
		wg.Done()
		defer a.Close()
	}()

	wg.Wait()
}
