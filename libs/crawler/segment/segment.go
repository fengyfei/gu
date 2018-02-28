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
	"fmt"
	"net/http"

	"github.com/asciimoo/colly"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	store "github.com/fengyfei/gu/libs/store/badger"
)

type segmentCrawler struct {
	collector *colly.Collector

	dataCh   chan *crawler.Data
	finishCh chan struct{}

	errCh   chan error
	blogURL chan *string

	visitSiteReady  chan struct{}
	visitSiteFinish chan struct{}
	crawlerFinish   chan struct{}

	db      *store.BadgerDB
	oldIncr string
	newIncr string
}

const (
	site           = "https://segment.com"
	defaultOldIncr = "default-old-increment"
)

// NewSegmentCrawler generates a crawler for Segment blogs.
func NewSegmentCrawler(dataCh chan *crawler.Data, finishCh chan struct{}) crawler.Crawler {
	return &segmentCrawler{
		collector: colly.NewCollector(),

		dataCh:   dataCh,
		finishCh: finishCh,

		errCh:   make(chan error),
		blogURL: make(chan *string),

		visitSiteReady:  make(chan struct{}),
		visitSiteFinish: make(chan struct{}),
		crawlerFinish:   make(chan struct{}),
	}
}

// Crawler interface Init
func (c *segmentCrawler) Init() error {
	c.collector.OnHTML("body", c.parseURL)

	err := c.prepare()
	if err != nil {
		logger.Error("Error in preparing to start:", err)
		return err
	}
	return nil
}

// Crawler interface Start
func (c *segmentCrawler) Start() error {

	go func() {
		c.visitSiteReady <- struct{}{}
	}()

	go c.startBlog()

	var page int
	for {
		select {
		case <-c.visitSiteReady:
			var url string
			page++
			if page == 1 {
				url = site + "/blog/"
			} else {
				url = fmt.Sprintf("%s/blog/page/%d", site, page)
			}
			err := c.collector.Visit(url)
			if err != nil {
				if err.Error() == http.StatusText(http.StatusNotFound) {
					err := c.finish()
					if err != nil {
						logger.Error("Error in the end:", err)
					}
					c.finishCh <- struct{}{}
					return nil
				}
				logger.Error("Error in getting blog url", err)
				return err
			}
		case err := <-c.errCh:
			return err
		case <-c.visitSiteFinish:
			return nil
		}
	}
}

func (c *segmentCrawler) prepare() error {
	db, err := store.NewBadgerDB(options.FileIO, "segment-blog-badger", true)
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

func (c *segmentCrawler) finish() error {
	return c.db.Set([]byte("increment-key"), []byte(c.newIncr))
}

func (c *segmentCrawler) parseURL(e *colly.HTMLElement) {
	if e.Request.URL.String() == "https://segment.com/blog/" {
		u, _ := e.DOM.Find("a.Link--primary.Link--animatedHover.ArticleInList-readMoreLink").Eq(0).Attr("href")
		c.newIncr = site + u
	}

	if _, ready := e.DOM.Find("a.Link--primary.Link--animatedHover.ArticleInList-readMoreLink").Eq(0).Attr("href"); !ready {
		c.visitSiteFinish <- struct{}{}
	}

	go func() {
		c.visitSiteReady <- struct{}{}
	}()

	for i := 0; ; i++ {
		u, ready := e.DOM.Find("a.Link--primary.Link--animatedHover.ArticleInList-readMoreLink").Eq(i).Attr("href")
		if !ready {
			return
		}
		url := site + u
		if url == c.oldIncr {
			err := c.finish()
			if err != nil {
				logger.Error("Error in the end:", err)
				c.errCh <- err
				return
			}
			c.visitSiteFinish <- struct{}{}
			return
		}
		c.blogURL <- &url
		<-c.crawlerFinish
	}
}
