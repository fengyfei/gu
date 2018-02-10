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
	"os"

	"github.com/asciimoo/colly"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
)

type segmentCrawler struct {
	collector *colly.Collector
	urlCh     chan *string
	urlOver   chan done
	blogOver  chan done
	errCh     chan error
}

type done struct{}

const (
	site = "https://segment.com"
)

// NewSegmentCrawler generates a crawler for Segment blogs.
func NewSegmentCrawler() crawler.Crawler {
	return &segmentCrawler{
		collector: colly.NewCollector(),
		urlCh:     make(chan *string),
		urlOver:   make(chan done),
		blogOver:  make(chan done),
	}
}

// Crawler interface Init
func (c *segmentCrawler) Init() error {
	c.collector.OnHTML("a.Link--primary.Link--animatedHover.ArticleInList-readMoreLink", c.parseURL)
	return os.MkdirAll("blog", 0755)
}

// Crawler interface Start
func (c *segmentCrawler) Start() error {
	go c.startBlog()
	go c.startURL()

	for {
		if err, ok := <-c.errCh; ok {
			return err
		}
	}
}

func (c *segmentCrawler) parseURL(e *colly.HTMLElement) {
	url := site + e.Attr("href")
	c.urlCh <- &url
	<-c.blogOver
}

func (c *segmentCrawler) startURL() {
	var (
		page int
		url  string
	)

	for {
		page++
		if page == 1 {
			url = site + "/blog/"
		} else {
			url = fmt.Sprintf("%s/blog/page/%d", site, page)
		}
		err := c.collector.Visit(url)
		if err != nil {
			logger.Error("error in getting blog url", err)
			c.errCh <- err
		}
	}
}
