package controllers

import (
	"encoding/json"
	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/models/shop/address"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/libs/logger"
)

type (
	AddressController struct {
		base.Controller
	}

	addReq struct {
		Address   string `json:"address"`
		IsDefault bool   `json:"isDefault"`
	}

	setDefaultReq struct {
		Id int `json:"id"`
	}

	modifyReq struct {
		Id        int    `json:"id"`
		Address   string `json:"address"`
		IsDefault bool   `json:"isDefault"`
	}
)

func (this *AddressController) AddAddress() {
	var (
		err      error
		req      addReq
		conn     orm.Connection
		claims   jwt.MapClaims
		userName string
		ok       bool
	)

	token, err := this.ParseToken()
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}
		goto finish
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}
	}
	userName = claims["username"].(string)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = address.Service.Add(conn, userName, req.Address, req.IsDefault)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}

func (this *AddressController) SetDefault() {
	var (
		err      error
		req      setDefaultReq
		conn     orm.Connection
		claims   jwt.MapClaims
		userName string
		ok       bool
	)

	token, err := this.ParseToken()
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}
		goto finish
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}
	}
	userName = claims["username"].(string)

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = address.Service.SetDefault(conn, userName, req.Id)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}

func (this *AddressController) Modify() {
	var (
		err  error
		req  modifyReq
		conn orm.Connection
	)

	_, err = this.ParseToken()
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrToken}
		goto finish
	}

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}
		goto finish
	}

	err = address.Service.Modify(conn, req.Id, req.Address)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed}

finish:
	this.ServeJSON(true)
}
