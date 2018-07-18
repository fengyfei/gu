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
 *     Initial: 2017/12/28        Jia Chenhui
 */

package trending

import (
	"time"

	"github.com/TechCatsLab/apix/http/server"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/github/crawler"
	"github.com/fengyfei/gu/applications/nats"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/crawler/github"
	"github.com/fengyfei/gu/libs/logger"
	nc "github.com/fengyfei/gu/libs/nats"
)

type (
	// langReq - The request struct that get the trending of the day of a language.
	langReq struct {
		Lang *string `json:"lang" validate:"required"`
	}

	// infoResp - The response struct that represents the trending of the day of a language.
	infoResp struct {
		Title    string `json:"title"`
		Abstract string `json:"abstract"`
		Lang     string `json:"lang"`
		Date     string `json:"date"`
		Stars    int    `json:"stars"`
		Today    int    `json:"today"`
	}
)

// LangInfo - Get library trending based on the language.
// If there is no data in cache, get data from GitHub.
func LangInfo(c *server.Context) error {
	var (
		err   error
		ok    bool
		req   langReq
		info  infoResp
		sub   *nc.Subscriber
		t     github.Trending
		tList []github.Trending
		resp  = make([]infoResp, 0)
		today = time.Now().Format("20060102")
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if sub, err = crawler.SubNatsWithSubject(req.Lang); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrSubNats, nil)
	}
	defer sub.Unsubscribe()

readFromCache:
	if tList, ok = crawler.TrendingCache.Read(req.Lang); !ok {
		nats.StartTrendingCrawler(req.Lang)
		time.Sleep(5 * time.Second)

		goto readFromCache
	}

	for _, t = range tList {
		info = infoResp{
			Title:    t.Title,
			Abstract: t.Abstract,
			Lang:     t.Lang,
			Date:     t.Date,
			Stars:    t.Stars,
			Today:    t.Today,
		}

		resp = append(resp, info)
	}

	// Check whether the data in the cache is up-to-date.
	if t.Date != today {
		crawler.TrendingCache.Flush(req.Lang)
		nats.StartTrendingCrawler(req.Lang)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
