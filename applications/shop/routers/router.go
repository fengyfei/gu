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
	r.Post("/shop/account/wechatlogin", handler.WechatLogin)
	r.Post("/shop/account/addphone", handler.AddPhone)
	r.Post("/shop/account/changeinfo", handler.ChangeInfo)
	r.Post("/shop/account/register", handler.PhoneRegister)
	r.Post("/shop/account/phonelogin", handler.PhoneLogin)
	r.Post("/shop/account/changepass", handler.ChangePassword)

	r.Post("/shop/address/add", handler.AddAddress)
	r.Post("/shop/address/setdefault", handler.SetDefaultAddress)
	r.Post("/shop/address/modify", handler.ModifyAddress)
	r.Post("/shop/address/delete", handler.DeleteAddress)
	r.Get("/shop/address/get", handler.GetAddress)

	r.Post("/shop/category/add", handler.AddCategory)
	r.Post("/shop/category/modify", handler.ModifyCategory)
	//r.Post("/shop/category/get", handler.GetCategory)

	// ware api for admin
	r.Post("/shop/api/ware/create", handler.CreateWare)
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
	r.Post("/shop/cart/add", handler.AddCart)
	r.Post("/shop/cart/userid", handler.GetByUser)
	r.Post("/shop/cart/remove", handler.Remove)
	r.Post("/shop/cart/removemany", handler.RemoveMany)

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
