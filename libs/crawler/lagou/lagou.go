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
 *     Initial: 2018/01/24        Shi Ruitao
 */

package lagou

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/fengyfei/gu/libs/crawler"
	"github.com/fengyfei/gu/libs/logger"
)

var DataPipe chan Result = make(chan Result)

const (
	jobHtml = "https://www.lagou.com/jobs/%d.html"

	Method     = "POST"
	RequestUrl = "https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false&isSchoolJob=0"
	Referer    = "https://www.lagou.com/jobs/list_java?city=%E5%85%A8%E5%9B%BD&cl=false&fromSearch=true&labelWords=&suginput="
	Cookie     = "JSESSIONID=ABAAABAACDBABJBF293E69FB60C42BE49201229C99D8FCD; user_trace_token=20171205104807-b96aac92-d966-11e7-9c13-5254005c3644; LGUID=20171205104807-b96aaf6b-d966-11e7-9c13-5254005c3644; index_location_city=%E5%85%A8%E5%9B%BD; _putrc=68D38A1FF9900849; login=true; unick=%E6%97%B6%E7%91%9E%E6%B6%9B; _ga=GA1.2.206867280.1512442088; _gid=GA1.2.136452891.1514977679; Hm_lvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1512442089,1512443109,1514977679; LGSID=20180104095511-4cb7a44f-f0f2-11e7-bbf1-525400f775ce; _gat=1; TG-TRACK-CODE=index_search; Hm_lpvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1515033736; LGRID=20180104104216-e02d395f-f0f8-11e7-bc05-525400f775ce; SEARCH_ID=5a8f1728749c46b299033be73e954568"
	UserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"
)

// What kind of job do you want to crawl?
type lagouClient struct {
	job *string
}

// NewLaGouCrawler generates a crawler for lagou.
func NewLaGouCrawler(job string) crawler.Crawler {
	return &lagouClient{
		job: &job,
	}
}

// Crawler interface Init
func (lc *lagouClient) Init() error {
	return nil
}

// Start interface Start
func (lc *lagouClient) Start() error {
	var sum int
	for i := 1; true; i++ {
		urls, _ := lc.getJobHttp(i)
		if len(urls) == 0 {
			break
		}
		sum += len(urls)
	}
	fmt.Println("sum:", sum)
	return nil
}

func (lc *lagouClient) getJobHttp(pn int) ([]string, error) {
	hc := &http.Client{}

	data := url.Values{}
	data.Set("fires", "fasle")
	data.Set("pn", fmt.Sprintf("%d", pn))
	data.Set("kd", *lc.job)
	b := strings.NewReader(data.Encode())

	req := newRequest(Method, RequestUrl, b)
	lc.setHeader(req, Referer, Cookie, RequestUrl)

	resp, err := hc.Do(req)
	if err != nil {
		logger.Error("error in doing a http request.", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error in reading response body.", err)
		return nil, err
	}

	var object Object
	err = json.Unmarshal(body, &object)
	if err != nil {
		logger.Error("error in unmarshalling response body.", err)
		return nil, err
	}

	var urls = make([]string, 0)
	for _, value := range object.Content.PositionResult.Result {
		urls = append(urls, fmt.Sprintf(jobHtml, value.PositionId))
	}

	go chRoutine()

	for i := 0; i < len(urls); i++ {
		DataPipe <- object.Content.PositionResult.Result[i]
	}
	return urls, nil
}

func chRoutine() {
	for {
		fmt.Println(<-DataPipe)
	}
}

func newRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln("error in creating a new request.", err)
	}
	return req
}

func (lc *lagouClient) setHeader(req *http.Request, referer, cookie, user_agent string) {
	var headers = make(map[string]string)

	for k, v := range defaultHeader() {
		headers[k] = v
	}

	headers["Referer"] = referer
	headers["Cookie"] = cookie
	headers["User-Agent"] = user_agent

	header(req, headers)
}

func defaultHeader() map[string]string {
	var defaultHeaders = make(map[string]string)

	defaultHeaders["Accept"] = "application/json, text/javascript, */*; q=0.01"
	defaultHeaders["Accept-Encoding"] = "gzip, deflate, br"
	defaultHeaders["Accept-Language"] = "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7"
	defaultHeaders["Connection"] = "keep-alive"
	defaultHeaders["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	defaultHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"

	return defaultHeaders
}

func header(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
