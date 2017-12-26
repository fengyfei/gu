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
 *     From:   https://github.com/rs/cors/blob/master/cors.go
 *     Modify: 2017/12/21        Jia Chenhui
 */

package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fengyfei/gu/libs/http/server"
)

// CORSConfig is a configuration container to setup the CORS middleware.
type CORSConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
}

// CORS http handler
type CORS struct {
	// skipper defines a function to skip middleware.
	skipper Skipper
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard
	// Optional origin validator function
	allowOriginFunc func(origin string) bool
	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	// Normalized list of allowed headers
	allowedHeaders []string
	// Normalized list of allowed methods
	allowedMethods []string
	// Normalized list of exposed headers
	exposedHeaders    []string
	allowCredentials  bool
	maxAge            int
	optionPassthrough bool
}

// CORSWithConfig creates a new CORS handler with the provided config.
func CORSWithConfig(config CORSConfig) *CORS {
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}
	c := &CORS{
		skipper:           config.Skipper,
		exposedHeaders:    convert(config.ExposedHeaders, http.CanonicalHeaderKey),
		allowOriginFunc:   config.AllowOriginFunc,
		allowCredentials:  config.AllowCredentials,
		maxAge:            config.MaxAge,
		optionPassthrough: config.OptionsPassthrough,
	}

	// Normalize options
	// Note: for origins and methods matching, the spec requires a case-sensitive matching.
	// As it may error prone, we chose to ignore the spec here.

	// Allowed Origins
	if len(config.AllowedOrigins) == 0 {
		if config.AllowOriginFunc == nil {
			// Default is all origins
			c.allowedOriginsAll = true
		}
	} else {
		c.allowedOrigins = []string{}
		c.allowedWOrigins = []wildcard{}
		for _, origin := range config.AllowedOrigins {
			// Normalize
			origin = strings.ToLower(origin)
			if origin == "*" {
				// If "*" is present in the list, turn the whole list into a match all
				c.allowedOriginsAll = true
				c.allowedOrigins = nil
				c.allowedWOrigins = nil
				break
			} else if i := strings.IndexByte(origin, '*'); i >= 0 {
				// Split the origin in two: start and end string without the *
				w := wildcard{origin[0:i], origin[i+1:]}
				c.allowedWOrigins = append(c.allowedWOrigins, w)
			} else {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	}

	// Allowed Headers
	if len(config.AllowedHeaders) == 0 {
		// Use sensible defaults
		c.allowedHeaders = []string{"Origin", "Accept", "Content-Type"}
	} else {
		// Origin is always appended as some browsers will always request for this header at preflight
		c.allowedHeaders = convert(append(config.AllowedHeaders, "Origin"), http.CanonicalHeaderKey)
		for _, h := range config.AllowedHeaders {
			if h == "*" {
				c.allowedHeadersAll = true
				c.allowedHeaders = nil
				break
			}
		}
	}

	// Allowed Methods
	if len(config.AllowedMethods) == 0 {
		// Default is spec's "simple" methods
		c.allowedMethods = []string{"GET", "POST", "HEAD"}
	} else {
		c.allowedMethods = convert(config.AllowedMethods, strings.ToUpper)
	}

	return c
}

// CORSDefault creates a new CORS handler with default config.
func CORSDefault() *CORS {
	return CORSWithConfig(CORSConfig{})
}

// CORSAllowAll create a new CORS handler with permissive configuration allowing all
// origins with all standard methods with any header and credentials.
func CORSAllowAll() *CORS {
	return CORSWithConfig(CORSConfig{
		Skipper:          DefaultSkipper,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{server.HEAD, server.GET, server.POST, server.PUT, server.PATCH, server.DELETE},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
}

// Handler apply the CORS specification on the request, and add relevant CORS headers
// as necessary.
func (c *CORS) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == server.OPTIONS && r.Header.Get(server.HeaderAccessControlRequestMethod) != "" {
			c.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			if c.optionPassthrough {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			c.handleActualRequest(w, r)
			h.ServeHTTP(w, r)
		}
	})
}

// HandlerFunc provides Martini compatible handler.
func (c *CORS) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == server.OPTIONS && r.Header.Get(server.HeaderAccessControlRequestMethod) != "" {
		c.handlePreflight(w, r)
	} else {
		c.handleActualRequest(w, r)
	}
}

// Negroni compatible interface.
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := server.NewContext(w, r)
	if c.skipper(ctx) {
		next(w, r)
		return
	}

	if r.Method == server.OPTIONS && r.Header.Get(server.HeaderAccessControlRequestMethod) != "" {
		c.handlePreflight(w, r)
		// Preflight requests are standalone and should stop the chain as some other
		// middleware may not handle OPTIONS requests correctly. One typical example
		// is authentication middleware ; OPTIONS requests won't carry authentication
		// headers (see #1)
		if c.optionPassthrough {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		c.handleActualRequest(w, r)
		next(w, r)
	}
}

// handlePreflight handles pre-flight CORS requests.
func (c *CORS) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get(server.HeaderOrigin)

	if r.Method != server.OPTIONS {
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/CORS/issues/10,
	//     https://github.com/rs/CORS/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add(server.HeaderVary, server.HeaderOrigin)
	headers.Add(server.HeaderVary, server.HeaderAccessControlRequestMethod)
	headers.Add(server.HeaderVary, server.HeaderAccessControlRequestHeaders)

	if origin == "" {
		return
	}
	if !c.isOriginAllowed(origin) {
		return
	}

	reqMethod := r.Header.Get(server.HeaderAccessControlRequestMethod)
	if !c.isMethodAllowed(reqMethod) {
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get(server.HeaderAccessControlRequestHeaders))
	if !c.areHeadersAllowed(reqHeaders) {
		return
	}
	if c.allowedOriginsAll && !c.allowCredentials {
		headers.Set(server.HeaderAccessControlAllowOrigin, "*")
	} else {
		headers.Set(server.HeaderAccessControlAllowOrigin, origin)
	}
	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set(server.HeaderAccessControlAllowMethods, strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {

		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers.Set(server.HeaderAccessControlAllowHeaders, strings.Join(reqHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set(server.HeaderAccessControlAllowCredentials, "true")
	}
	if c.maxAge > 0 {
		headers.Set(server.HeaderAccessControlMaxAge, strconv.Itoa(c.maxAge))
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects.
func (c *CORS) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get(server.HeaderOrigin)

	if r.Method == server.OPTIONS {
		return
	}
	// Always set Vary, see https://github.com/rs/CORS/issues/10
	headers.Add(server.HeaderVary, server.HeaderOrigin)
	if origin == "" {
		return
	}
	if !c.isOriginAllowed(origin) {
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !c.isMethodAllowed(r.Method) {
		return
	}
	if c.allowedOriginsAll && !c.allowCredentials {
		headers.Set(server.HeaderAccessControlAllowOrigin, "*")
	} else {
		headers.Set(server.HeaderAccessControlAllowOrigin, origin)
	}
	if len(c.exposedHeaders) > 0 {
		headers.Set(server.HeaderAccessControlExposeHeaders, strings.Join(c.exposedHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set(server.HeaderAccessControlAllowCredentials, "true")
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint.
func (c *CORS) isOriginAllowed(origin string) bool {
	if c.allowOriginFunc != nil {
		return c.allowOriginFunc(origin)
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range c.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

// isMethodAllowed checks if a given method can be used as part of a cross-domain request
// on the endpoing.
func (c *CORS) isMethodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == server.OPTIONS {
		// Always allow preflight requests
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (c *CORS) areHeadersAllowed(requestedHeaders []string) bool {
	if c.allowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
