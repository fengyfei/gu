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
 *     Initial: 2017/12/28        Jia Chenhui
 */

package repos

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/github/repos"
)

type (
	// createReq - The request struct that create repos information.
	createReq struct {
		Avatar *string  `json:"avatar" validate:"required,url"`
		Name   *string  `json:"name" validate:"required,printascii,excludesall=;0x2C"`
		Image  *string  `json:"image"`
		Intro  *string  `json:"intro"`
		Lang   []string `json:"lang"`
	}

	// activateReq - The request struct that modify repos status.
	activateReq struct {
		ID     string `json:"id" validate:"required,alphanum,len=24"`
		Active bool   `json:"active"`
	}

	// infoReq - The request struct that get a list of repos detail information.
	infoReq struct {
		ID string `json:"id"`
	}

	// infoResp - The more detail of repos.
	infoResp struct {
		ID      string    `json:"id"`
		Avatar  string    `json:"avatar"`
		Name    string    `json:"name"`
		Image   string    `json:"image"`
		Intro   string    `json:"intro"`
		Lang    []string  `json:"lang"`
		Created time.Time `json:"created"`
		Active  bool      `json:"active"`
	}
)

// Create - Create repos information.
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
	case req.Lang == nil:
		req.Lang = make([]string, 0)
	}

	id, err := repos.Service.Create(req.Avatar, req.Name, req.Image, req.Intro, req.Lang)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// ModifyActive - Modify repos status.
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

	if err = repos.Service.ModifyActive(&req.ID, req.Active); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// List - Get all the repos.
func List(c *server.Context) error {
	var resp []infoResp = make([]infoResp, 0)

	rlist, err := repos.Service.List()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, r := range rlist {
		info := infoResp{
			ID:      r.ID.Hex(),
			Avatar:  r.Avatar,
			Name:    r.Name,
			Image:   r.Image,
			Intro:   r.Intro,
			Lang:    r.Lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// ActiveList - Get all the active repos.
func ActiveList(c *server.Context) error {
	var resp []infoResp = make([]infoResp, 0)

	rlist, err := repos.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, r := range rlist {
		info := infoResp{
			ID:      r.ID.Hex(),
			Avatar:  r.Avatar,
			Name:    r.Name,
			Image:   r.Image,
			Intro:   r.Intro,
			Lang:    r.Lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// Info - Get ten records that are greater than the specified ID.
func Info(c *server.Context) error {
	var (
		err  error
		req  infoReq
		resp []infoResp = make([]infoResp, 0)
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	rlist, err := repos.Service.GetByID(req.ID)
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, r := range rlist {
		info := infoResp{
			ID:      r.ID.Hex(),
			Avatar:  r.Avatar,
			Name:    r.Name,
			Image:   r.Image,
			Intro:   r.Intro,
			Lang:    r.Lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
