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
 *     Initial: 2017/09/19        Li Zebang
 */

package search

import (
	"sort"
)

// Data -
type Data struct {
	Items []interface{}
}

// Compare -
type Compare interface {
	Less(interface{}) bool
}

// Arith -
type Arith interface {
	Sub(interface{}) float64
}

// Len is in order to implement the sort interface.
func (d Data) Len() int {
	return len(d.Items)
}

// Less is in order to implement the sort interface.
func (d Data) Less(i, j int) bool {
	return d.Items[i].(Compare).Less(d.Items[j])
}

// Swap is in order to implement the sort interface.
func (d Data) Swap(i, j int) {
	d.Items[i], d.Items[j] = d.Items[j], d.Items[i]
}

// LinearSearch returns whether there is a matching item.
func (d *Data) LinearSearch(key interface{}) bool {
	for _, item := range d.Items {
		if item == key {
			return true
		}
	}
	return false
}

// BinarySearch sorts Data's Items in ascending order and
// returns whether there is a matching item.
func (d *Data) BinarySearch(key interface{}) bool {
	sort.Sort(d)

	head := 0
	tail := len(d.Items) - 1

	for head <= tail {
		median := (head + tail) / 2

		if d.Items[median].(Compare).Less(key) {
			head = median + 1
		} else {
			tail = median - 1
		}
	}

	if head == len(d.Items) || d.Items[head] != key {
		return false
	}

	return true
}

// InterpolationSearch sorts Data's Items in ascending order
// and returns the index of matched item. If it returns -1,
// there is no item matched the key.
func (d *Data) InterpolationSearch(key interface{}) int {
	sort.Sort(d)

	head, tail := 0, len(d.Items)-1
	min, max := d.Items[head], d.Items[tail]

	for {
		if key.(Compare).Less(min) {
			if d.Items[head] != key {
				return -1
			}

			return head
		}

		if max.(Compare).Less(key) {
			if d.Items[tail] != key {
				return -1
			}

			return tail + 1
		}

		var guess int
		if head == tail {
			guess = tail
		} else {
			size := tail - head
			offset := int(float64(size-1) * (key.(Arith).Sub(min) / max.(Arith).Sub(min)))
			guess = head + offset
		}

		if d.Items[guess] == key {
			for guess > 0 && d.Items[guess-1] == key {
				guess--
			}
			return guess
		}

		if key.(Compare).Less(d.Items[guess]) {
			tail = guess - 1
			max = d.Items[tail]
		} else {
			head = guess + 1
			min = d.Items[head]
		}
	}
}
