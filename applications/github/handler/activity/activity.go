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
 *     Initial: 2018/03/06        Lin Hao
 */

package activity

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/TechCatsLab/apix/http/server"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/github/activity"
)

type (
	// createReq - The request struct that create activity information.
	createReq struct {
		Title *string `json:"title" validate:"required,min=1,max=256"`
		Image *string `json:"image"`
		Intro *string `json:"intro"`
	}

	// activateReq - The request struct that modify activity status.
	activateReq struct {
		ID     string `json:"id" validate:"required,alphanum,len=24"`
		Active bool   `json:"active"`
	}

	// infoResp - The more detail of activity.
	infoResp struct {
		ID      string    `json:"id"`
		Title   string    `json:"title"`
		Image   string    `json:"image"`
		Intro   string    `json:"intro"`
		Active  bool      `json:"active"`
		Created time.Time `json:"created"`
	}
)

// Create - Create activity information.
func Create(c *server.Context) error {
	var (
		err      error
		req      createReq
		emptyStr = new(string)
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	switch {
	case req.Image == nil:
		req.Image = emptyStr
	case req.Intro == nil:
		req.Intro = emptyStr
	}

	id, err := activity.Service.Create(req.Title, req.Image, req.Intro)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// ModifyActive - Modify activity status.
func ModifyActive(c *server.Context) error {
	var (
		err error
		req activateReq
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = activity.Service.ModifyActive(&req.ID, req.Active); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// List - Get all the activities.
func List(c *server.Context) error {
	var resp = make([]infoResp, 0)

	alist, err := activity.Service.List()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, a := range alist {
		info := infoResp{
			ID:      a.ID.Hex(),
			Title:   a.Title,
			Image:   a.Image,
			Intro:   a.Intro,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// ActiveList - Get all the active activities.
func ActiveList(c *server.Context) error {
	var resp = make([]infoResp, 0)

	alist, err := activity.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, a := range alist {
		info := infoResp{
			ID:      a.ID.Hex(),
			Title:   a.Title,
			Image:   a.Image,
			Intro:   a.Intro,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
