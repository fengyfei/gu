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

type ok struct{}

type vuejsCrawler struct {
	collector  *colly.Collector
	ready      bool
	readyCh    chan ok
	newsFinish chan ok
	newsCh     chan *News
	endCh      chan bool
	db         *store.BadgerDB
	oldIncr    string
	newIncr    string
}

type News struct {
	Title       string
	Date        string
	Description string
	URL         string
	Content     map[string][]Content
}

type Content struct {
	Class       string
	URL         string
	Title       string
	Description string
}

const (
	defaultOldIncr = "0"
	site           = "https://news.vuejs.org/issues/"
)

// NewVuejsCrawler generates a crawler for vuejs news.
func NewVuejsCrawler(ch chan *News, end chan bool) crawler.Crawler {
	return &vuejsCrawler{
		collector:  colly.NewCollector(),
		readyCh:    make(chan ok),
		newsFinish: make(chan ok),
		newsCh:     ch,
		endCh:      end,
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
	var (
		incr int
		err  error
	)

	if c.oldIncr != "" {
		incr, err = strconv.Atoi(c.oldIncr)
		if err != nil {
			logger.Error("Error in converting string to int:", err)
			return err
		}
	}

	go func() {
		c.readyCh <- ok{}
	}()

	for {
		select {
		case <-c.readyCh:
			go func() {
				err := c.collector.Visit(fmt.Sprintf("%s%d", site, incr))
				if err != nil {
					if c.ready && err.Error() == http.StatusText(http.StatusNotFound) {
						c.newsFinish <- ok{}
						return
					} else if err.Error() != http.StatusText(http.StatusNotFound) {
						logger.Error("Error in crawling the news:", err)
					}
				}

				if !c.ready {
					c.readyCh <- ok{}
				}
				incr++
			}()
		case <-c.newsFinish:
			err := c.finish()
			if err != nil {
				logger.Error("Error in the end:", err)
			}
			c.endCh <- true
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
	c.newIncr = e.Request.URL.String()[30:]
	if c.oldIncr == c.newIncr {
		c.newsFinish <- ok{}
		return
	}
	news := News{
		Title:       e.DOM.Find("div.issue-title").Text(),
		Date:        e.DOM.Find("span.issue-date").Text(),
		Description: e.DOM.Find("p.issue-description").Text(),
		URL:         e.Request.URL.String(),
		Content:     make(map[string][]Content),
	}
	for i := 0; ; i++ {
		s := e.DOM.Find("div").Eq(i)
		class, _ := s.Attr("class")
		if class == "story" || class == "library" {
			url, _ := s.Find("a").Attr("href")
			content := Content{
				Class:       class,
				URL:         url,
				Title:       s.Find("h1").Text(),
				Description: s.Find("p").Text(),
			}
			news.Content[class] = append(news.Content[class], content)
		} else if class == "" {
			break
		}
	}
	c.newsCh <- &news
	c.readyCh <- ok{}
}
