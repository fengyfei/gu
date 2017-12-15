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
 *     Initial: 2017/11/01        Jia Chenhui
 */

package router

import (
	"github.com/astaxie/beego"

	"github.com/fengyfei/gu/applications/beego/polaris/controller"
	"github.com/fengyfei/gu/applications/echo/core"
)

func InitRouter() {
	// JWT middleware does not affect these route.
	core.URLMap["/api/v1/office/staff/login"] = struct{}{}

	// Staff
	beego.Router("/api/v1/office/staff/login", &controller.Staff{}, "post:Login")
	beego.Router("/api/v1/office/staff/create", &controller.Staff{}, "post:Create")
	beego.Router("/api/v1/office/staff/modify/info", &controller.Staff{}, "post:Modify")
	beego.Router("/api/v1/office/staff/modify/pwd", &controller.Staff{}, "post:ModifyPwd")
	beego.Router("/api/v1/office/staff/modify/mobile", &controller.Staff{}, "post:ModifyMobile")
	beego.Router("/api/v1/office/staff/activate", &controller.Staff{}, "post:ModifyActive")
	beego.Router("/api/v1/office/staff/dismiss", &controller.Staff{}, "post:Dismiss")
	beego.Router("/api/v1/office/staff/list", &controller.Staff{}, "get:List")
	beego.Router("/api/v1/office/staff/info", &controller.Staff{}, "post:Info")

	// Relation
	beego.Router("/api/v1/office/staff/relation/create", &controller.Staff{}, "post:AddRole")
	beego.Router("/api/v1/office/staff/relation/remove", &controller.Staff{}, "post:RemoveRole")
	beego.Router("/api/v1/office/staff/relation/list", &controller.Staff{}, "post:RoleList")

	// Role
	beego.Router("/api/v1/office/role/create", &controller.Role{}, "post:Create")
	beego.Router("/api/v1/office/role/modify/info", &controller.Role{}, "post:Modify")
	beego.Router("/api/v1/office/role/activate", &controller.Role{}, "post:ModifyActive")
	beego.Router("/api/v1/office/role/list", &controller.Role{}, "get:List")
	beego.Router("/api/v1/office/role/detail", &controller.Role{}, "post:Info")

	// Permission
	beego.Router("/api/v1/office/permission/create", &controller.Permission{}, "post:Create")
	beego.Router("/api/v1/office/permission/remove", &controller.Permission{}, "post:Remove")
	beego.Router("/api/v1/office/permission/list", &controller.Permission{}, "get:List")
}
