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
	urlPipe    chan string = make(chan string)
	readyPipe  chan bool   = make(chan bool)
	DataPipe   chan *GoCN  = make(chan *GoCN)
)

func NewGoCNCrawler() crawler.Crawler {
	return &gocnCrawler{
		collector: colly.NewCollector(),
	}
}

func (c *gocnCrawler) Init() error {
	go func() {
		urlPipe <- "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-1"
	}()
	fmt.Println(123456)
	return nil
}

func (c *gocnCrawler) Start() error {
	var page int
	go func() {
		for {
			if url, ok := <-urlPipe; ok {
				fmt.Println(url)
				c.init(url)

				if strings.Contains(url, "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-") {
					page++
					urlPipe <- fmt.Sprintf("https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-%d", page)
				}
				go func() {
					c.start(url)
				}()
			}
		}
	}()
	for {

	}

	return nil
}

func (c *gocnCrawler) parseURL(e *colly.HTMLElement) {
	if strings.Contains(e.Text, "每日新闻") {
		if e.Attr("href") != "https://gocn.io/explore/category-14" && e.Attr("href") != "https://gocn.io/topic/每日新闻" {
			url := e.Text + "\t" + e.Attr("href")
			urlPipe <- url
		}
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	var cn = &GoCN{
		URL:     e.Request.URL.String(),
		Content: make(map[string]string),
	}
	val, exists := e.DOM.Find("a").Attr("href")
	if exists {
		index := strings.Index(e.Text, "http://") + strings.Index(e.Text, "https://")
		if index > 0 {
			cn.Content[validKey(string([]byte(e.Text)[:index]))] = val
		}
	}
	DataPipe <- cn
}

func validKey(str string) string {
	for _, value := range invalidKey {
		str = strings.Replace(str, value, "", -1)
	}

	return str
}

func (c *gocnCrawler) init(url string) {
	if strings.Contains(url, "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-") {
		c.collector.OnHTML("a", c.parseURL)
	} else if strings.Contains(url, "https://gocn.io/question/") {
		c.collector.OnHTML("li", c.parseNews)
	} else {
		c.collector.OnHTML("p", c.parseNews)
	}
}

func (c *gocnCrawler) start(url string) error {
	return c.collector.Visit(url)
}
