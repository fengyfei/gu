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
 *     Initial: 2018/02/01        Shi Ruitao
 */

package routers

import (
	"github.com/fengyfei/gu/applications/shop/handler"
	"github.com/fengyfei/gu/libs/http/server"
)

var (
	Router *server.Router
)

func init() {
	Router = server.NewRouter()
	register(Router)
}

func register(r *server.Router) {
	r.Post("/shop/user/wechatlogin", handler.WechatLogin)
	r.Post("/shop/user/addphone", handler.AddPhone)
	r.Post("/shop/user/changeinfo", handler.ChangeInfo)
	r.Post("/shop/user/register", handler.PhoneRegister)
	r.Post("/shop/user/phonelogin", handler.PhoneLogin)
	r.Post("/shop/user/changepass", handler.ChangePassword)

	r.Post("/shop/user/address/add", handler.AddAddress)
	r.Post("/shop/user/address/setdefault", handler.SetDefault)
	r.Post("/shop/user/address/modify", handler.Modify)
	r.Get("/shop/user/address/alladdress", handler.AddressRead)

	//r.Post("/shop/api/category/add", "post:AddCategory")
	//r.Get("/shop/category/getmainclass",  "get:GetMainCategories")
	//r.Post("/shop/category/getsubclass",  "post:GetSubCategories")

	// ware api for admin
	//r.Post("/shop/api/ware/create",  "post:CreateWare")
	//r.Get("/shop/api/ware/getall",  "get:GetAllWare")
	//r.Post("/shop/api/ware/updateinfo",  "post:UpdateWithID")
	//r.Post("/shop/api/ware/modifyprice",  "post:ModifyPrice")
	//r.Post("/shop/api/ware/changestatus",  "post:ChangeStatus")

	// ware api for user
	//r.Post("/shop/ware/getbycid", "post:GetWareByCategory")
	//r.Get("/shop/ware/getpromotion",  "get:GetPromotion")
	//r.Post("/shop/ware/homelist",  "post:HomePageList")
	//r.Post("/shop/ware/recommend", "post:RecommendList")
	//r.Post("/shop/ware/getdetail", "post:GetDetail")

	// cart api
	//r.Post("/shop/cart/add", &handler.CartController{}, "post:Add")
	//r.Post("/shop/cart/remove", &handler.CartController{}, "post:Remove")
	//r.Post("/shop/cart/get", &handler.CartController{}, "get:GetByUser")

	// order api for user
	//r.Post("/shop/user/order/add", &handler.OrderController{}, "post:CreateOrder")
	//r.Post("/shop/user/order/get", &handler.OrderController{}, "get:GetUserOrder")

	// order api for admin
	//r.Post("/shop/api/order/confirm", &handler.OrderController{}, "post:ConfirmOrder")

	// collection api for user
	//r.Post("/shop/collection/get", &handler.CollectionController{}, "get:GetByUserID")
	//r.Post("/shop/collection/add", &handler.CollectionController{}, "post:Add")
	//r.Post("/shop/collection/remove", &handler.CollectionController{}, "post:Remove")

	// panel api for admin
	//r.Post("/shop/api/panel/create", &handler.PanelController{}, "post:AddPanel")
	//r.Post("/shop/api/panel/addpromotion", &handler.PanelController{}, "post:AddPromotion")
	//r.Post("/shop/api/panel/addrecommend", &handler.PanelController{}, "post:AddRecommend")

	// panel api for user
	//r.Post("/shop/panel/getpage", &handler.PanelController{}, "get:GetPanelPage")
}
