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

type gocnCrawler struct {
	collectorURL  *colly.Collector
	collectorNews *colly.Collector

	dataCh   chan *crawler.Data
	finishCh chan struct{}

	errCh   chan error
	newsURL chan *string

	visitSiteReady  chan struct{}
	visitSiteFinish chan struct{}
	crawlerFinish   chan struct{}

	db      *store.BadgerDB
	oldIncr string
	newIncr string
}

const (
	site           = "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-"
	defaultOldIncr = "https://gocn.io/question/731"
)

// NewGoCNCrawler generates a crawler for gocn news.
func NewGoCNCrawler(dataCh chan *crawler.Data, finishCh chan struct{}) crawler.Crawler {
	return &gocnCrawler{
		collectorURL:  colly.NewCollector(),
		collectorNews: colly.NewCollector(),

		dataCh:   dataCh,
		finishCh: finishCh,

		errCh:   make(chan error),
		newsURL: make(chan *string),

		visitSiteReady:  make(chan struct{}),
		visitSiteFinish: make(chan struct{}),
		crawlerFinish:   make(chan struct{}),
	}
}

// Crawler interface Init
func (c *gocnCrawler) Init() error {
	c.collectorURL.OnHTML("div.aw-common-list", c.parseURL)
	c.collectorNews.OnHTML("div.aw-mod.aw-question-detail", c.parseNews)

	err := c.prepare()
	if err != nil {
		logger.Error("Error in preparing to start:", err)
		return err
	}
	return nil
}

// Crawler interface Start
func (c *gocnCrawler) Start() error {
	go func() {
		c.visitSiteReady <- struct{}{}
	}()

	go c.startURL()
	go c.startNews()

	for {
		select {
		case err := <-c.errCh:
			return err
		case <-c.crawlerFinish:
			err := c.finish()
			if err != nil {
				logger.Error("Error in the end:", err)
			}
			c.finishCh <- struct{}{}
			return nil
		}
	}
}

func (c *gocnCrawler) prepare() error {
	db, err := store.NewBadgerDB(options.FileIO, "gocn-news-badger", true)
	if err != nil {
		return err
	}
	c.db = db

	value, err := c.db.Get([]byte("increment-key"))
	if len(value) != 0 && err == nil {
		c.oldIncr = string(value)
		return nil
	}

	if err != badger.ErrKeyNotFound {
		return err
	}

	c.oldIncr = defaultOldIncr
	return nil
}

func (c *gocnCrawler) finish() error {
	return c.db.Set([]byte("increment-key"), []byte(c.newIncr))
}

func (c *gocnCrawler) startURL() {
	var page int

	for {
		select {
		case <-c.visitSiteReady:
			page++
			go func() {
				err := c.collectorURL.Visit(fmt.Sprintf("%s%d", site, page))
				if err != nil {
					logger.Error("Error in crawling the URL", err)
					c.errCh <- err
				}
			}()
		case <-c.visitSiteFinish:
			return
		}
	}
}

func (c *gocnCrawler) parseURL(e *colly.HTMLElement) {
	if e.Request.URL.String() == fmt.Sprintf("%s%d", site, 1) {
		url, _ := e.DOM.Find("div.aw-question-content").Eq(0).Find("a").Attr("href")
		c.newIncr = url
	}

	if e.DOM.Find("div.aw-question-content").Eq(0).Find("a").Text() == "" {
		c.visitSiteFinish <- struct{}{}
	}

	for i := 0; ; i++ {
		url, b := e.DOM.Find("div.aw-question-content").Eq(i).Find("a").Attr("href")
		if b && !strings.Contains(e.DOM.Find("div.aw-question-content").Eq(i).Find("a").Text(), "每日新闻") {
			continue
		} else if !b {
			break
		}
		c.newsURL <- &url
		if url == c.oldIncr {
			c.visitSiteFinish <- struct{}{}
			return
		}
	}

	c.visitSiteReady <- struct{}{}
}

func (c *gocnCrawler) startNews() {
	for {
		select {
		case url := <-c.newsURL:
			err := c.collectorNews.Visit(*url)
			if err != nil {
				logger.Error("Error in crawling the news", err)
				c.errCh <- err
			}
		case <-time.NewTimer(time.Second).C:
			return
		}
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	if c.oldIncr == e.Request.URL.String() && c.oldIncr == defaultOldIncr {
		c.crawlerFinish <- struct{}{}
	} else if c.oldIncr == e.Request.URL.String() {
		c.crawlerFinish <- struct{}{}
		return
	}

	times := strings.SplitN(e.DOM.Find("h1").Text(), "-", 3)
	data := &crawler.Data{
		Source: "GoCN Daily News",
		Date:   fmt.Sprintf("%s-%s-%s", parseTime(times[0][len(times[0])-4:]), parseTime(times[1]), parseTime(times[2][:2])),
		URL:    e.Request.URL.String(),
	}
	data.Title = "GoCN 每日新闻 " + data.Date

	element, text := parseNodes(e.DOM.Nodes)
	for i := 0; i < 5; i++ {
		if i > len(text)-1 {
			break
		} else if i > len(element)-1 {
			data.Text += fmt.Sprintf("%s\n", text[i])
		} else {
			data.Text += fmt.Sprintf("%s %s\n", text[i], element[i])
		}
	}

	c.dataCh <- data
}

func parseNodes(s []*html.Node) ([]string, []string) {
	var (
		f           func(*html.Node)
		elementPipe = make(chan string)
		textPipe    = make(chan string)
		element     []string
		text        []string
	)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, v := range n.Attr {
				if v.Key == "href" {
					elementPipe <- v.Val
				}
			}
		}

		if n.Type == html.TextNode {
			text := n.Data
			if strings.Count(text, string(byte(9)))+strings.Count(text, string(byte(10)))+strings.Count(text, string(byte(32))) != len(text) && !strings.Contains(text, "每日新闻") && !strings.Contains(text, "http://") && !strings.Contains(text, "https://") {
				textPipe <- text
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

	for {
		select {
		case e := <-elementPipe:
			element = append(element, e)
			if len(element)+len(text) > 9 {
				return element, text
			}
		case t := <-textPipe:
			text = append(text, t)
			if len(element)+len(text) > 9 {
				return element, text
			}
		}
	}
}

func parseTime(s string) string {
	var time string
	ss := strings.Split(s, "")
	for _, v := range ss {
		if !strings.Contains("0123456789", v) {
			break
		}
		time += v
	}
	if len(time) < 2 {
		time = "0" + time
	}
	return time
}
