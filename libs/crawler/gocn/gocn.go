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
 *     Initial: 2018/01/21        Li Zebang
 */

package gocn

import (
	"fmt"
	"strings"

	"github.com/asciimoo/colly"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	"golang.org/x/net/html"
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
	invalidKey [10]string = [10]string{"/", ".", "\"", "$", "*", "<", ">", ":", "|", "?"}
	// invalidLi  [5]string   = [5]string{"\nGopherChina2018来了！ https://www.bagevent.com/event/1086224\n", "GopherChina Telegram群现已上线 https://t.me/gopherchina ", " 微博", " QZONE", " 微信"}
	errorPipe chan error  = make(chan error)
	urlPipe   chan string = make(chan string)
	overURL   chan ok     = make(chan ok)
	overNews  chan ok     = make(chan ok)
	readyPipe chan ok     = make(chan ok)
	DataPipe  chan *GoCN  = make(chan *GoCN)
)

// NewGoCNCrawler generates a crawler for gocn news.
func NewGoCNCrawler() crawler.Crawler {
	return &gocnCrawler{
		collectorURL:  colly.NewCollector(),
		collectorNews: colly.NewCollector(),
	}
}

// Crawler interface Init
func (c *gocnCrawler) Init() error {
	c.collectorURL.OnHTML("a", c.parseURL)
	c.collectorNews.OnHTML("div.aw-mod.aw-question-detail", c.parseNews)

	return nil
}

// Crawler interface Start
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
					logger.Error("error in crawling the URL", err)
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
	for u := range urlPipe {
		url := u
		go func() {
			err := c.collectorNews.Visit(url)
			if err != nil {
				logger.Error("error in crawling the news", err)
				errorPipe <- err
			}
		}()
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	times := strings.SplitN(e.DOM.Find("h1").Text(), "-", 3)
	data := GoCN{
		Time:    fmt.Sprintf("%s-%s-%s", times[0][len(times[0])-4:], times[1], times[2][:2]),
		URL:     e.Request.URL.String(),
		Content: make(map[string]string),
	}

	s := data.parseNodes(e.DOM.Nodes)

	for i := 0; i < 5; i++ {
		data.Content[validKey(s[2*i])] = s[2*i+1]
	}

	DataPipe <- &data
}

func (g *GoCN) parseNodes(s []*html.Node) []string {
	var (
		f          func(*html.Node)
		stringPipe = make(chan string)
		news       []string
	)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, v := range n.Attr {
				if v.Key == "href" {
					stringPipe <- v.Val
				}
			}
		}

		if n.Type == html.TextNode {
			text := n.Data
			if strings.Count(text, string(byte(9)))+strings.Count(text, string(byte(10)))+strings.Count(text, string(byte(32))) != len(text) && !strings.Contains(text, "每日新闻") && !strings.Contains(text, "http://") && !strings.Contains(text, "https://") {
				stringPipe <- text
			}
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	for _, n := range s {
		go f(n)
	}

	for s := range stringPipe {
		news = append(news, s)
		if len(news) > 9 {
			break
		}
	}

	return news
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
