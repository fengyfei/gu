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

	// address api for user
	r.Post("/shop/address/add", handler.AddAddress)
	r.Post("/shop/address/setdefault", handler.SetDefaultAddress)
	r.Post("/shop/address/modify", handler.ModifyAddress)
	r.Post("/shop/address/delete", handler.DeleteAddress)
	r.Get("/shop/address/get", handler.GetAddress)

	// category api for admin
	r.Post("/shop/category/add", handler.AddCategory)
	r.Post("/shop/category/maincategory", handler.GetMainCategories)
	r.Post("/shop/category/subcategory", handler.GetSubCategories)
	r.Post("/shop/category/delete", handler.Delete)

	// ware api for admin
	r.Post("/shop/api/ware/create", handler.CreateWare)
	r.Get("/shop/api/ware/getall", handler.GetAllWare)
	r.Post("/shop/api/ware/updateinfo", handler.UpdateWithID)
	r.Post("/shop/api/ware/modifyprice", handler.ModifyPrice)
	r.Post("/shop/api/ware/changestatus", handler.ChangeStatus)

	// ware api for user
	r.Post("/shop/ware/getbycid", handler.GetWareByCategory)
	r.Get("/shop/ware/getpromotion", handler.GetPromotion)
	r.Post("/shop/ware/gethomelist", handler.HomePageList)
	r.Get("/shop/ware/getrecommend", handler.RecommendList)
	r.Post("/shop/ware/getdetail", handler.GetDetail)
	r.Get("/shop/ware/getnew", handler.GetNewWares)

	// cart api
	r.Post("/shop/cart/add", handler.AddCart)
	r.Post("/shop/cart/userid", handler.GetByUser)
	r.Post("/shop/cart/remove", handler.Remove)
	r.Post("/shop/cart/removemany", handler.RemoveWhenOrder)

	// order api for user
	r.Post("/shop/user/order/add", handler.CreateOrder)
	r.Get("/shop/user/order/get", handler.GetUserOrder)

	// order api for admin
	r.Post("/shop/api/order/confirm", handler.ConfirmOrder)

	// collection api for user
	r.Get("/shop/collection/get", handler.GetByUserID)
	r.Post("/shop/collection/add", handler.AddColl)
	r.Post("/shop/collection/remove", handler.RemoveColl)

	// panel api for admin
	r.Post("/shop/api/panel/create", handler.AddPanel)
	r.Post("/shop/api/panel/addpromotion", handler.AddPromotion)
	r.Post("/shop/api/panel/addrecommend", handler.AddRecommend)

	// panel api for user
	r.Get("/shop/panel/getpage", handler.GetPanelPage)
}
