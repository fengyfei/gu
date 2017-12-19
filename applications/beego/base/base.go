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

package base

import (
	"github.com/astaxie/beego"
	"gopkg.in/go-playground/validator.v9"

	"github.com/fengyfei/gu/libs/constants"
)

// Controller wraps general functionality.
type Controller struct {
	beego.Controller
}

// Validate the parameters.
func (base Controller) Validate(val interface{}) error {
	v := validator.New()

	return v.Struct(val)
}

// WriteStatusAndDataJSON write JSON encoding status (and data) to controller and sends a response.
// The data is optional.
func (base Controller) WriteStatusAndDataJSON(status int, data interface{}) {
	if data == nil {
		base.Data["json"] = map[string]interface{}{
			constants.RespKeyStatus: status,
		}

		base.ServeJSON(true)
	}

	base.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: status,
		constants.RespKeyData:   data,
	}

	base.ServeJSON(true)
}

// WriteStatusAndTokenJSON write JSON encoding status and token to controller and sends a response.
func (base Controller) WriteStatusAndTokenJSON(status int, token interface{}) {
	base.Data["json"] = map[string]interface{}{
		constants.RespKeyStatus: status,
		constants.RespKeyToken:  token,
	}

	base.ServeJSON(true)
}
