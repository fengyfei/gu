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
 *     Modify:  2017/12/23        Jia Chenhui
 */

package server

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	json "github.com/json-iterator/go"
	"gopkg.in/go-playground/validator.v9"
)

const defaultMemory = 32 << 20 // 32 MB

var (
	errNoBody              = errors.New("request body is empty")
	errEmptyResponse       = errors.New("empty JSON response body")
	errNotJsonBody         = errors.New("request body is not JSON")
	errInvalidRedirectCode = errors.New("invalid redirect status code")
)

// Context wraps http.Request and http.ResponseWriter.
// Returned set to true if response was written to then don't execute other handler.
type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	Validator      *validator.Validate
	LastError      error
}

// NewContext create a new context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		responseWriter: w,
		request:        r,
		Validator:      validator.New(),
	}
}

// Reset the context.
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.responseWriter = w
	c.request = r
	c.LastError = nil
}

// JSONBody parses the JSON request body.
func (c *Context) JSONBody(v interface{}) error {
	if c.request.ContentLength == 0 {
		return errNoBody
	}

	if conType := c.request.Header.Get(HeaderContentType); !strings.EqualFold(conType, MIMEApplicationJSON) {
		return errNotJsonBody
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

	c.responseWriter.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	_, err = c.responseWriter.Write(resp)
	return err
}

// Redirect does redirection to localurl with http header status code.
func (c *Context) Redirect(status int, url string) error {
	if status < 300 || status > 308 {
		return errInvalidRedirectCode
	}
	c.responseWriter.Header().Set(HeaderLocation, url)
	c.responseWriter.WriteHeader(status)
	return nil
}

// Cookies return all cookies.
func (c *Context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

// GetCookie Get cookie from request by a given key.
func (c *Context) GetCookie(key string) (*http.Cookie, error) {
	return c.request.Cookie(key)
}

// SetCookie Set cookie for response.
func (c *Context) SetCookie(name string, value string) {
	cook := http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(c.responseWriter, &cook)
}

// Request return the request.
func (c *Context) Request() *http.Request {
	return c.request
}

// SetRequest set the r as the new request.
func (c *Context) SetRequest(r *http.Request) {
	c.request = r
}

// Response return the responseWriter
func (c *Context) Response() http.ResponseWriter {
	return c.responseWriter
}

// FormValue returns the first value for the named component of the query.
func (c *Context) FormValue(name string) string {
	return c.request.FormValue(name)
}

// FormParams return the parsed form data
func (c *Context) FormParams() (url.Values, error) {
	if strings.HasPrefix(c.request.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := c.request.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := c.request.ParseForm(); err != nil {
			return nil, err
		}
	}
	return c.request.Form, nil
}

// SetHeader Set header for response.
func (c *Context) SetHeader(key, val string) {
	c.responseWriter.Header().Set(key, val)
}

// WriteHeader sends an HTTP response header with status code.
func (c *Context) WriteHeader(code int) error {
	c.responseWriter.WriteHeader(code)
	return nil
}

// GetHeader Get header from request by a given key.
func (c *Context) GetHeader(key string) string {
	return c.request.Header.Get(key)
}

// Validate the parameters.
func (c *Context) Validate(val interface{}) error {
	return c.Validator.Struct(val)
}
