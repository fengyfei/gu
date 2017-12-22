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
 *     Initial: 2017/09/17        Li Zebang
 */

package middleware

import (
	"net/http"
)

// BasicAuth checks the username and password provided in the request's Authorization header.
type BasicAuth struct {
	handler http.Handler
}

func (h *BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _, ok := r.BasicAuth()
	if ok {
		h.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

// SetCookie sets username provided in the request's Authorization header as the Cookie's Value.
func SetCookie(handler http.Handler) http.Handler {
	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
		})
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(nextHandler)
}
