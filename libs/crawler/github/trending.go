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
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/asciimoo/colly"
	"github.com/fengyfei/gu/libs/crawler"
)

// Trending used to store the acquired data.
type Trending struct {
	Title    string
	Abstract string
	Lang     string
	Date     string
	Stars    int
	Today    int
}

var (
	DataPipe chan *Trending = make(chan *Trending)
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
	rawTitle, _ := s.Children().Eq(0).Find("a").Attr("href")
	title := strings.TrimPrefix(rawTitle, "/")

	rawAbstract := s.Children().Eq(2).Find("p").Text()
	abstract := strings.TrimSpace(rawAbstract)

	rawStars := s.Children().Eq(3).Children().Eq(1).Text()
	trimStar := strings.TrimSpace(rawStars)
	stars := star2Int(trimStar)

	rawToday := s.Children().Eq(3).Find("span.float-sm-right").Text()
	trimToday := strings.TrimSpace(rawToday)
	today := today2Int(trimToday)

	date := time.Now().Format("20060102")
	info := &Trending{
		Title:    title,
		Abstract: abstract,
		Lang:     *c.topic,
		Date:     date,
		Stars:    stars,
		Today:    today,
	}

	DataPipe <- info

}

func star2Int(star string) int {
	var str string

	count := strings.Count(star, ",")

	switch count {
	case 0:
		str = star
	case 1:
		list := strings.Split(star, ",")
		str = list[0] + list[1]
	case 2:
		list := strings.Split(star, ",")
		str = list[0] + list[1] + list[2]
	}

	i, _ := strconv.ParseInt(str, 10, 0)
	return int(i)
}

func today2Int(today string) int {
	list := strings.Split(today, " ")
	str := list[0]
	i, _ := strconv.ParseInt(str, 10, 0)
	return int(i)
}
