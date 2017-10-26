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
 *     Initial: 2017/10/26        Jia Chenhui
 */

package article_test

import (
	"testing"

	"github.com/fengyfei/gu/models/blog/article"
	"github.com/fengyfei/gu/pkg/log"
)

var (
	running = article.MDCreateArticle{
		Author:   "jch",
		Title:    "running",
		Content:  "running running",
		Abstract: "run",
		Tag:      []string{"run", "sport"},
	}
)

func TestCreateGetByID(t *testing.T) {
	m := "TestCreateGetByID"
	article.Prepare()

	id, err := article.Service.Create(&running)
	checkError("Create", err)

	a, err := article.Service.GetByID(id)
	checkError("GetByID", err)

	if !checkGetByIDResp(a) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestGetListGetActiveList(t *testing.T) {
	m := "TestGetListGetActiveList"
	article.Prepare()

	list, err := article.Service.GetList()
	checkError("GetList", err)

	activeList, err := article.Service.GetActiveList()
	checkError("GetActiveList", err)

	if len(list) != len(activeList) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestCreateModify(t *testing.T) {
	m := "TestCreateModify"
	article.Prepare()

	id, err := article.Service.Create(&running)
	checkError("Create", err)

	stop := article.MDModifyArticle{
		ArticleID: id,
		Title:     "stop",
		Content:   "stop stop",
		Abstract:  "stop",
	}

	err = article.Service.Modify(&stop)
	checkError("Modify", err)

	a, err := article.Service.GetByID(id)
	checkError("GetByID", err)

	if !checkModifyResp(a) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestGetListGetByTags(t *testing.T) {
	m := "TestCreateGetByTags"
	tags := []string{"run", "sport"}
	article.Prepare()

	l1, err := article.Service.GetList()
	checkError("GetList", err)

	l2, err := article.Service.GetByTags(tags)
	checkError("GetByTags", err)

	if len(l1) != len(l2) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestCreateAddTagsRemoveTags(t *testing.T) {
	m := "TestCreateAddTagsRemoveTags"
	tags := []string{"a", "b"}
	article.Prepare()

	l1, err := article.Service.GetByTags(tags)
	checkError("GetByTags 1", err)

	id, err := article.Service.Create(&running)
	checkError("Create", err)

	err = article.Service.AddTags(id, tags)
	checkError("AddTags", err)

	a, err := article.Service.GetByTags(tags)
	checkError("GetByTags 2", err)
	if len(a) != 1 {
		log.Logger.Debug("%s failure.", "GetByTags 2")
	} else {
		log.Logger.Debug("%s success.", "GetByTags 2")
	}

	err = article.Service.RemoveTags(id, tags)
	checkError("RemoveTags", err)

	l2, err := article.Service.GetByTags(tags)
	checkError("GetByTags 3", err)

	if len(l1) != len(l2) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func checkError(method string, err error) {
	if err != nil {
		log.Logger.Debug("%s returned error: %s", method, err)
	}

	log.Logger.Debug("%s execute success.", method)
}

func checkGetByIDResp(resp article.MDArticle) bool {
	return resp.Author == "jch" && resp.Title == "running" && resp.Content == "running running"
}

func checkModifyResp(resp article.MDArticle) bool {
	return resp.Title == "stop" && resp.Content == "stop stop" && resp.Abstract == "stop" && resp.Active == false
}
