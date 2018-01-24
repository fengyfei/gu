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
	collectorURL  *colly.Collector
	collectorNews *colly.Collector
}

type GoCN struct {
	Time    string            `json:"time"`
	URL     string            `json:"url"`
	Content map[string]string `json:"content"`
}

type ok struct{}

var (
	counterURL int
	invalidKey [10]string  = [10]string{"/", ".", "\"", "$", "*", "<", ">", ":", "|", "?"}
	invalidLi  [5]string   = [5]string{"\nGopherChina2018来了！ https://www.bagevent.com/event/1086224\n", "GopherChina Telegram群现已上线 https://t.me/gopherchina ", " 微博", " QZONE", " 微信"}
	errorPipe  chan error  = make(chan error)
	urlPipe    chan string = make(chan string)
	overURL    chan ok     = make(chan ok)
	overNews   chan ok     = make(chan ok)
	readyPipe  chan ok     = make(chan ok)
	DataPipe   chan *GoCN  = make(chan *GoCN)
)

func NewGoCNCrawler() crawler.Crawler {
	return &gocnCrawler{
		collectorURL:  colly.NewCollector(),
		collectorNews: colly.NewCollector(),
	}
}

func (c *gocnCrawler) Init() error {
	c.collectorURL.OnHTML("a", c.parseURL)
	c.collectorNews.OnHTML("div.aw-mod", c.parseNews)

	return nil
}

func (c *gocnCrawler) Start() error {
	go func() {
		readyPipe <- ok{}
	}()

	go c.startURL()
	go c.startNews()

	for {
		select {
		case err := <-errorPipe:
			if err != nil {
				return err
			}
		case <-overNews:
			return nil
		}
	}
}

func (c *gocnCrawler) startURL() {
	var (
		page int
	)

	for {
		select {
		case <-readyPipe:
			page++
			go func() {
				err := c.collectorURL.Visit(fmt.Sprintf("https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-%d", page))
				if err != nil {
					errorPipe <- err
				}
			}()
		case <-overURL:
			goto EXIT
		}
	}

EXIT:
	overNews <- ok{}
}

func (c *gocnCrawler) parseURL(e *colly.HTMLElement) {
	if strings.Contains(e.Text, "每日新闻") {
		if e.Attr("href") == "https://gocn.io/explore/category-14" {
			counterURL += 100
		} else if strings.Contains(e.Attr("href"), "https://gocn.io/topic/") {
			counterURL += 100
		} else {
			counterURL++
			url := e.Attr("href")
			urlPipe <- url
		}

		if counterURL > 200 {
			counterURL = 0
			readyPipe <- ok{}
		} else if counterURL == 200 {
			counterURL = 0
			overURL <- ok{}
		}
	}
}

func (c *gocnCrawler) startNews() {
	for {
		select {
		case url := <-urlPipe:
			go func() {
				err := c.collectorNews.Visit(url)
				if err != nil {
					errorPipe <- err
				}
			}()
		case <-time.NewTimer(3 * time.Second).C:
			goto EXIT
		}
	}

EXIT:
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	if strings.Contains(e.Attr("class"), "aw-mod aw-question-detail") {
		var data = &GoCN{
			Time:    parseTitle(e.DOM.Find("h1").Text()),
			URL:     e.Request.URL.String(),
			Content: make(map[string]string),
		}

		for i := 0; ; i++ {
			url, _ := e.DOM.Find("a").Eq(i).Attr("href")
			text := e.DOM.Find("li").Eq(i).Text()

			for _, value := range invalidLi {
				if strings.Contains(text, value) {
					text = ""
					break
				}
			}
			if text == "" {
				text = e.DOM.Find("p").Eq(i + 1).Text()
			}

			index := strings.Index(text, "http://") + strings.Index(text, "https://") + 1
			if index > 0 {
				text = string([]byte(text)[:index])
			}

			if url == "" || text == "" {
				data.parseNews("p", e)
				if len(data.Content) != 5 {
					data.parseNews("code", e)
				}
				break
			}

			data.Content[validKey(text)] = url
			if len(data.Content) == 5 {
				break
			}
		}

		DataPipe <- data
	}
}

func (data *GoCN) parseNews(query string, e *colly.HTMLElement) {
	data.Content = make(map[string]string)
	urls := strings.Split(e.DOM.Find("a").Text(), "http")
	for k, v := range strings.Split(e.DOM.Find(query).Text(), "http") {
		if strings.Contains(v, "每日新闻") {
			data.Content[validKey(strings.Split(v, ")")[1])] = "http" + urls[k+1]
		} else {
			vs := strings.Split(v, "\n")
			data.Content[validKey(vs[len(vs)-1])] = "http" + urls[k+1]
		}
		if k == 4 {
			break
		}
	}
}

func parseTitle(title string) string {
	return strings.Split(strings.Split(title, "(")[1], ")")[0]
}

func validKey(str string) string {
	for _, value := range invalidKey {
		str = strings.Replace(str, value, "", -1)
	}
	return str
}
