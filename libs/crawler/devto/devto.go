/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/10/28        Feng Yifei
 */

package devto

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asciimoo/colly"

	"github.com/fengyfei/gu/libs/crawler"
)

type devToCrawler struct {
	collector *colly.Collector
	topic     *string
}

// NewDevToCrawler generates a crawler for dev.to.
func NewDevToCrawler(topic string) crawler.Crawler {
	return &devToCrawler{
		collector: colly.NewCollector(),
		topic:     &topic,
	}
}

// Crawler interface Init
func (c *devToCrawler) Init() error {
	c.collector.OnHTML("div.single-article", c.parse)
	return nil
}

// Crawler interface Start
func (c *devToCrawler) Start() error {
	return c.collector.Visit("https://dev.to/t/" + *c.topic)
}

func (c *devToCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Each(c.parseContent)
}

func (c *devToCrawler) parseContent(_ int, s *goquery.Selection) {
	cls, _ := s.Attr("class")

	if strings.Contains(cls, "single-article-small-pic") {
		fmt.Println(s.Children().Eq(1).Attr("href"))
		fmt.Println(s.Find("h3").Text())
	}
}
