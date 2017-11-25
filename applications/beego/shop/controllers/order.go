package controllers

import (
	"github.com/fengyfei/gu/applications/beego/base"
	"github.com/fengyfei/gu/libs/constants"
	"encoding/json"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/dgrijalva/jwt-go"
	"github.com/fengyfei/gu/libs/orm"
	Order "github.com/fengyfei/gu/models/shop/order"
)

type (
	OrderController struct {
		base.Controller
	}

	createReq struct {
		Orders []Order.OrderItem `json:"orders" validate:"required"`
	}
)

func (this *OrderController) CreateOrder() {
	var (
		req    createReq
		err    error
		IP     string
		claims jwt.MapClaims
		ok     bool
		userId int32
		conn   orm.Connection
		signStr string
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
	userId = int32(claims["userid"].(float64)) //strange, maybe it is conversed int32 to float64 when parsing the token

	conn, err = mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	IP = this.Ctx.Input.IP()

	err = this.Validate(&req)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrInvalidParam}

		goto finish
	}

	signStr, err = Order.Service.OrderByWechat(conn, userId, IP, &req.Orders)
	if err != nil {
		logger.Error(err)
		this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatPay}

		goto finish
	}
	this.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrSucceed, constants.RespPaySign: signStr}

finish:
	this.ServeJSON(true)
}
