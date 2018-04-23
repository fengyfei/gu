/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2018/4/23        Zhang Hao
 */

package xteam

import "fmt"

const (
	site = "https://x-team.com/blog/page/%d/"
)

func (xc *xteamCrawler) Init() error {
	err := xc.prepare()
	if err != nil {
		defer close(xc.dataCh)
		defer close(xc.finishCh)
		return err
	}
	xc.onRequest()
	xc.onHtml()
	xc.detailOnHtml()
	return nil
}

func (xc *xteamCrawler) Start() error {
	defer close(xc.dataCh)
	defer close(xc.finishCh)
	defer xc.closeBoltDB()
	for pageNumber := 1; !xc.isAllDateGet; pageNumber++ {
		url := fmt.Sprintf(site, pageNumber)
		err := xc.visit(url)
		if err != nil {
			if err.Error() == "Not Found" {
				break
			} else {
				return err
			}
		}
	}
	err := xc.putLastUrl()
	return err
}
