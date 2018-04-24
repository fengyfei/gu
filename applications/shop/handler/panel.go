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
 *     Initial: 2018/03/16        Shi Ruitao
 */

package handler

import (
	"errors"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/shop/panel"
	"github.com/fengyfei/gu/models/shop/ware"
)

// add promotion panel
func AddPanel(c *server.Context) error {
	var (
		err    error
		addReq panel.Panel
		conn   orm.Connection
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	err = c.JSONBody(&addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = panel.Service.CreatePanel(conn, &addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("create panel", addReq.Title, "success")
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// add promotion list
func AddPromotion(c *server.Context) error {
	var (
		err    error
		addReq panel.PromotionReq
		conn   orm.Connection
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	err = c.JSONBody(&addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = panel.Service.AddPromotionList(conn, addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("add promotion list success")
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// add recommend
func AddRecommend(c *server.Context) error {
	var (
		err    error
		addReq panel.RecommendReq
		conn   orm.Connection
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	err = c.JSONBody(&addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(addReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	addReq.Picture, err = util.SavePicture(addReq.Picture, "recommend/")
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = panel.Service.AddRecommend(conn, addReq)
	if err != nil {
		logger.Error(err)
		if !util.DeletePicture(addReq.Picture) {
			logger.Error(errors.New("add recommend failed and delete it's pictures go wrong, please delete picture manually"))
		}
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("add recommend of panel success")
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// TODO: add second-hand
func AddSecondHand(c *server.Context) {}

// get panel page
func GetPanelPage(c *server.Context) error {
	var (
		err error
		res []panel.PanelsPage
	)

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = panel.Service.GetPanels(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	for i := range res {
		if res[i].Type == 1 {
			detail, err := panel.Service.GetDetail(conn, res[i].ID)
			if err != nil {
				logger.Error(err)
				res[i].Content = []interface{}{}
			} else {
				ids := strings.Split(detail.Content, "#")

				wares, err := ware.Service.GetByIDs(conn, ids)
				if err != nil {
					logger.Error(err)
					res[i].Content = []interface{}{}
				} else {
					for k := range wares {
						res[i].Content = append(res[i].Content, wares[k])
					}
				}
			}
		}
		if res[i].Type == 2 {
			detail, err := panel.Service.GetDetail(conn, res[i].ID)
			if err != nil {
				logger.Error(err)
				res[i].Content = []interface{}{}
			} else {
				res[i].Content = append(res[i].Content, detail.Picture)
			}
		}
		if res[i].Type == 3 {
			if newWares, newErr := ware.Service.GetNewWares(conn); len(newWares) > 0 {
				if newErr != nil {
					logger.Error(newErr)
					res[i].Content = []interface{}{}
				} else {
					res[i].Content = append(res[i].Content, newWares)
				}
			} else {
				res[i].Content = []interface{}{}
			}
		}
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}
