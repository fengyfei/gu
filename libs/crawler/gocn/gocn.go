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
	"strconv"
	"strings"

	"github.com/asciimoo/colly"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/util/bolt"
	"github.com/fengyfei/gu/libs/logger"
)

// GocnData -
type GocnData struct {
	// Source -
	Source string
	// Date -
	Date string
	// Title -
	Title string
	// URL -
	URL string
	// Text -
	Text []News
	// FileType -
	FileType string
}

// News -
type News struct {
	// Title -
	Title string
	// URL -
	URL string
}

func (data *GocnData) String() string {
	var text string
	for _, news := range data.Text {
		text += fmt.Sprintf("%s %s\n", news.Title, news.URL)
	}
	return fmt.Sprintf("Source: %s\nDate: %s\nTitle: %s\nURL: %s\n%s", data.Source, data.Date, data.Title, data.URL, text)
}

// File return "", "", "".
func (data *GocnData) File() (title, filetype, content string) {
	return "", "", ""
}

// IsFile return false.
func (data *GocnData) IsFile() bool {
	return false
}

type gocnCrawler struct {
	collector       *colly.Collector
	detailCollector *colly.Collector
	dataCh          chan crawler.Data
	finishCh        chan struct{}
	errCh           chan error
	db              *bolt.DB
	earlierURL      string
	currentURL      string
}

const (
	key  = "gocn"
	site = "https://gocn.io/sort_type-new__category-14__day-0__is_recommend-0__page-%d"
)

// NewGoCNCrawler generates a crawler for gocn news.
func NewGoCNCrawler(dataCh chan crawler.Data, finishCh chan struct{}) crawler.Crawler {
	c := colly.NewCollector()
	return &gocnCrawler{
		collector:       c,
		detailCollector: c.Clone(),
		dataCh:          dataCh,
		finishCh:        finishCh,
		errCh:           make(chan error),
	}
}

func (c *gocnCrawler) prepare() error {
	db, err := bolt.Open(crawler.CrawlerPath)
	if err != nil {
		logger.Error("error in opening a boltdb", err)
		return err
	}
	c.db = db

	value, err := c.db.Get([]byte(crawler.CrawlerBucket), []byte(key))
	if err != nil {
		logger.Error("error in getting earlier url", err)
		return err
	}
	c.earlierURL = string(value)
	return nil
}

func (c *gocnCrawler) shutdown() error {
	return c.db.Set([]byte(crawler.CrawlerBucket), []byte(key), []byte(c.currentURL))
}

// Crawler interface Init
func (c *gocnCrawler) Init() error {
	c.detailCollector.OnHTML("div.aw-mod.aw-question-detail.aw-item", c.parseNews)
	c.collector.OnHTML("div.aw-mod.aw-explore-list", c.visitNews)
	c.collector.OnHTML("ul.pagination.pull-right", c.visitNext)

	return c.prepare()
}

// Crawler interface Start
func (c *gocnCrawler) Start() error {
	defer close(c.dataCh)

	go func() {
		err := c.collector.Visit(fmt.Sprintf(site, 1))
		if err != nil {
			logger.Error("error in starting a visit", err)
			c.errCh <- err
		}
	}()

	for err := range c.errCh {
		close(c.errCh)
		return err
	}

	close(c.errCh)
	return nil
}

func (c *gocnCrawler) visitNews(e *colly.HTMLElement) {
	selector := "div.aw-item > div.aw-question-content > h4 > a"
	if e.Request.URL.String() == fmt.Sprintf(site, 1) {
		c.currentURL, _ = e.DOM.Find(selector).Eq(0).Attr("href")
	}
	if c.currentURL == c.earlierURL {
		defer close(c.finishCh)
		err := c.shutdown()
		if err != nil {
			logger.Error("error in shutdown", err)
			c.errCh <- err
		}
		return
	}

	var (
		subURL = ""
		ok     = true
	)
	for index := 0; ok; index++ {
		subURL, ok = e.DOM.Find(selector).Eq(index).Attr("href")
		if !strings.Contains(e.DOM.Find(selector).Eq(index).Text(), "每日新闻") {
			continue
		}
		if ok {
			err := c.detailCollector.Visit(e.Request.AbsoluteURL(subURL))
			if err != nil {
				logger.Error("error in visiting a blog", err)
				c.errCh <- err
			}
		}
	}
}

func (c *gocnCrawler) visitNext(e *colly.HTMLElement) {
	if e.Text[len(e.Text)-1:len(e.Text)] == ">" {
		page, err := strconv.Atoi(e.Request.URL.String()[len(site)-2:])
		if err != nil {
			logger.Error("error in preparing for a visit", err)
			c.errCh <- err
		}
		err = c.collector.Visit(fmt.Sprintf(site, page+1))
		if err != nil {
			logger.Error("error in starting a visit", err)
			c.errCh <- err
		}
	} else {
		defer close(c.finishCh)
		err := c.shutdown()
		if err != nil {
			logger.Error("error in shutdown", err)
			c.errCh <- err
		}
	}
}

func (c *gocnCrawler) parseNews(e *colly.HTMLElement) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic when parsing gocn news", err)
			c.errCh <- fmt.Errorf("%v", err)
		}
	}()

	getDate := func(head string) string {
		getNum := func(str string) string {
			var s string
			slice := strings.Split(str, "")
			for _, v := range slice {
				if !strings.Contains("0123456789", v) {
					break
				}
				s += v
			}
			if len(s) < 2 {
				s = "0" + s
			}
			return s
		}

		ts := strings.SplitN(head, "-", 3)
		return fmt.Sprintf("%s-%s-%s", getNum(ts[0][len(ts[0])-4:len(ts[0])]), getNum(ts[1]), getNum(ts[2][:2]))
	}

	getText := func(e *colly.HTMLElement) []News {
		getTitle := func(s string) string {
			index := strings.Index(s, "http://") + strings.Index(s, "https://")
			if index < -1 {
				return strings.TrimSpace(s)
			}
			return strings.TrimSpace(s[:index+1])
		}

		news := []News{}
		url := ""
		ok := true
		for index := 0; ok && index < 5; index++ {
			url, ok = e.DOM.Find("a").Eq(index).Attr("href")
			news = append(news, News{getTitle(e.DOM.Find("li").Eq(index).Text()), url})
		}
		return news
	}

	c.dataCh <- &GocnData{
		Source: "GoCN Daily News",
		Date:   getDate(e.DOM.Find("h1").Text()),
		Title:  strings.TrimSpace(e.DOM.Find("h1").Text()),
		URL:    e.Request.URL.String(),
		Text:   getText(e),
	}
}
