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
 *     Initial: 2017/12/19        Feng Yifei
 *     Modify:  2017/12/20        Jia Chenhui
 */

package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Router register routes to be matched and dispatched to a handler.
type Router struct {
	router     *mux.Router
	ctxPool    sync.Pool
	errHandler func(*Context)
	filters    [FinishRouter + 1][]*FilterRouter
}

// NewRouter returns a router.
func NewRouter() *Router {
	r := &Router{
		router:     mux.NewRouter(),
		errHandler: func(_ *Context) {},
	}

	r.ctxPool.New = func() interface{} {
		return NewContext(nil, nil)
	}

	r.router.NotFoundHandler = http.NotFoundHandler()
	r.router.MethodNotAllowedHandler = MethodNotAllowedHandler()

	return r
}

// Handler returns a http.Handler.
func (rt *Router) Handler() http.Handler {
	return rt.router
}

// SetErrorHandler attach a global error handler on router.
func (rt *Router) SetErrorHandler(h func(*Context)) {
	rt.errHandler = h
}

// Get adds a route path access via GET method.
func (rt *Router) Get(pattern string, handler HandlerFunc) {
	rt.router.HandleFunc(pattern, rt.wrapHandlerFunc(handler)).Methods("GET")
}

// Post adds a route path access via POST method.
func (rt *Router) Post(pattern string, handler HandlerFunc) {
	rt.router.HandleFunc(pattern, rt.wrapHandlerFunc(handler)).Methods("POST")
}

// Wraps a HandlerFunc to a http.HandlerFunc.
func (rt *Router) wrapHandlerFunc(f HandlerFunc) http.HandlerFunc {
	var (
		err error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		c := rt.ctxPool.Get().(*Context)
		defer rt.ctxPool.Put(c)
		c.Reset(r, w)

		if len(rt.filters[BeforeStatic]) > 0 && rt.execFilter(c, BeforeStatic) {
			err = errFilterNotPassed
			goto finish
		}

		if len(rt.filters[BeforeRouter]) > 0 && rt.execFilter(c, BeforeRouter) {
			err = errFilterNotPassed
			goto finish
		}

		if len(rt.filters[BeforeExec]) > 0 && rt.execFilter(c, BeforeExec) {
			err = errFilterNotPassed
			goto finish
		}

		// execute HandlerFunc
		if err = f(c); err != nil {
			goto finish
		}

		if len(rt.filters[AfterExec]) > 0 && rt.execFilter(c, AfterExec) {
			err = errFilterNotPassed
			goto finish
		}

		if len(rt.filters[FinishRouter]) > 0 && rt.execFilter(c, FinishRouter) {
			err = errFilterNotPassed
			goto finish
		}

	finish:
		c.LastError = err
		rt.errHandler(c)
	}
}

// MethodNotAllowedHandler returns a simple request handler
// that replies to each request with a ``405 method not allowed'' reply.
func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	})
}
