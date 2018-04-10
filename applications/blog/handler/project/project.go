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
 *     Initial: 2018/03/27        Chen Yanchen
 */

package project

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/models/blog/project"
)

func Create(c *server.Context) error {
	var req struct {
		Title    string `json:"title" validate:"required,max=32"`
		Abstract string `json:"abstract" validate:"max=64"`
		Content  string `json:"content" validate:"required"`
		Image    string `json:"image"`
	}

	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	authorID := c.Request().Context().Value("staff").(jwt.MapClaims)["staffid"].(float64)

	p := &project.Project{
		AuthorID: int32(authorID),
		Title:    req.Title,
		Abstract: req.Abstract,
		Content:  req.Content,
		Image:    req.Image,
	}
	err := project.ProjectServer.Creat(p)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func Delete(c *server.Context) error {
	var req struct {
		ID string `json:"id"`
	}

	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err := project.ProjectServer.Delete(req.ID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func Modify(c *server.Context) error {
	var req struct {
		ID       bson.ObjectId `json:"id" validate:"required"`
		Title    string        `json:"title" validate:"required,max=32"`
		Abstract string        `json:"abstract"`
		Content  string        `json:"content" validate:"required"`
		Image    string        `json:"image"`
	}

	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	p := &project.Project{
		ID:       req.ID,
		Title:    req.Title,
		Abstract: req.Abstract,
		Content:  req.Content,
		Image:    req.Image,
	}
	err := project.ProjectServer.Modify(p)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func GetID(c *server.Context) error {
	var req struct {
		Title string `json:"title" validate:"required"`
	}
	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	id, err := project.ProjectServer.GetID(req.Title)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, id)
}

func GetByID(c *server.Context) error {
	var req struct {
		ID string `json:"id" validate:"required"`
	}
	if err := c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err := project.ProjectServer.GetByID(bson.ObjectIdHex(req.ID))
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

func AbstractList(c *server.Context) error {
	res, err := project.ProjectServer.AbstractList()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}
