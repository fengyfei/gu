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
 *     Initial: 2018/02/25        Li Zebang
 */

package client

import (
	"fmt"
	"time"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/social/slack"
)

type Client struct {
	Crawler   crawler.Crawler
	DataCh    *chan *crawler.Data
	FinishCh  *chan struct{}
	DB        string
	C         string
	BotsToken string
	Channel   string
}

var Clients = make(map[string]*Client)

func Start(client string) error {
	var (
		cli   = Clients[client]
		errCh = make(chan error)
	)

	go func() {
		err := crawler.StartCrawler(cli.Crawler)
		if err != nil {
			logger.Error("Error in running the crawler:", err)
			errCh <- err
			return
		}
	}()

	for {
		select {
		case data := <-*cli.DataCh:
			err := cli.store(data)
			if err != nil {
				logger.Error(fmt.Sprintf("Error in storing the %s: %s", data.Source, err.Error()))
				return err
			}
			err = cli.release(data)
			if err != nil {
				logger.Error("Error in releasing to slack:", err)
				return err
			}
			logger.Info(fmt.Sprintf("Success the %s: %s", data.Source, data.Date))
		case <-*cli.FinishCh:
			time.Sleep(5 * time.Second)
			return nil
		case err := <-errCh:
			return err
		}
	}
}

func (c *Client) store(data *crawler.Data) error {
	db := session.DB(c.DB).C(c.C)

	err := db.Insert(data)
	session.Refresh()

	return err
}

func (c *Client) release(data *crawler.Data) error {
	cli := slack.NewClient(c.BotsToken)

	if data.FileType != "" {
		text := fmt.Sprintf("Source: %s\nDate: %s\nTitle: %s\nURL: %s\n", data.Source, data.Date, data.Title, data.URL)

		err := cli.PostMessage(c.Channel, text)
		if err != nil {
			return err
		}

		return cli.UploadFile(c.Channel, data.Title, data.FileType, data.Text)
	}

	text := fmt.Sprintf("Source: %s\nDate: %s\nTitle: %s\nURL: %s\n%s", data.Source, data.Date, data.Title, data.URL, data.Text)

	return cli.PostMessage(c.Channel, text)
}
