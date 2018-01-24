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

package gocn

import (
	"fmt"
	"strings"
	"time"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/gocolly/colly"
)

type gocnCrawler struct {
	collector *colly.Collector
}

type GoCN struct {
	Time    string            `json:"time"`
	URL     string            `json:"url"`
	Content map[string]string `json:"content"`
}

var (
	invalidKey [10]string  = [10]string{"/", ".", "\"", "$", "*", "<", ">", ":", "|", "?"}
	invalidP               = []string{"活动预告", "编辑", "招聘信息", "订阅新闻", "GopherChina"}
	urlPipe    chan string = make(chan string)
	overPipe   chan bool   = make(chan bool)
	DataPipe   chan *GoCN  = make(chan *GoCN)
)

func NewGoCNCrawler() crawler.Crawler {
	return &gocnCrawler{
		collector: colly.NewCollector(),
	}
}

func (c *gocnCrawler) Init() error {
	c.collector.OnHTML("a", c.parseURL)
	return nil
}

func (c *gocnCrawler) Start() error {
	var (
		page int
	)

	go ready()
	go c.startNews()

	for {
		select {
		case <-overPipe:
			page++
			err := c.start(fmt.Sprintf("%s%d", "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-", page))
			if err != nil {
				return err
			}
		case <-time.NewTimer(3 * time.Second).C:
			goto AAA
		}
	}
AAA:
	fmt.Println("break")
	return nil
}

func ready() {
	for i := 0; i < 20; i++ {
		overPipe <- true
	}
}

func (c *gocnCrawler) parseURL(e *colly.HTMLElement) {
	if strings.Contains(e.Text, "每日新闻") {
		if e.Attr("href") != "https://gocn.io/explore/category-14" && !strings.Contains(e.Attr("href"), "https://gocn.io/topic/") {
			url := e.Attr("href")
			urlPipe <- url
		}
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	var data = &GoCN{
		Time:    e.DOM.Find("p").Eq(0).Text(),
		URL:     e.Request.URL.String(),
		Content: make(map[string]string),
	}
	for i := 0; ; i++ {
		var (
			url  string
			text string
		)
		url = e.DOM.Find("a").Eq(i).Text()
		text = e.DOM.Find("p").Eq(i + 1).Text()
		for _, value := range invalidP {
			if strings.Contains(text, value) {
				text = ""
				break
			}
		}
		if text == "" {
			text = e.DOM.Find("li").Eq(i).Text()
		}
		index := strings.Index(text, "http://") + strings.Index(text, "https://")
		if index > 0 {
			text = string([]byte(text)[:index])
		}
		if url == "" || text == "" {
			break
		}
		data.Content[validKey(text)] = url
		if len(data.Content) == 5 {
			break
		}
	}
	sum++
	fmt.Println("sum", sum, *data)
	DataPipe <- data
}

func validKey(str string) string {
	for _, value := range invalidKey {
		str = strings.Replace(str, value, "", -1)
	}

	return str
}

func (c *gocnCrawler) start(url string) error {
	return c.collector.Visit(url)
}

func (c *gocnCrawler) startNews() error {
	for {
		select {
		case url := <-urlPipe:
			c.startxx(url)
		case <-time.NewTimer(3 * time.Second).C:
			goto AAA
		}
	}
AAA:
	return nil
}

func (c *gocnCrawler) startxx(url string) (err error) {
	c.collector.OnHTML("div.content", c.parseNews)
	go func() {
		err = c.start(url)
		if err != nil {
			return
		}
	}()
	return nil
}
