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
 *     Initial: 2018/04/23        Li Zebang
 */

package segment

import (
	"fmt"
	"os"
	"testing"

	"github.com/fengyfei/gu/libs/crawler"
)

func TestSegmentCrawler(t *testing.T) {
	var (
		dataCh   = make(chan crawler.Data)
		finishCh = make(chan struct{})
	)
	c := NewSegmentCrawler(dataCh, finishCh)
	go func() {
		err := crawler.StartCrawler(c)
		if err != nil {
			panic(err)
		}
	}()

	err := os.MkdirAll("blog", 0755)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case data := <-dataCh:
			blog := data.(*SegmentData)
			file, err := os.OpenFile(fmt.Sprintf("blog/%s.md", blog.Title), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			n, err := file.Write([]byte(blog.Text))
			if n != len(blog.Text) || err != nil {
				panic(err)
			}
		case <-finishCh:
			fmt.Println(2)
			return
		}
	}
}
