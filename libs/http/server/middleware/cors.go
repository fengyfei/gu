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
 *     From: https://github.com/labstack/echo/blob/master/middleware/cors.go
 *     Modify:  2017/12/19        Jia Chenhui
 */

package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fengyfei/gu/libs/http/server"
)

// CORSConfig defines the config for CORS middleware.
type CORSConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// AllowOrigin defines a list of origins that may access the resource.
	// Optional. Default value []string{"*"}.
	AllowOrigins []string `json:"allow_origins"`

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	// Optional. Default value DefaultCORSConfig.AllowMethods.
	AllowMethods []string `json:"allow_methods"`

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This in response to a preflight request.
	// Optional. Default value []string{}.
	AllowHeaders []string `json:"allow_headers"`

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	// Optional. Default value false.
	AllowCredentials bool `json:"allow_credentials"`

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Optional. Default value []string{}.
	ExposeHeaders []string `json:"expose_headers"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	// Optional. Default value 0.
	MaxAge int `json:"max_age"`
}

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:      DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{server.GET, server.POST},
	}
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func CORS() server.MiddlewareFunc {
	return CORSWithConfig(DefaultCORSConfig)
}

// CORSWithConfig returns a CORS middleware with config.
func CORSWithConfig(config CORSConfig) server.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCORSConfig.Skipper
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(c *server.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			resp := c.Response()
			origin := req.Header.Get(server.HeaderOrigin)
			allowOrigin := ""

			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
			}

			// Simple request
			if req.Method != server.OPTIONS {
				resp.Header().Add(server.HeaderVary, server.HeaderOrigin)
				resp.Header().Set(server.HeaderAccessControlAllowOrigin, allowOrigin)
				if config.AllowCredentials {
					resp.Header().Set(server.HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					resp.Header().Set(server.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				return next(c)
			}

			// Preflight request
			resp.Header().Add(server.HeaderVary, server.HeaderOrigin)
			resp.Header().Add(server.HeaderVary, server.HeaderAccessControlRequestMethod)
			resp.Header().Add(server.HeaderVary, server.HeaderAccessControlRequestHeaders)
			resp.Header().Set(server.HeaderAccessControlAllowOrigin, allowOrigin)
			resp.Header().Set(server.HeaderAccessControlRequestMethod, allowMethods)
			if config.AllowCredentials {
				resp.Header().Set(server.HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				resp.Header().Set(server.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := req.Header.Get(server.HeaderAccessControlRequestHeaders)
				if h != "" {
					resp.Header().Set(server.HeaderAccessControlAllowHeaders, h)
				}
			}

			if config.MaxAge > 0 {
				resp.Header().Set(server.HeaderAccessControlMaxAge, maxAge)
			}
			return c.WriteHeader(http.StatusNoContent)
		}
	}
}
