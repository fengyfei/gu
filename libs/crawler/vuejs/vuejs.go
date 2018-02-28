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
 *     Initial: 2018/02/07        Li Zebang
 */

package vuejs

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asciimoo/colly"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	store "github.com/fengyfei/gu/libs/store/badger"
)

type vuejsCrawler struct {
	collector *colly.Collector

	dataCh   chan *crawler.Data
	finishCh chan struct{}

	ready          bool
	visitSiteReady chan struct{}
	crawlerFinish  chan struct{}

	db      *store.BadgerDB
	oldIncr string
	newIncr string
}

const (
	defaultOldIncr = "0"
	site           = "https://news.vuejs.org/issues/"
)

// NewVuejsCrawler generates a crawler for vuejs news.
func NewVuejsCrawler(dataCh chan *crawler.Data, finishCh chan struct{}) crawler.Crawler {
	return &vuejsCrawler{
		collector: colly.NewCollector(),
		dataCh:    dataCh,
		finishCh:  finishCh,

		visitSiteReady: make(chan struct{}),
		crawlerFinish:  make(chan struct{}),
	}
}

// Crawler interface Init
func (c *vuejsCrawler) Init() error {
	c.collector.OnHTML("article.issue", c.parse)

	err := c.prepare()
	if err != nil {
		logger.Error("Error in preparing to start:", err)
		return err
	}
	return nil
}

// Crawler interface Start
func (c *vuejsCrawler) Start() error {
	go func() {
		c.visitSiteReady <- struct{}{}
	}()

	incr, _ := strconv.Atoi(c.oldIncr)

	for {
		select {
		case <-c.visitSiteReady:
			go func() {
				err := c.collector.Visit(fmt.Sprintf("%s%d", site, incr))
				if err != nil {
					if c.ready && err.Error() == http.StatusText(http.StatusNotFound) {
						c.crawlerFinish <- struct{}{}
						return
					} else if err.Error() != http.StatusText(http.StatusNotFound) {
						logger.Error("Error in crawling the news:", err)
					}
				}

				if !c.ready {
					c.visitSiteReady <- struct{}{}
				}
				incr++
			}()
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

func (c *vuejsCrawler) prepare() error {
	db, err := store.NewBadgerDB(options.FileIO, "vuejs-news-badger", true)
	if err != nil {
		return err
	}
	c.db = db

	value, err := db.Get([]byte("increment-key"))
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

func (c *vuejsCrawler) finish() error {
	return c.db.Set([]byte("increment-key"), []byte(c.newIncr))
}

func (c *vuejsCrawler) parse(e *colly.HTMLElement) {
	c.ready = true
	c.newIncr = e.Request.URL.String()[len(site):]
	if c.oldIncr == c.newIncr {
		c.crawlerFinish <- struct{}{}
		return
	}

	data := crawler.Data{
		Source: "Vuejs News",
		Date:   e.DOM.Find("span.issue-date").Text(),
		Title:  e.DOM.Find("div.issue-title").Text(),
		URL:    e.Request.URL.String(),
		Text:   "Description: " + e.DOM.Find("p.issue-description").Text() + "\nContent:\n",
	}

	data.Text += "Story\n"
	for i := 0; ; i++ {
		s := e.DOM.Find("div.story").Eq(i)
		url, _ := s.Find("a").Attr("href")
		if url == "" {
			break
		}
		data.Text += fmt.Sprintf("%d. %s\n%s\n%s\n", i+1, s.Find("h1").Text(), url, s.Find("p").Text())
	}

	data.Text += "Library\n"
	for i := 0; ; i++ {
		s := e.DOM.Find("div.library").Eq(i)
		url, _ := s.Find("a").Attr("href")
		if url == "" {
			break
		}
		data.Text += fmt.Sprintf("%d. %s\n%s\n%s\n", i+1, s.Find("h1").Text(), url, s.Find("p").Text())
	}

	c.dataCh <- &data
	c.visitSiteReady <- struct{}{}
}
