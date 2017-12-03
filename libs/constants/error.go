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
 *     Initial: 2017/10/28        Feng Yifei
 */

package constants

const (
	// ErrSucceed - Succeed
	ErrSucceed = 0

	// ErrPermission - Permission Denied
	ErrPermission = 401
	ErrForbidden  = 438

	// ErrToken - Invalid Token
	ErrToken = 420

	// ErrInvalidParam - Invalid Parameter
	ErrInvalidParam = 421

	// ErrAccount - No This User or Password Error
	ErrAccount = 422

	// ErrInternalServerError - Internal error.
	ErrInternalServerError = 500

	// ErrWechatPay - Wechat Pay error.
	ErrWechatPay = 520

	// ErrWechatAuth - Wechat Auth error.
	ErrWechatAuth = 521

	// ErrMongoDB - MongoDB operations error.
	ErrMongoDB = 600

	// ErrMysql - Mysql operations error.
	ErrMysql = 700
)
