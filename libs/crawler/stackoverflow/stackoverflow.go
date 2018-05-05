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
 *     Initial: 2018/4/23        Zhang Hao
 */

package stackoverflow

import (
	"fmt"
	"strings"
	"github.com/fengyfei/gu/libs/crawler/util/bolt"
	"github.com/gocolly/colly"
	"github.com/fengyfei/gu/libs/crawler"
)

type Blog struct {
	Source  string `bson:"source"`
	Title   string `bson:"title"`
	Author  string `bson:"author"`
	Date    string `bson:"date"`
	Photo   string `bson:"photo"`
	Content string `bson:"content"`
}

type stackOverFlowCrawler struct {
	collector       *colly.Collector
	detailCollector *colly.Collector
	db              *bolt.DB
	newestUrl       string
	lastUrl         string
	counter         int64
	isAllDateGet    bool
	dataCh          chan crawler.Data
	finishCh        chan struct{}
}

func NewStackOverFlow(dataCh chan crawler.Data, finishCh chan struct{}) *stackOverFlowCrawler {
	boltDB, err := initBoltDB()
	if err != nil {
		panic(err)
	}
	return &stackOverFlowCrawler{
		collector:       colly.NewCollector(),
		detailCollector: colly.NewCollector(),
		db:              boltDB,
		counter:         0,
		isAllDateGet:    false,
		dataCh:          dataCh,
		finishCh:        finishCh,
	}
}

func initBoltDB() (*bolt.DB, error) {
	return bolt.Open(crawler.CrawlerPath)
}

func (sc *stackOverFlowCrawler) closeBoltDB() {
	err := sc.db.Close()
	if err != nil {
		defer close(sc.dataCh)
		defer close(sc.finishCh)
		panic(err)
	}
}

func (sc *stackOverFlowCrawler) preUpdate() error {
	lastUrlSlice, err := sc.db.Get([]byte(crawler.CrawlerBucket), []byte("stackoverflowLastUrl"))
	if err != nil {
		return err
	}
	if lastUrlSlice == nil {
		sc.lastUrl = ""
		fmt.Println("*** Starting to crawl for the first time. ***")
	} else {
		sc.lastUrl = string(lastUrlSlice)
	}
	return nil
}

func (sc *stackOverFlowCrawler) putLastUrl() error {
	err := sc.db.Set([]byte(crawler.CrawlerBucket), []byte("stackoverflowLastUrl"), []byte(sc.newestUrl))
	return err
}

func (sc *stackOverFlowCrawler) visit(url string) error {
	err := sc.collector.Visit(url)
	if err != nil {
		return err
	}
	return nil
}

func (sc *stackOverFlowCrawler) onRequest() {
	sc.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
}

func (sc *stackOverFlowCrawler) onHtml() {
	sc.collector.OnHTML("article div.m-post-card__content-column", sc.parse)
}

func (sc *stackOverFlowCrawler) detailOnHtml() {
	sc.detailCollector.OnHTML("main#main.site-main", sc.parseDetail)
}

func (sc *stackOverFlowCrawler) parse(e *colly.HTMLElement) {
	if sc.counter == 0 {
		sc.newestUrl, _ = e.DOM.Find("h2 a").Attr("href")
	}
	link, _ := e.DOM.Find("h2 a").Attr("href")
	if link != sc.lastUrl && (sc.isAllDateGet == false) {
		sc.detailCollector.Visit(link)
	} else {
		sc.isAllDateGet = true
	}
	sc.counter++
}

func (sc *stackOverFlowCrawler) parseDetail(e *colly.HTMLElement) {
	var Blog = &Blog{}
	Blog.Source = "stackOverFlow blog"
	Blog.Title = strings.TrimSpace(e.DOM.Find("div.column h1.section-title").Text())
	Blog.Photo, _ = e.DOM.Find("div span span a img.avatar__image").Attr("src")
	Blog.Author = e.DOM.Find("div.m-post__meta span.author-name a").Text()
	Blog.Date = e.DOM.Find("div.m-post__meta span.date time.entry-date").Text()
	Blog.Content = e.DOM.Find("div.m-post-content").Text()
	sc.dataCh <- Blog
}

func (b *Blog) String() string {
	return fmt.Sprintf("Source: %s\nTitle: %s\nAuthor: %s\nDate: %s\nPhoto: %s\nContent: %s\n", b.Source, b.Title, b.Author, b.Date, b.Photo, b.Content)
}

func (b *Blog) File() (title, filetype, content string) {
	return "", "", ""
}

func (b *Blog) IsFile() bool {
	return false
}
