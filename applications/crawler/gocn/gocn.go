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
 *     Initial: 2018/02/10        Li Zebang
 */

package main

import (
	"fmt"
	"time"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/gocn"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/social/slack"
	mgo "gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
)

func init() {
	var err error

	session, err = mgo.DialWithTimeout("mongodb://127.0.0.1:27017", time.Second)
	if err != nil {
		logger.Error("Conn't get connection to mongodb:", err)
	}

	session.SetMode(mgo.Monotonic, true)

	logger.Info("The mongoDB is connected!")
}

func main() {
	var (
		newsCh = make(chan *gocn.GoCN)
		endCh  = make(chan bool)
	)

	c := gocn.NewGoCNCrawler(newsCh, endCh)
	go func() {
		crawler.StartCrawler(c)
	}()

	for {
		select {
		case news := <-newsCh:
			err := store(news)
			if err != nil {
				logger.Error("Error in storing GoCN news:", err)
			}
			err = release(news)
			if err != nil {
				logger.Error("Error in releasing GoCN news to slack:", err)
			}
			logger.Info("Success GoCN news:", news.Date)
		case <-endCh:
			return
		}
	}
}

func store(news *gocn.GoCN) error {
	c := session.DB("crawler").C("GoCN每日新闻")

	err := c.Insert(news)
	session.Refresh()

	return err
}

func release(news *gocn.GoCN) error {
	// If you don't have a custom bot, you can add one through
	// https://<your-workspace>.slack.com/services/new/bot
	cli := slack.NewClient("your custom bot token")

	text := "source: GoCN每日新闻\n"
	text += fmt.Sprintf("date: %s\n", news.Date)
	text += fmt.Sprintf("url: %s\n", news.URL)
	for k, v := range news.Content {
		text += fmt.Sprintf("%s: %s\n", k, v)
	}

	return cli.PostMessage("your channel", text)
}
