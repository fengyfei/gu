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
 *     Initial: 2018/04/20        Li Zebang
 */

package increment

import (
	"strconv"
	"testing"
	"time"
)

func TestIncrementNoDone(t *testing.T) {
	incr := NewIncrement("test.db", "test")

	err := incr.StartIncrement()
	if err != nil {
		panic(err)
	}
	defer incr.Close()

	go func() {
		for i := 0; i < 1000; i++ {
			err := incr.SendSignal(&Signal{
				IncrementKey: strconv.Itoa(i),
				Done:         false,
			})
			if err != nil {
				panic(err)
			}
		}
	}()

	time.Sleep(1e9)
}

func TestIncrementRepeatDone(t *testing.T) {
	incr := NewIncrement("test.db", "test")

	err := incr.StartIncrement()
	if err != nil {
		panic(err)
	}
	defer incr.Close()

	go func() {
		for i := 0; i < 1000; i++ {
			err := incr.SendSignal(&Signal{
				IncrementKey: strconv.Itoa(i),
				Done:         false,
			})
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < 1000; i++ {
			err := incr.SendSignal(&Signal{
				IncrementKey: strconv.Itoa(i + 1000),
				Done:         true,
			})
			if err != nil && err != ErrSignalChClosed {
				panic(err)
			}
		}
	}()

	time.Sleep(1e9)
}

func TestIncrementAllRight(t *testing.T) {
	incr := NewIncrement("test.db", "test")

	err := incr.StartIncrement()
	if err != nil {
		panic(err)
	}
	defer incr.Close()

	go func() {
		for i := 0; i < 1000; i++ {
			err := incr.SendSignal(&Signal{
				IncrementKey: strconv.Itoa(i),
				Done:         false,
			})
			if err != nil {
				panic(err)
			}
		}
		err := incr.SendSignal(&Signal{
			IncrementKey: strconv.Itoa(11),
			Done:         true,
		})
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1e9)
}
