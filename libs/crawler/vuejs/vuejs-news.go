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
	"github.com/asciimoo/colly"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	store "github.com/fengyfei/gu/libs/store/badger"
)

type vuejsCrawler struct {
	collector *colly.Collector
	increment string
	newsCh    chan *News
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
	noneBadger = "none badger"
)

// NewVuejsCrawler generates a crawler for vuejs news.
func NewVuejsCrawler(ch chan *News) crawler.Crawler {
	return &vuejsCrawler{
		collector: colly.NewCollector(),
		newsCh:    ch,
	}
}

// Crawler interface Init
func (c *vuejsCrawler) Init() error {
	c.collector.OnHTML("article.issue", c.parse)
	return nil
}

// Crawler interface Start
func (c *vuejsCrawler) Start() error {
	c.collector.Visit("https://news.vuejs.org/issues/73")
	return nil
}

func prepare() (string, error) {
	db, err := store.NewBadgerDB(options.FileIO, "./vuejs-news-badger", true)
	if err != nil {
		logger.Error("Error in opening badger database:", err)
		return "", err
	}
	value, err := db.Get([]byte("increment-key"))
	if len(value) != 0 && err == nil {
		return string(value), nil
	}
	if err != badger.ErrKeyNotFound {
		logger.Error("Error in getting the badger key:", err)
		return "", err
	}
	return noneBadger, nil
}

func (c *vuejsCrawler) parse(e *colly.HTMLElement) {
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
}
