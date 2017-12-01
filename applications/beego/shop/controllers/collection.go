/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the 'Software'), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/12/01      Lin Hao
 */

package controllers

import (
	"encoding/json"

	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	Collection "github.com/fengyfei/gu/models/shop/collection"
)

type (
	// CollectionController - the controller of collection
	CollectionController struct {
		base.Controller
	}

	addCollectionReq struct {
		WareID int32 `json:"wareId"`
	}

	removeCollectionReq struct {
		ID int32 `json:"id"`
	}
)

// GetByUserID - get collections by userID.
func (cc *CollectionController) GetByUserID() {
	var (
		items  []Collection.CollectionItem
		err    error
		userID int32
		conn   orm.Connection
	)

	userID = cc.Ctx.Request.Context().Value("userId").(int32)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}

		goto finish
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	items, err = Collection.Service.GetByUserID(conn, userID)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	cc.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: constants.ErrSucceed,
		constants.RespKeyData:   items,
	}

finish:
	cc.ServeJSON()
}

// Add - add collection.
func (cc *CollectionController) Add() {
	var (
		req    addCollectionReq
		err    error
		userID int32
		conn   orm.Connection
	)

	userID = cc.Ctx.Request.Context().Value("userId").(int32)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}

		goto finish
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	err = json.Unmarshal(cc.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = cc.Validate(&req)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = Collection.Service.Add(conn, userID, req.WareID)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	cc.ServeJSON()
}

// Remove - remove collections.
func (cc *CollectionController) Remove() {
	var (
		req  removeCollectionReq
		err  error
		conn orm.Connection
	)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	err = json.Unmarshal(cc.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	err = cc.Validate(&req)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
	}

	err = Collection.Service.Remove(conn, req.ID)
	if err != nil {
		logger.Error(err)

		cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}

		goto finish
	}

	cc.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	cc.ServeJSON()
}
