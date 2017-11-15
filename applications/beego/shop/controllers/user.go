package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/models/shop/user"
	"github.com/fengyfei/gu/libs/constants"
	"fmt"
	"net/http"
	"io/ioutil"
)

var (
	APPID  = ""
	SECRET = ""
)

type (
	UserController struct {
		beego.Controller
	}

	WechatLoginReq struct {
		UserName   string `json:"userName" validate:"required,alphanum,min=6,max=30"`
		WechatCode string `json:"wechatCode" validate:"required"`
	}

	wechatLogin struct {
		data wechatLoginData
	}

	wechatLoginData struct {
		errmsg  string
		unionid string
	}
)

func (u *UserController) WechatLogin() {
	var
	(
		wechatUser WechatLoginReq
		err        error
		userName   string
		url        string
		wechatData wechatLogin
		wechatRes  *http.Response
		con        []byte
	)
	json.Unmarshal(u.Ctx.Input.RequestBody, &wechatUser)

	conn, err := mysql.Pool.Get()
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}

	url = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", APPID, SECRET, wechatUser.WechatCode)

	wechatRes, err = http.Get(url)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	con, _ = ioutil.ReadAll(wechatRes.Body)
	err = json.Unmarshal(con, &wechatData)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	if wechatData.data.errmsg != "" {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrWechatAuth}
		goto finish
	}

	userName, err = user.Service.WechatLogin(conn, &wechatData.data.unionid)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	u.Data["json"] = map[string]string{"userName": userName}

finish:
	u.ServeJSON(true)
}

