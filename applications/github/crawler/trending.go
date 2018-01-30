/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 *     Initial: 2017/12/29        Jia Chenhui
 */

package crawler

import (
	"github.com/fengyfei/gu/applications/nats"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/crawler/github"
	"github.com/fengyfei/gu/libs/logger"
	nc "github.com/fengyfei/gu/libs/nats"
	"github.com/fengyfei/gu/models/github/trending"
	gonc "github.com/nats-io/go-nats"
)

var (
	// TrendingCache use to store the trending data of specified language.
	TrendingCache *trendingCache
)

func init() {
	TrendingCache = newTrendingCache()
}

// SubNatsWithSubject subscribe to subject and begin to execute the
// crawler after receiving the specified language.
func SubNatsWithSubject(subject *string) (*nc.Subscriber, error) {
	var (
		err        error
		subscriber *nc.Subscriber
	)

	dataPipe := make(chan *github.Trending)
	msgHandler := func(msg *gonc.Msg) {
		go storeTrending(dataPipe)
		startTrendingCrawler(string(msg.Data), dataPipe)
	}

	subscriber, err = nats.Conn.Subscribe(*subject, msgHandler)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

// storeTrending store it to TrendingCache and MongoDB when the trending
// data is received.
func storeTrending(dataPipe chan *github.Trending) {
	var (
		err   error
		tInfo *github.Trending
	)

	for {
		select {
		case tInfo = <-dataPipe:
			// write to cache
			TrendingCache.Write(&tInfo.Lang, tInfo)

			// write to database
			t := &trending.Trending{
				Title:    tInfo.Title,
				Abstract: tInfo.Abstract,
				Lang:     tInfo.Lang,
				Date:     tInfo.Date,
				Stars:    tInfo.Stars,
				Today:    tInfo.Today,
			}
			err = trending.Service.Create(t)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// startTrendingCrawler start the crawler.
func startTrendingCrawler(lang string, dataPipe chan *github.Trending) error {
	c := github.NewTrendingCrawler(lang, dataPipe)

	return crawler.StartCrawler(c)
}
