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
 *     Initial: 2018/01/28        Li Zebang
 */

package segment

import (
	"github.com/asciimoo/colly"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/util/bolt"
	"github.com/fengyfei/gu/libs/logger"
)

type segmentCrawler struct {
	collector  *colly.Collector
	dataCh     chan crawler.Data
	finishCh   chan struct{}
	errCh      chan error
	db         *bolt.DB
	earlierURL string
	currentURL string
}

const (
	key  = "segment"
	site = "https://segment.com/blog/"
)

// NewSegmentCrawler generates a crawler for Segment blogs.
func NewSegmentCrawler(dataCh chan crawler.Data, finishCh chan struct{}) crawler.Crawler {
	c := colly.NewCollector()
	return &segmentCrawler{
		collector: c,
		dataCh:    dataCh,
		finishCh:  finishCh,
		errCh:     make(chan error),
	}
}

func (c *segmentCrawler) prepare() error {
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

func (c *segmentCrawler) shutdown() error {
	return c.db.Set([]byte(crawler.CrawlerBucket), []byte(key), []byte(c.currentURL))
}

// Crawler interface Init
func (c *segmentCrawler) Init() error {
	c.collector.OnHTML("section.Articles-list.clearfix", c.visitBlog)
	c.collector.OnHTML("div.Pagination", c.visitNext)

	return c.prepare()
}

// Crawler interface Start
func (c *segmentCrawler) Start() error {
	defer close(c.dataCh)

	go func() {
		err := c.collector.Visit(site)
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

func (c *segmentCrawler) visitNext(e *colly.HTMLElement) {
	subURL, exist := e.DOM.Find("a.Link--primary.Link--animatedHover.Pagination-older").Attr("href")
	if exist {
		err := c.collector.Visit(e.Request.AbsoluteURL(subURL))
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

func (c *segmentCrawler) visitBlog(e *colly.HTMLElement) {
	selector := "a.Link--primary.Link--animatedHover.ArticleInList-readMoreLink"
	if e.Request.URL.String() == site {
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
		if ok {
			err := c.parseBlog(e.Request.AbsoluteURL(subURL))
			if err != nil {
				logger.Error("error in visiting a blog", err)
				c.errCh <- err
			}
		}
	}
}
