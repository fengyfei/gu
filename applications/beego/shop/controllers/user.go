package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/fengyfei/gu/applications/beego/shop/mysql"
	"github.com/fengyfei/gu/models/shop/user"
	_ "github.com/fengyfei/gu/libs/constants"
	_ "github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/constants"
	"fmt"
	"net/http"
	"io/ioutil"
)

var (
	APPID  = ""
	SECRET = ""
)

// Operations about Users
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
		errmsg string
		unionid string
	}
)

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /addwechatuser [post]
func (u *UserController) WechatLogin() {
	var
	(
		wechatUser WechatLoginReq
		err        error
		uid        string
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

	uid, err = user.Service.WechatLogin(conn, &wechatData.data.unionid)
	if err != nil {
		u.Data["json"] = map[string]interface{}{constants.RespKeyStatus: constants.ErrMysql}
		goto finish
	}
	u.Data["json"] = map[string]string{"uid": uid}

finish:
	u.ServeJSON(true)
}
