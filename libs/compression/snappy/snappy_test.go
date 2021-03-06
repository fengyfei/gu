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

package snappy

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	var (
		dst, want []byte
		err       error
	)
	src := []byte("Hello world, with gosnappy, have a good day.")

	if dst, err = Encode(src); err != nil {
		t.Fatal(err)
	}

	fmt.Println("Source length:", len(src), ",Encoded length:", len(dst))

	if want, err = Decode(dst); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(src, want) {
		t.Fatal("Not same with original")
	}
}

func TestEncoder(t *testing.T) {
	var (
		src, dst *os.File
		err      error
	)

	// curl www.sina.com.cn > sina.html
	src, err = os.Open("./sina.html")
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	dst, err = os.Create("./sina.enc")
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()

	encoder := NewEncoder(src, dst)
	if err = encoder.Encode(); err != nil {
		t.Fatal(err)
	}
}

func TestDecoder(t *testing.T) {
	var (
		src, dst *os.File
		err      error
	)

	// curl www.sina.com.cn > sina.html
	src, err = os.Open("./sina.enc")
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	dst, err = os.Create("./sina.dec")
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()

	decoder := NewDecoder(src, dst)
	if err = decoder.Decode(); err != nil {
		t.Fatal(err)
	}
}
