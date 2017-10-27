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

package bio_test

import (
	"testing"

	"github.com/fengyfei/gu/models/blog/bio"
	"github.com/fengyfei/gu/pkg/log"
)

var (
	me = bio.MDCreateBio{
		Title: "about me",
		Bio:   "biobio",
	}

	newMe = bio.MDCreateBio{
		Title: "new me",
		Bio:   "bio",
	}
)

func TestCreateAndGetBio(t *testing.T) {
	m := "TestCreateAndGetBio"
	bio.Prepare()

	err := bio.Service.Create(&me)
	checkError("Create 1", err)

	b1, err := bio.Service.GetBio()
	checkError("GetBio 1", err)

	err = bio.Service.Create(&newMe)
	checkError("Create 2", err)

	b2, err := bio.Service.GetBio()
	checkError("GetBio 2", err)

	if b1.BioID != b2.BioID {
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
