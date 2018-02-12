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
 *     Initial: 2018/02/12        Li Zebang
 */

package main

import (
	"fmt"
	"time"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/segment"
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
	var blogCh = make(chan *segment.Blog)

	c := segment.NewSegmentCrawler(blogCh)
	go func() {
		err := crawler.StartCrawler(c)
		logger.Error("Error in running the crawler:", err)
		return
	}()

	for {
		select {
		case blog := <-blogCh:
			err := store(blog)
			if err != nil {
				logger.Error("Error in storing the Segment blog:", err)
			}
			err = release(blog)
			if err != nil {
				logger.Error("Error in releasing the Segment blog to slack:", err)
			}
			logger.Info("Success the Segment blog:", blog.Date)
		case <-time.NewTimer(3 * time.Second).C:
			return
		}
	}
}

func store(blog *segment.Blog) error {
	c := session.DB("crawler").C("Segment blog")

	err := c.Insert(blog)
	session.Refresh()

	return err
}

func release(blog *segment.Blog) error {
	// If you don't have a custom bot, you can add one through
	// https://<your-workspace>.slack.com/services/new/bot
	cli := slack.NewClient("your custom bot token")

	text := fmt.Sprintf("source: Segment blog\ntitle: %s\ndate: %s\nurl: %s\n", blog.Title, blog.Date, blog.URL)

	err := cli.PostMessage("your channel", text)
	if err != nil {
		return err
	}

	return cli.UploadFile("your channel", text, "markdown", blog.Blog)
}
