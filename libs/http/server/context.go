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
 */

package server

import (
	"errors"
	"net/http"

	json "github.com/json-iterator/go"
)

var (
	errNoBody        = errors.New("Request body is empty")
	errEmptyResponse = errors.New("Empty JSON response body")
)

// Context wraps http.Request and http.ResponseWriter.
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	LastError      error
}

// NewContext create a new context.
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
	}
}

// Reset the context.
func (c *Context) Reset(r *http.Request, w http.ResponseWriter) {
	c.request = r
	c.responseWriter = w
	c.LastError = nil
}

// JSONBody parses the JSON request body.
func (c *Context) JSONBody(v interface{}) error {
	if c.request.ContentLength == 0 {
		return errNoBody
	}

	return json.NewDecoder(c.request.Body).Decode(v)
}

// ServeJSON sends a JSON response.
func (c *Context) ServeJSON(v interface{}) error {
	if v == nil {
		return errEmptyResponse
	}

	resp, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.responseWriter.Write(resp)
	return err
}
