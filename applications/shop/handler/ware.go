/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2018/03/06        Shi Ruitao
 */

package handler

import (
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/shop/mysql"
	"github.com/fengyfei/gu/applications/shop/util"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/shop/ware"
)

type (
	categoryReq struct {
		ParentCID uint32 `json:"parent_cid" validate:"required"`
		CID       uint32 `json:"cid"`
	}

	homeReq struct {
		LastID uint32 `json:"last_id"`
	}

	detailReq struct {
		ID uint32 `json:"id" validate:"required"`
	}

	ChangeStatusReq struct {
		IDs    []uint32 `json:"i_ds" validate:"required,min=1"`
		Status int8   `json:"status" validate:"required,eq=-1|eq=1|eq=2|eq=3"`
	}
)

// add new ware
func CreateWare(c *server.Context) error {
	var (
		err    error
		addReq ware.Ware
		conn   orm.Connection
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	conn, err = mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
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

	if len(addReq.Avatar) > 0 {
		addReq.Avatar, err = util.SavePicture(addReq.Avatar, "ware/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}
	if len(addReq.Image) > 0 {
		addReq.Image, err = util.SavePicture(addReq.Image, "ware/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}
	if len(addReq.DetailPic) > 0 {
		addReq.DetailPic, err = util.SavePicture(addReq.DetailPic, "wareIntro/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}

	err = ware.Service.CreateWare(conn, addReq)
	if err != nil {
		logger.Error(err)
		if (len(addReq.Avatar) > 0 && !util.DeletePicture(addReq.Avatar)) ||
			(len(addReq.Image) > 0 && !util.DeletePicture(addReq.Image)) ||
			(len(addReq.DetailPic) > 0 && !util.DeletePicture(addReq.DetailPic)) {
			logger.Error(errors.New("create ware failed and delete it's pictures go wrong, please delete picture manually"))
		}
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("create ware", addReq.Name, "success")
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// get all wares
func GetAllWare(c *server.Context) error {
	var (
		err error
		res []ware.Ware
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = ware.Service.GetAllWare(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get ware by categoryID
func GetWareByCategory(c *server.Context) error {
	var (
		err    error
		cidReq categoryReq
		res    []ware.BriefInfo
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&cidReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(cidReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if cidReq.CID == 0 {
		res, err = ware.Service.GetByParentCID(conn, cidReq.ParentCID)
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	} else {
		res, err = ware.Service.GetByCID(conn, cidReq.CID)
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get new wares
func GetNewWares(c *server.Context) error {
	var (
		err error
		res []ware.BriefInfo
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = ware.Service.GetNewWares(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get promotion wares
func GetPromotion(c *server.Context) error {
	var (
		err error
		res []ware.BriefInfo
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = ware.Service.GetPromotionList(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// update ware info
func UpdateWithID(c *server.Context) error {
	var (
		err error
		req ware.UpdateReq
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if len(req.Avatar) > 0 {
		req.Avatar, err = util.SavePicture(req.Avatar, "ware/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}
	if len(req.Image) > 0 {
		req.Image, err = util.SavePicture(req.Image, "ware/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}
	if len(req.DetailPic) > 0 {
		req.DetailPic, err = util.SavePicture(req.DetailPic, "wareIntro/")
		if err != nil {
			logger.Error(err)
			return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
		}
	}

	err = ware.Service.UpdateWare(conn, req)
	if err != nil {
		logger.Error(err)
		if (len(req.Avatar) > 0 && !util.DeletePicture(req.Avatar)) ||
			(len(req.Image) > 0 && !util.DeletePicture(req.Image)) ||
			(len(req.DetailPic) > 0 && !util.DeletePicture(req.DetailPic)) {
			logger.Error(errors.New("update ware failed and delete it's pictures go wrong, please delete picture manually"))
		}
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("update ware info of id:", req.ID)
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// modify price by id
func ModifyPrice(c *server.Context) error {
	var (
		err error
		req ware.ModifyPriceReq
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = ware.Service.ModifyPrice(conn, req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("modify price of", req.ID)
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// get homepage list with last wareID
func HomePageList(c *server.Context) error {
	var (
		err   error
		idReq homeReq
		res   []ware.BriefInfo
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&idReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(idReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err = ware.Service.HomePageList(conn, idReq.LastID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get recommend list
func RecommendList(c *server.Context) error {
	var (
		err error
		res []ware.BriefInfo
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	res, err = ware.Service.GetRecommendList(conn)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// get detail info of ware
func GetDetail(c *server.Context) error {
	var (
		err error
		req detailReq
		res *ware.Ware
	)

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	res, err = ware.Service.GetByID(conn, req.ID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, res)
}

// change status of wares
func ChangeStatus(c *server.Context) error {
	var (
		err       error
		changeReq ChangeStatusReq
	)

	isAdmin := c.Request().Context().Value("user").(jwtgo.MapClaims)[util.IsAdmin].(bool)
	if !isAdmin {
		logger.Error("You don't have access")
		return core.WriteStatusAndDataJSON(c, constants.ErrToken, nil)
	}

	conn, err := mysql.Pool.Get()
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = c.JSONBody(&changeReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err = c.Validate(changeReq)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	// status: -1 -> delete or sold out
	// 1 -> common
	// 2 -> promotion
	// 3 -> today new wares
	err = ware.Service.ChangeStatus(conn, changeReq.IDs, changeReq.Status)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	logger.Info("change ware status success")
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
