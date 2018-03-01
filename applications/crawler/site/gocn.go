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
 *     Initial: 2018/03/01        Li Zebang
 */

package site

import (
	"github.com/fengyfei/gu/applications/crawler/client"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/gocn"
)

const (
	GoCNDailyNews = "gocn daily news"
)

func init() {
	client.Clients[GoCNDailyNews] = NewGoCNClinet
}

func NewGoCNClinet() *client.Client {
	var (
		dataCh   = make(chan *crawler.Data)
		finishCh = make(chan struct{})
	)

	crawler := gocn.NewGoCNCrawler(dataCh, finishCh)

	return &client.Client{
		Crawler:   crawler,
		DataCh:    &dataCh,
		FinishCh:  &finishCh,
		DB:        "Crawler",
		C:         "GoCN Daily News",
		BotsToken: "xoxb-312476598064-97wqE4OJeqhv4mTX1g2c9LZs",
		Channel:   "C97LN9DGF",
	}
}
