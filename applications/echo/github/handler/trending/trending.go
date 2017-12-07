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
 *     Initial: 2017/11/24        Jia Chenhui
 */

package trending

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/crawler/github"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/models/github/trending"
)

const (
	pageSize = 25
)

type (
	// langReq - The request struct that get the trending of the day of a language.
	langReq struct {
		Lang *string `json:"lang" validate:"required,alpha"`
	}

	// infoResp - The response struct that represents the trending of the day of a language.
	infoResp struct {
		Title    string `json:"title"`
		Abstract string `json:"abstract"`
		Lang     string `json:"lang"`
		Stars    int    `json:"stars"`
		Today    int    `json:"today"`
	}
)

// LangInfo - Get library trending based on the language.
// If there is no data in the database, get the data from GitHub.
func LangInfo(c echo.Context) error {
	var (
		err        error
		req        langReq
		resp       []infoResp = make([]infoResp, 0)
		t          trending.Trending
		tStore     *trending.Trending
		tStoreList []*trending.Trending
		info       infoResp
		tInfo      *github.Trending
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)

	if err = c.Bind(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&req); err != nil {
		return core.NewErrorWithMsg(http.StatusBadRequest, err.Error())
	}

	tlist, err := trending.Service.GetByLang(req.Lang)

	if err != nil {
		if err == mgo.ErrNotFound {
			goto crawler
		}

		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	} else if tlist == nil {

		goto crawler
	} else {
		for _, t = range tlist {
			info = infoResp{
				Title:    t.Title,
				Abstract: t.Abstract,
				Lang:     *req.Lang,
				Stars:    t.Stars,
				Today:    t.Today,
			}

			resp = append(resp, info)
		}

		goto finish
	}

crawler:
	go func() {
		for i := 0; i < pageSize; i++ {
			go func() {
				wg.Add(1)
				select {
				case tInfo = <-github.DataPipe:
					info = infoResp{
						Title:    tInfo.Title,
						Abstract: tInfo.Abstract,
						Lang:     *req.Lang,
						Stars:    tInfo.Stars,
						Today:    tInfo.Today,
					}
					resp = append(resp, info)

					tStore = &trending.Trending{
						Title:    tInfo.Title,
						Abstract: tInfo.Abstract,
						Lang:     tInfo.Lang,
						Stars:    tInfo.Stars,
						Today:    tInfo.Today,
					}
					tStoreList = append(tStoreList, tStore)
				}
				wg.Done()
				return
			}()
		}
	}()

	if err = startLangCrawler(*req.Lang); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	wg.Wait()

	if err = trending.Service.CreateList(tStoreList); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

finish:
	return c.JSON(http.StatusOK, map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   resp,
	})
}

func startLangCrawler(tag string) error {
	c := github.NewTrendingCrawler(tag)

	return crawler.StartCrawler(c)
}
