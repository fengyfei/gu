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
	"time"

	"github.com/asciimoo/colly"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	store "github.com/fengyfei/gu/libs/store/badger"
	"golang.org/x/net/html"
)

type GoCN struct {
	Date    string            `json:"time"`
	URL     string            `json:"url"`
	Content map[string]string `json:"content"`
}

type ok struct{}

type gocnCrawler struct {
	collectorURL  *colly.Collector
	collectorNews *colly.Collector
	errCh         chan error
	urlCh         chan *string
	urlFinish     chan ok
	readyCh       chan ok
	newsCh        chan *GoCN
	oldIncr       string
	newIncr       string
}

const (
	noneBadger = "none badger"
)

var (
	invalidKey = []string{"/", ".", "\"", "$", "*", "<", ">", ":", "|", "?"}
)

// NewGoCNCrawler generates a crawler for gocn news.
func NewGoCNCrawler(ch chan *GoCN) crawler.Crawler {
	return &gocnCrawler{
		collectorURL:  colly.NewCollector(),
		collectorNews: colly.NewCollector(),
		errCh:         make(chan error),
		urlCh:         make(chan *string),
		urlFinish:     make(chan ok),
		readyCh:       make(chan ok),
		newsCh:        ch,
	}
}

// Crawler interface Init
func (c *gocnCrawler) Init() error {
	c.collectorURL.OnHTML("div.aw-common-list", c.parseURL)
	c.collectorNews.OnHTML("div.aw-mod.aw-question-detail", c.parseNews)

	return nil
}

// Crawler interface Start
func (c *gocnCrawler) Start() error {
	err := c.prepare()
	if err != nil {
		logger.Error("Error in preparing to start:", err)
		return err
	}

	go func() {
		c.readyCh <- ok{}
	}()

	go c.startURL()
	go c.startNews()

	for {
		select {
		case err = <-c.errCh:
			if err != nil {
				return err
			}
		case <-time.NewTimer(10 * time.Second).C:
			err = c.finish()
			if err != nil {
				logger.Error("Error in the end:", err)
			}
			return nil
		}
	}
}

func (c *gocnCrawler) prepare() error {
	db, err := store.NewBadgerDB(options.FileIO, "./gocn-news-badger", true)
	if err != nil {
		return err
	}

	value, err := db.Get([]byte("increment-key"))
	if len(value) != 0 && err == nil {
		c.oldIncr = string(value)
		return nil
	}

	if err != badger.ErrKeyNotFound {
		return err
	}

	c.oldIncr = noneBadger
	return nil
}

func (c *gocnCrawler) finish() error {
	db, err := store.NewBadgerDB(options.FileIO, "./gocn-news-badger", true)
	if err != nil {
		return err
	}

	return db.Set([]byte("increment-key"), []byte(c.newIncr))
}

func (c *gocnCrawler) startURL() {
	var (
		page int
	)

	for {
		select {
		case <-c.readyCh:
			page++
			go func() {
				err := c.collectorURL.Visit(fmt.Sprintf("https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-%d", page))
				if err != nil {
					logger.Error("error in crawling the URL", err)
					c.errCh <- err
				}
			}()
		case <-c.urlFinish:
			goto EXIT
		}
	}

EXIT:
}

func (c *gocnCrawler) parseURL(e *colly.HTMLElement) {
	if e.Request.URL.String() == "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-1" {
		url, _ := e.DOM.Find("div.aw-question-content").Eq(0).Find("a").Attr("href")
		c.newIncr = url
	}
	if e.DOM.Find("div.aw-question-content").Eq(0).Find("a").Text() == "" {
		c.urlFinish <- ok{}
	}
	for i := 0; ; i++ {
		url, b := e.DOM.Find("div.aw-question-content").Eq(i).Find("a").Attr("href")
		if url == c.oldIncr {
			c.urlFinish <- ok{}
		}
		if b && !strings.Contains(e.DOM.Find("div.aw-question-content").Eq(i).Find("a").Text(), "每日新闻") {
			continue
		} else if !b {
			break
		}
		c.urlCh <- &url
	}
	c.readyCh <- ok{}
}

func (c *gocnCrawler) startNews() {
	for u := range c.urlCh {
		url := u
		go func() {
			err := c.collectorNews.Visit(*url)
			if err != nil {
				logger.Error("error in crawling the news", err)
				c.errCh <- err
			}
		}()
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	times := strings.SplitN(e.DOM.Find("h1").Text(), "-", 3)
	data := GoCN{
		Date:    fmt.Sprintf("%s-%s-%s", times[0][len(times[0])-4:], times[1], times[2][:2]),
		URL:     e.Request.URL.String(),
		Content: make(map[string]string),
	}

	s := data.parseNodes(e.DOM.Nodes)

	for i := 0; i < 5; i++ {
		data.Content[validKey(s[2*i])] = s[2*i+1]
	}

	c.newsCh <- &data
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

func validKey(str string) string {
	for _, value := range invalidKey {
		str = strings.Replace(str, value, "", -1)
	}
	return str
}
