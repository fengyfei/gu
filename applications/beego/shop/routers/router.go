// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
  "github.com/fengyfei/gu/applications/beego/shop/controllers"

  "github.com/astaxie/beego"
)

func init() {
  beego.Router("/shop/user/wechatlogin", &controllers.UserController{}, "post:WechatLogin")
  beego.Router("/shop/user/register", &controllers.UserController{}, "post:PhoneRegister")
  beego.Router("/shop/user/login", &controllers.UserController{}, "post:PhoneLogin")
  beego.Router("/shop/category/getmainclass", &controllers.CategoryController{}, "get:GetMainCategories")
  beego.Router("/shop/category/getsubclass", &controllers.CategoryController{}, "post:GetSubCategories")
  beego.Router("/shop/category/add", &controllers.CategoryController{}, "post:AddCategory")
  beego.Router("/shop/user/changepass", &controllers.UserController{}, "post:ChangePassword")
}
