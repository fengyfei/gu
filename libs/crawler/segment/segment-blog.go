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

type ok struct{}

type segmentCrawler struct {
	collector  *colly.Collector
	errCh      chan error
	urlCh      chan *string
	urlReady   chan ok
	urlFinish  chan ok
	blogFinish chan ok
	blogCh     chan *Blog
	db         *store.BadgerDB
	oldIncr    string
	newIncr    string
}

const (
	defaultOldIncr = "default-old-increment"
	site           = "https://segment.com"
)

// NewSegmentCrawler generates a crawler for Segment blogs.
func NewSegmentCrawler(ch chan *Blog) crawler.Crawler {
	return &segmentCrawler{
		collector:  colly.NewCollector(),
		urlCh:      make(chan *string),
		urlReady:   make(chan ok),
		urlFinish:  make(chan ok),
		blogFinish: make(chan ok),
		blogCh:     ch,
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
	var (
		page int
		url  string
	)

	go func() {
		c.urlReady <- ok{}
	}()

	go c.startBlog()

	for {
		select {
		case <-c.urlFinish:
			return nil
		case <-c.urlReady:
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
					return nil
				}
				logger.Error("Error in getting blog url", err)
				return err
			}
		case err := <-c.errCh:
			return err
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
		c.urlFinish <- ok{}
	}

	go func() {
		c.urlReady <- ok{}
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
			c.urlFinish <- ok{}
			return
		}
		c.urlCh <- &url
		<-c.blogFinish
	}
}
