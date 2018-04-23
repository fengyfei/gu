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

package xteam

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/util/bolt"
	"github.com/gocolly/colly"
)

type Blog struct {
	Source  string   `bson:"source"`
	Title   string   `bson:"title"`
	Author  string   `bson:"author"`
	Date    string   `bson:"date"`
	Photo   string   `bson:"photo"`
	Content string   `bson:"content"`
	Tags    []string `bson:"tags"`
}

type xteamCrawler struct {
	collector       *colly.Collector
	detailCollector *colly.Collector
	db              *bolt.DB
	lastUrl         string
	newestUrl       string
	counter         int64
	isAllDateGet    bool
	dataCh          chan crawler.Data
	finishCh        chan struct{}
}

func NewXteam(dataCh chan crawler.Data, finishCh chan struct{}) *xteamCrawler {
	boltDB, err := initBoltDB()
	if err != nil {
		panic(err)
	}

	return &xteamCrawler{
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

func (xc *xteamCrawler) closeBoltDB() {
	err := xc.db.Close()
	if err != nil {
		defer close(xc.dataCh)
		defer close(xc.finishCh)
		panic(err)
	}
}

func (xc *xteamCrawler) prepare() error {
	lastUrlSlice, err := xc.db.Get([]byte(crawler.CrawlerBucket), []byte("xteamLastUrl"))
	if err != nil {
		return err
	}
	if lastUrlSlice == nil {
		xc.lastUrl = ""
		fmt.Println("*** Starting to crawl for the first time. ***")
	} else {
		xc.lastUrl = string(lastUrlSlice)
	}
	return nil
}

func (xc *xteamCrawler) putLastUrl() error {
	err := xc.db.Set([]byte(crawler.CrawlerBucket), []byte("xteamLastUrl"), []byte(xc.newestUrl))
	return err
}

func (xc *xteamCrawler) visit(url string) error {
	err := xc.collector.Visit(url)
	if err != nil {
		return err
	}
	return nil
}

func (xc *xteamCrawler) onRequest() {
	xc.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
}

func (xc *xteamCrawler) onHtml() {
	xc.collector.OnHTML("main article h2 a", xc.parse)
}

func (xc *xteamCrawler) detailOnHtml() {
	xc.detailCollector.OnHTML("main article", xc.parseDetail)
}

func (xc *xteamCrawler) parse(e *colly.HTMLElement) {
	if xc.counter == 0 {
		xc.newestUrl = e.Attr("href")
	}
	link := e.Attr("href")
	if link != xc.lastUrl && (xc.isAllDateGet == false) {
		xc.detailCollector.Visit("https://x-team.com" + link)
	} else {
		xc.isAllDateGet = true
	}
	xc.counter++
}

func (xc *xteamCrawler) parseDetail(e *colly.HTMLElement) {
	var Blog = &Blog{}
	Blog.Source = "x-team blog"
	Blog.Title = strings.TrimSpace(e.DOM.Find("h1.wrapper-m.title.post-title").Text())
	Blog.Photo, _ = e.DOM.Find("img.post-author-avatar").Attr("src")
	Blog.Author = e.DOM.Find("ul li.post-author-name span[itemprop]").Text()
	Blog.Date = e.DOM.Find("ul li.post-date span").Text()
	Blog.Content = e.DOM.Find("section div.kg-card-markdown").Text()
	e.DOM.Find("ul.button-action li ul.option-list li a[title]").Each(func(i int, selection *goquery.Selection) {
		tag, _ := selection.Attr("title")
		Blog.Tags = append(Blog.Tags, tag)
	})
	xc.dataCh <- Blog
}

func (b *Blog) String() string {
	return fmt.Sprintf("Source: %s\nTitle: %s\nAuthor: %s\nDate: %s\nPhoto: %s\nTags: %s\nContent: %s\n", b.Source, b.Title, b.Author, b.Date, b.Photo, b.Tags, b.Content)
}

func (b *Blog) File() (title, filetype, content string) {
	return "", "", ""
}

func (b *Blog) IsFile() bool {
	return false
}
