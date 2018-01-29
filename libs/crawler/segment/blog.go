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
 *     Initial: 2017/01/21        Li Zebang
 */

package segment

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fengyfei/gu/libs/logger"
)

type Blog struct {
	Title string
	Blog  string
	Label []string
}

var (
	blogDonePipe chan done         = make(chan done)
	topMap       map[string]string = map[string]string{"&amp;": "&", "&#x27;": "'"}
	bodyMap      map[string]string = map[string]string{"&amp;": "&", "&#x27;": "'"}
)

func (c *segmentCrawler) startBlog() {
	for {
		select {
		case url := <-urlPipe:
			err := c.getBlog(url)
			if err != nil {
				errorPipe <- err
			}
			blogDonePipe <- done{}
		case <-time.NewTimer(3 * time.Second).C:
			goto EXIT
		}
	}

EXIT:
}

func (c *segmentCrawler) getBlog(url *string) error {
	cli := &http.Client{}

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		logger.Error("error in getting blog", err)
		return err
	}

	resp, err := cli.Do(req)
	if err != nil {
		logger.Error("error in getting blog", err)
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error in getting blog", err)
		return err
	}

	b := &Blog{}
	b.parseBlog(string(data))

	return nil
}

func (b *Blog) parseBlog(s string) {
	s = strings.SplitN(s, "<h1 class=\"Article-title\" data-reactid=\"39\">", 2)[1]
	s = strings.SplitN(s, "<footer class=\"Article-footer\" data-reactid=\"", 2)[0]
	data := strings.SplitN(s, "<div class=\"Article-body Content\" data-swiftype-name=\"body\" data-swiftype-type=\"text\" data-reactid=", 2)

	b.parseTop(data[0])
	b.parseBody(data[1])

	f, _ := os.OpenFile("./blog/"+b.Title+".md", os.O_CREATE|os.O_RDWR, 0644)
	f.Write([]byte(b.Blog))
}

func (b *Blog) parseTop(s string) {
	text := strings.SplitN(s, "</h1>", 2)
	b.Title = text[0]
	for k, v := range topMap {
		strings.Replace(b.Title, k, v, -1)
	}
	b.Blog += "# " + b.Title + "\n\n"

	count := strings.Count(text[1], "<a class=\"Author-name\" href=\"")
	for i := 0; i < count; i++ {
		if i != 0 {
			b.Blog += " and "
		}
		text = strings.SplitN(text[1], "<a class=\"Author-name\" href=\"", 2)
		text = strings.SplitN(text[1], "\"", 2)
		url := site + text[0]
		text = strings.SplitN(text[1], ">", 2)
		text = strings.SplitN(text[1], "<", 2)
		name := text[0]
		b.Blog = fmt.Sprintf("%s[%s](%s)", b.Blog, name, url)
	}

	text = strings.SplitN(text[1], "<!-- /react-text -->", 3)
	text = strings.SplitN(text[1], "-->", 2)
	date := text[1]
	b.Blog = fmt.Sprintf("%s on %s\n", b.Blog, date)
}

func (b *Blog) parseBody(s string) {
	b.Blog += parseMD(s[5 : len(s)-6])
}
