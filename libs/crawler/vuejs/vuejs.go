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

	"github.com/asciimoo/colly"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/util/bolt"
)

type vuejsCrawler struct {
	collector *colly.Collector

	dataCh   chan crawler.Data
	finishCh chan struct{}

	db         *bolt.DB
	earlierURL string
	currentURL string
}

const (
	key      = "vuejs"
	site     = "https://news.vuejs.org/"
	previous = "a.issue-nav-link.issue-nav-link--next"
)

// NewVuejsCrawler generates a crawler for vuejs news.
func NewVuejsCrawler(dataCh chan crawler.Data, finishCh chan struct{}) crawler.Crawler {
	return &vuejsCrawler{
		collector: colly.NewCollector(),
		dataCh:    dataCh,
		finishCh:  finishCh,
	}
}

func (c *vuejsCrawler) prepare() error {
	db, err := bolt.Open(crawler.CrawlerPath)
	if err != nil {
		return err
	}
	c.db = db

	value, err := c.db.Get([]byte(crawler.CrawlerBucket), []byte(key))
	if err != nil {
		return err
	}
	c.earlierURL = string(value)
	return nil
}

func (c *vuejsCrawler) shutdown() error {
	fmt.Println("-=-=-=-=-=-=-=-shutdown=-=-=-=-=-=-=-=-")

	return c.db.Set([]byte(crawler.CrawlerBucket), []byte(key), []byte(c.currentURL))
}

// Crawler interface Init
func (c *vuejsCrawler) Init() error {
	c.collector.OnHTML("a.issue-nav-link.issue-nav-link--next", c.visitNext)
	c.collector.OnHTML("article.issue", c.parseNews)

	return c.prepare()
}

// Crawler interface Start
func (c *vuejsCrawler) Start() error {
	return c.collector.Visit(site)
}

func (c *vuejsCrawler) visitNext(e *colly.HTMLElement) {
	if e.Attr("href") == "" {
		c.shutdown()
		fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

		c.finishCh <- struct{}{}
		close(c.finishCh)
	}
	c.collector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
}

func (c *vuejsCrawler) parseNews(e *colly.HTMLElement) {
	if e.Request.URL.String() == site {
		c.currentURL, _ = e.DOM.Find("a").Eq(0).Attr("href")
	}
	subURL, _ := e.DOM.Find("a").Eq(0).Attr("href")
	if subURL == c.earlierURL {
		c.shutdown()
		c.finishCh <- struct{}{}
		close(c.finishCh)
	}

	data := crawler.DefaultData{
		Source: "Vuejs News",
		Date:   e.DOM.Find("span.issue-date").Text(),
		Title:  e.DOM.Find("div.issue-title").Text(),
		URL:    e.Request.AbsoluteURL(subURL),
		Text:   "Description: " + e.DOM.Find("div.issue-description").Text() + "\nContent:\n",
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
}
