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
 *     Initial: 2018/01/26        Shi Ruitao
 */

package lagou

type Object struct {
	Content Content `json:"content"`
}

type Content struct {
	PositionResult PositionResult `json:"positionResult"`
}

type PositionResult struct {
	Result []Result `json:"result"`
}

// LagouResult used to store the  data.
type Result struct {
	CompanyId             int         `json:"companyId"`
	CompanyShortName      string      `json:"companyShortName"`
	CreateTime            string      `json:"createTime"`
	PositionId            int         `json:"positionId"`
	Score                 int         `json:"score"`
	PositionAdvantage     string      `json:"positionAdvantage"`
	Salary                string      `json:"salary"`
	WorkYear              string      `json:"workYear"`
	Education             string      `json:"education"`
	City                  string      `json:"city"`
	PositionName          string      `json:"positionName"`
	CompanyLogo           string      `json:"companyLogo"`
	FinanceStage          string      `json:"financeStage"`
	IndustryField         string      `json:"industryField"`
	JobNature             string      `json:"jobNature"`
	CompanySize           string      `json:"companySize"`
	Approve               int         `json:"approve"`
	CompanyLabelList      []string    `json:"companyLabelList"`
	PublisherId           int         `json:"publisherId"`
	District              interface{} `json:"district"`
	PositionLabels        []string    `json:"positionLables"`
	IndustryLabels        []string    `json:"industryLables"`
	BusinessZones         interface{} `json:"businessZones"`
	AdWord                int         `json:"adWord"`
	Longitude             string      `json:"longitude"`
	Latitude              string      `json:"latitude"`
	ImState               string      `json:"imState"`
	LastLogin             uint64      `json:"lastLogin"`
	Explain               interface{} `json:"explain"`
	Plus                  interface{} `json:"plus"`
	PcShow                int         `json:"pcShow"`
	AppShow               int         `json:"appShow"`
	Deliver               int         `json:"deliver"`
	GradeDescription      interface{} `json:"gradeDescription"`
	PromotionScoreExplain interface{} `json:"promotionScoreExplain"`
	FirstType             string      `json:"firstType"`
	SecondType            string      `json:"secondType"`
	IsSchoolJob           int         `json:"isSchoolJob"`
	SubwayLine            string      `json:"subwayline"`
	StationName           string      `json:"stationname"`
	LineStation           string      `json:"linestaion"`
	FormatCreateTime      string      `json:"formatCreateTime"`
	CompanyFullName       string      `json:"companyFullName"`
}
