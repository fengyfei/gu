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

package article

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/TechCatsLab/apix/http/server"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/github/article"
)

type (
	// createReq - The request struct that create article information.
	createReq struct {
		Title  *string `json:"title" validate:"required,min=1,max=256"`
		URL    *string `json:"url" validate:"required,url"`
		Source *string `json:"source"`
	}

	// activateReq - The request struct that modify article status.
	activateReq struct {
		ID     string `json:"id" validate:"required,alphanum,len=24"`
		Active bool   `json:"active"`
	}

	// infoReq - The request struct that get a list of articles detail information.
	infoReq struct {
		ID string `json:"id"`
	}

	// infoResp - The more detail of article.
	infoResp struct {
		ID      string    `json:"id"`
		Title   string    `json:"title"`
		URL     string    `json:"url"`
		Source  string    `json:"source"`
		Active  bool      `json:"active"`
		Created time.Time `json:"created"`
	}
)

// Create - Create article information.
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

	if req.Source == nil {
		req.Source = emptyStr
	}

	id, err := article.Service.Create(req.Title, req.URL, req.Source)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// ModifyActive - Modify article status.
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

	if err = article.Service.ModifyActive(&req.ID, req.Active); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// List - Get all the articles.
func List(c *server.Context) error {
	var resp = make([]infoResp, 0)

	alist, err := article.Service.List()
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
			URL:     a.URL,
			Source:  a.Source,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// ActiveList - Get all the active articles.
func ActiveList(c *server.Context) error {
	var resp = make([]infoResp, 0)

	alist, err := article.Service.ActiveList()
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
			URL:     a.URL,
			Source:  a.Source,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// Info - Get detail information for specified article.
func Info(c *server.Context) error {
	var (
		err  error
		req  infoReq
		resp = make([]infoResp, 0)
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	alist, err := article.Service.GetByID(req.ID)
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
			URL:     a.URL,
			Source:  a.Source,
			Active:  a.Active,
			Created: a.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
