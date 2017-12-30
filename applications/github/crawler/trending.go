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
	// Subscriber represents the subscriber of subject SubjectTrending.
	Subscriber    *nc.Subscriber
	TrendingCache *trendingCache = newCache()
)

func init() {
	subAndStartCrawler()
}

// subAndStartCrawler subscribe to SubjectTrending and begin to execute the
// crawler after receiving the specified language.
func subAndStartCrawler() {
	var err error

	msgHandler := func(msg *gonc.Msg) {
		go storeTrending()
		startLangCrawler(string(msg.Data))
	}

	Subscriber, err = nats.Conn.Subscribe(nats.SubjectTrending, msgHandler)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	logger.Debug("Successfully subscribe to NATS subject: SubjectTrending.")
}

func storeTrending() {
	var (
		err   error
		tInfo *github.Trending
		tList []*github.Trending
	)

	for {
		select {
		case tInfo = <-github.DataPipe:
			// write to cache
			tList = append(tList, tInfo)
			TrendingCache.Store(tInfo.Lang, tList)

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
		default:
		}
	}

}

func startLangCrawler(lang string) error {
	c := github.NewTrendingCrawler(lang)

	return crawler.StartCrawler(c)
}
