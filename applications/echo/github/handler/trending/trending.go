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

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/crawler/github"
	"github.com/fengyfei/gu/applications/echo/core"
	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/models/github/trending"
)

type (
	createReq struct {
		Title    *string `json:"title" validate:"required,alphanum"`
		Abstract *string `json:"abstract"`
		Lang     *string `json:"lang" validate:"required,alphanum"`
		Stars    int     `json:"stars" validate:"required"`
		Today    int     `json:"today" validate:"required"`
	}

	// infoReq - The request struct that get one repos detail information.
	infoReq struct {
		ID string `json:"id" validate:"required,alphanum,len=24"`
	}

	langReq struct {
		Lang *string `json:"lang" validate:"required,alphanum"`
	}

	infoResp struct {
		Title    string
		Abstract string
		Lang     string
		Stars    int
		Today    int
	}
)

// LangInfo - Get library trending based on the language.
// If there is no data in the database, get the data from GitHub.
func LangInfo(c echo.Context) error {
	var (
		err        error
		req        langReq
		t          trending.Trending
		tStore     *trending.Trending
		tStoreList []*trending.Trending
		info       infoResp
		resp       []infoResp
		tInfo      *github.Trending
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
	loop:
		for {
			select {
			case <-github.DataPipe.Done:
				break loop
			case tInfo = <-github.DataPipe.DataCh:
				info = infoResp{
					Title:    tInfo.Title,
					Abstract: tInfo.Abstract,
					Lang:     tInfo.Lang,
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
		}
	}()

	if err = langCrawler(*req.Lang); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

	if err = trending.Service.CreateList(tStoreList); err != nil {
		return core.NewErrorWithMsg(http.StatusInternalServerError, err.Error())
	}

finish:
	return c.JSON(http.StatusOK, resp)
}

func langCrawler(tag string) error {
	c := github.NewTrendingCrawler(tag)

	err := crawler.StartCrawler(c)
	if err != nil {
		return err
	}

	return nil
}
