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
 *     Initial: 2017/10/27        Jia Chenhui
 */

package tag_test

import (
	"testing"

	"github.com/fengyfei/gu/models/blog/tag"
	"github.com/fengyfei/gu/pkg/log"
)

var (
	walk = tag.MDCreateTag{
		Tag: "walk",
	}
)

func TestCreateAndGetByID(t *testing.T) {
	m := "TestCreateAndGetByID"
	tag.Prepare()

	id, err := tag.Service.Create(walk.Tag)
	checkError("Create", err)

	result, err := tag.Service.GetByID(id)
	checkError("GetByID", err)

	if !checkGetByIDResp(result) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestGetListGetActiveList(t *testing.T) {
	m := "TestGetListGetActiveList"
	tag.Prepare()

	list, err := tag.Service.GetList()
	checkError("GetList", err)

	activeList, err := tag.Service.GetActiveList()
	checkError("GetActiveList", err)

	if len(list) != len(activeList) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func TestCreateModify(t *testing.T) {
	m := "TestCreateModify"
	tag.Prepare()

	id, err := tag.Service.Create(walk.Tag)
	checkError("Create", err)

	stop := tag.MDModifyTag{
		TagID:  id,
		Tag:    "stop",
		Active: false,
	}

	err = tag.Service.Modify(&stop)
	checkError("Modify", err)

	newStop := tag.MDModifyTag{
		TagID:  id,
		Tag:    "walk",
		Active: true,
	}

	err = tag.Service.Modify(&newStop)
	checkError("Modify", err)

	a, err := tag.Service.GetByID(id)
	checkError("GetByID", err)

	if !checkModifyResp(a) {
		log.Logger.Debug("%s failure.", m)
	} else {
		log.Logger.Debug("%s success.", m)
	}
}

func checkError(method string, err error) {
	if err != nil {
		log.Logger.Debug("%s returned error: %s", method, err)
	} else {
		log.Logger.Debug("%s execute success.", method)
	}
}

func checkGetByIDResp(resp tag.MDTag) bool {
	return resp.Tag == walk.Tag && resp.Active == true
}

func checkModifyResp(resp tag.MDTag) bool {
	return resp.Tag == walk.Tag && resp.Active == true
}
