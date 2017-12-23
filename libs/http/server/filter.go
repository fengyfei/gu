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
 *     From: https://github.com/astaxie/beego/blob/master/router.go
 *     Modify: 2017/12/22        Jia Chenhui
 */

package server

import (
	"errors"
)

// default filter execution points.
const (
	BeforeStatic = iota
	BeforeRouter
	BeforeExec
	AfterExec
	FinishRouter
)

var (
	errNotRightFilterPosition = errors.New("can not find this filter position")
	errFilterNotPassed        = errors.New("filter check is not passed")
)

// FilterFunc defines a filter function which is invoked before the controller
// handler is executed.
type FilterFunc func(*Context)

// FilterRouter defines a filter operation which is invoked before the
// controller handler is executed.
type FilterRouter struct {
	filterFunc FilterFunc
	pattern    string
}

// InsertFilter Add a FilterFunc with pattern rule and action constant.
func (rt *Router) InsertFilter(pattern string, pos int, filter FilterFunc) error {
	fr := &FilterRouter{
		filterFunc: filter,
		pattern:    pattern,
	}

	if pos < BeforeStatic || pos > FinishRouter {
		return errNotRightFilterPosition
	}

	rt.filters[pos] = append(rt.filters[pos], fr)

	return nil
}

// execFilter execution filter. If the filter is passed, return false, else true.
func (rt *Router) execFilter(ctx *Context, pos int) (returned bool) {
	for _, fr := range rt.filters[pos] {
		if ctx.Returned {
			return true
		}

		fr.filterFunc(ctx)

		// check whether the filterFunc changes ctx.Returned
		if ctx.Returned {
			return true
		}
	}

	return false
}
