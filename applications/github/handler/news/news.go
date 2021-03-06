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
 *     Initial: 2018/05/14        Tong Yuehong
 */

package news

import (
	newsapi "github.com/kaelanb/newsapi-go"

	"github.com/TechCatsLab/apix/http/server"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
)

type query struct {
	Query []string `json:"query"`
}

var (
	apikey = "cb49c9acbcb64a91b049f272b54b2554"
	client = newsapi.New(apikey)
)

func Everything(this *server.Context) error {
	var (
		query        query
		newsResponse *newsapi.NewsResponse
		err          error
	)

	if err := this.JSONBody(&query); err != nil {
		logger.Error("[newsapi][everything] parameters error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	for i := 0; i < 3; i++ {
		newsResponse, err = client.GetEverything(query.Query)
		if err != nil {
			logger.Error("[newsapi][everything] query error:", err)
			continue
		} else {
			break
		}
	}

	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrInternalServerError, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, newsResponse)
}
