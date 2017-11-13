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
 *     Initial: 2017/11/13        Feng Yifei
 */

package google

import (
	"fmt"

	"github.com/asciimoo/colly"
	"github.com/fengyfei/gu/libs/crawler"
)

type golangBlogCrawler struct {
	collector *colly.Collector
}

// NewGolangBlogCrawler generates a crawler for github trending.
func NewGolangBlogCrawler() crawler.Crawler {
	return &golangBlogCrawler{
		collector: colly.NewCollector(),
	}
}

// Crawler interface Init
func (c *golangBlogCrawler) Init() error {
	c.collector.OnHTML("p.blogtitle", c.parse)
	return nil
}

// Crawler interface Start
func (c *golangBlogCrawler) Start() error {
	return c.collector.Visit("https://blog.golang.org/index")
}

func (c *golangBlogCrawler) parse(e *colly.HTMLElement) {
	fmt.Print(e.DOM.Find("a").Attr("href"))
	fmt.Print("\t")
	fmt.Print(e.DOM.Find("a").Text())
	fmt.Print("\t")
	fmt.Print(e.DOM.Find("span").Text())
	fmt.Println()
}
