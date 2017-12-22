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
	"testing"
)

type Int int

func (i Int) Less(j interface{}) bool {
	return i < j.(Int)
}

func (i Int) Sub(j interface{}) float64 {
	return float64(i - j.(Int))
}

func TestLinearSearch(t *testing.T) {
	d := Data{
		Items: []interface{}{Int(61), Int(72), Int(53), Int(94), Int(65)},
	}

	match := d.LinearSearch(Int(53))
	if match {
		t.Log("Linear Search Result:", match)
	} else {
		t.Fatal("Linear Search Error!")
	}

	match = d.LinearSearch(Int(0))
	if !match {
		t.Log("Linear Search Result:", match)
	} else {
		t.Fatal("Linear Search Error!")
	}
}

func TestBinarySearch(t *testing.T) {
	d := Data{
		Items: []interface{}{Int(61), Int(72), Int(53), Int(94), Int(65)},
	}

	match := d.BinarySearch(Int(53))
	if match {
		t.Log("Binary Search Result:", match)
	} else {
		t.Fatal("Binary Search Error!")
	}

	match = d.BinarySearch(Int(0))
	if !match {
		t.Log("Binary Search Result:", match)
	} else {
		t.Fatal("Binary Search Error!")
	}
}

func TestInterpolationSearch(t *testing.T) {
	d := Data{
		Items: []interface{}{Int(61), Int(72), Int(53), Int(94), Int(65)},
	}

	match := d.InterpolationSearch(Int(72))
	if match == 3 {
		t.Log("Interpolation Search Result:", match)
	} else {
		t.Fatal("Interpolation Search Error!")
	}

	match = d.InterpolationSearch(Int(92))
	if match == -1 {
		t.Log("Interpolation Search Result:", match)
	} else {
		t.Fatal("Interpolation Search Error!")
	}
}
