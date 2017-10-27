/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/10/22        Feng Yifei
 */

package routers

import (
	"github.com/astaxie/beego"
	"github.com/fengyfei/gu/applications/blog/controllers"
)

func init() {
	beego.Router("/", &controllers.Web{}, "*:Home")

	// tag router
	beego.Router("/blog/tag/list", &controllers.Tag{}, "get:List")
	beego.Router("/blog/tag/activelist", &controllers.Tag{}, "get:ActiveList")
	beego.Router("/blog/tag/info", &controllers.Tag{}, "post:TagInfo")
	beego.Router("/blog/tag/create", &controllers.Tag{}, "post:Create")
	beego.Router("/blog/tag/modify", &controllers.Tag{}, "post:Modify")

	// bio router
	beego.Router("/blog/bio/info", &controllers.Bio{}, "get:BioInfo")
	beego.Router("/blog/bio/create", &controllers.Bio{}, "post:Create")
}
