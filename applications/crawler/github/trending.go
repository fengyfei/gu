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
 *     Initial: 2017/11/08        Feng Yifei
 */

package github

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/asciimoo/colly"
	"github.com/fengyfei/gu/libs/crawler"
)

type trendingCrawler struct {
	collector *colly.Collector
	topic     *string
}

// NewTrendingCrawler generates a crawler for github trending.
func NewTrendingCrawler(tag string) crawler.Crawler {
	return &trendingCrawler{
		collector: colly.NewCollector(),
		topic:     &tag,
	}
}

// Crawler interface Init
func (c *trendingCrawler) Init() error {
	c.collector.OnHTML("ol.repo-list", c.parse)
	return nil
}

// Crawler interface Start
func (c *trendingCrawler) Start() error {
	return c.collector.Visit("https://github.com/trending/" + *c.topic)
}

func (c *trendingCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Children().Each(c.parseContent)
}

func (c *trendingCrawler) parseContent(_ int, s *goquery.Selection) {
	fmt.Print(s.Children().Eq(0).Find("a").Attr("href"))
	fmt.Print("\t")
	fmt.Print(s.Children().Eq(2).Find("p").Text())
	fmt.Print("\t")
	fmt.Print(s.Children().Eq(3).Children().Eq(1).Text())
	fmt.Print("\t")
	fmt.Print(s.Children().Eq(3).Find("span.float-sm-right").Text())
	fmt.Println()
}
