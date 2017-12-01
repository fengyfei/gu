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
 *     Initial: 2017/11/03        ShiChao
 */

package routers

import (
	"github.com/fengyfei/gu/applications/beego/shop/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/shop/user/wechatlogin", &controllers.UserController{}, "post:WechatLogin")
	beego.Router("/shop/user/register", &controllers.UserController{}, "post:PhoneRegister")
	beego.Router("/shop/user/login", &controllers.UserController{}, "post:PhoneLogin")
	beego.Router("/shop/user/changepass", &controllers.UserController{}, "post:ChangePassword")

	beego.Router("/shop/address/add", &controllers.AddressController{}, "post:AddAddress")
	beego.Router("/shop/address/setdefault", &controllers.AddressController{}, "post:SetDefault")
	beego.Router("/shop/address/modify", &controllers.AddressController{}, "post:Modify")

	beego.Router("/shop/api/category/add", &controllers.CategoryController{}, "post:AddCategory")
	beego.Router("/shop/category/getmainclass", &controllers.CategoryController{}, "get:GetMainCategories")
	beego.Router("/shop/category/getsubclass", &controllers.CategoryController{}, "post:GetSubCategories")

	// ware api for admin
	beego.Router("/shop/api/ware/create", &controllers.WareController{}, "post:CreateWare")
	beego.Router("/shop/api/ware/getall", &controllers.WareController{}, "get:GetWareList")
	beego.Router("/shop/api/ware/updateinfo", &controllers.WareController{}, "post:UpdateWithID")
	beego.Router("/shop/api/ware/modifyprice", &controllers.WareController{}, "post:ModifyPrice")
	beego.Router("/shop/api/ware/changestatus", &controllers.WareController{}, "post:ChangeStatus")

	// ware api for user
	beego.Router("/shop/ware/getbycid", &controllers.WareController{}, "post:GetWareByCategory")
	beego.Router("/shop/ware/getpromotion", &controllers.WareController{}, "get:GetPromotion")
	beego.Router("/shop/ware/homelist", &controllers.WareController{}, "post:HomePageList")
	beego.Router("/shop/ware/recommend", &controllers.WareController{}, "post:RecommendList")
	beego.Router("/shop/ware/getdetail", &controllers.WareController{}, "post:GetDetail")

	// cart api
	beego.Router("/shop/cart/add", &controllers.CartController{}, "post:Add")
	beego.Router("/shop/cart/remove", &controllers.CartController{}, "post:remove")
	beego.Router("/shop/cart/get", &controllers.CartController{}, "get:GetByUser")

	// collection api for user
	beego.Router("/shop/collection/get", &controllers.CollectionController{}, "get:GetByUserID")
	beego.Router("/shop/collection/add", &controllers.CollectionController{}, "post:Add")
	beego.Router("/shop/collection/remove", &controllers.CollectionController{}, "post:Remove")
}
