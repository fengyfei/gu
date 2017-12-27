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
 *     From:   https://github.com/labstack/echo/blob/master/middleware/jwt.go
 *     Modify: 2017/12/25        Jia Chenhui
 */

package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/fengyfei/gu/libs/http/server"
)

// Algorithms
const (
	AlgorithmHS256 = "HS256"
)

// JWTConfig is a configuration container to setup the JWT middleware.
type JWTConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Signing key to validate token.
	// Required.
	// Note: if SigningMethod is AlgorithmHS256, SigningKey type must be []byte.
	SigningKey interface{}

	// Signing method, used to check token signing method.
	// Optional. Default value HS256.
	SigningMethod string

	// Context key to store user information from the token into context.
	// Optional. Default value "user".
	ContextKey string

	// Claims are extendable claims data defining token content.
	// Optional. Default value jwt.MapClaims
	Claims jwt.Claims

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default value "Bearer".
	AuthScheme string

	keyFunc jwt.Keyfunc
}

// JWT represents a JWT HTTP handler.
type JWT struct {
	skipper       Skipper
	signingKey    interface{}
	signingMethod string
	contextKey    string
	claims        jwt.Claims
	tokenLookup   string
	authScheme    string
	keyFunc       jwt.Keyfunc
}

type jwtExtractor func(*server.Context) (string, error)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:       DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + server.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
	}
)

// JWTWithKey returns a JWT with the specified key and other default configuration.
func JWTWithKey(key interface{}) *JWT {
	config := DefaultJWTConfig
	config.SigningKey = key
	return JWTWithConfig(config)
}

// JWTWithConfig returns a JWT with specified config.
func JWTWithConfig(config JWTConfig) *JWT {
	// defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil {
		panic("jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.SigningKey, nil
	}

	return &JWT{
		skipper:       config.Skipper,
		signingKey:    config.SigningKey,
		signingMethod: config.SigningMethod,
		contextKey:    config.ContextKey,
		claims:        config.Claims,
		tokenLookup:   config.TokenLookup,
		authScheme:    config.AuthScheme,
		keyFunc:       config.keyFunc,
	}
}

// ServeHTTP implements negroni compatible interface.
func (j *JWT) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// execute skipper
	ctx := server.NewContext(w, r)
	if j.skipper(ctx) {
		next(w, r)
		return
	}

	// extract token string from request header/query/cookie
	parts := strings.Split(j.tokenLookup, ":")
	extractor := jwtFromHeader(parts[1], j.authScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}

	auth, err := extractor(ctx)
	if err != nil {
		ctx.WriteHeader(http.StatusForbidden)
		return
	}
	token := new(jwt.Token)

	// validate, and return a token
	if _, ok := j.claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(auth, j.keyFunc)
	} else {
		claims := reflect.ValueOf(j.claims).Interface().(jwt.Claims)
		token, err = jwt.ParseWithClaims(auth, claims, j.keyFunc)
	}

	if err == nil && token.Valid {
		// store user information from token into context
		ctxWithClaims := context.WithValue(context.Background(), j.contextKey, token.Claims)
		r = r.WithContext(ctxWithClaims)
		next(w, r)
		return
	}

	ctx.WriteHeader(http.StatusForbidden)
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c *server.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or invalid jwt in the request header")
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c *server.Context) (string, error) {
		token := c.FormValue(param)
		if token == "" {
			return "", errors.New("missing jwt in the query string")
		}
		return token, nil
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from the named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(c *server.Context) (string, error) {
		cookie, err := c.GetCookie(name)
		if err != nil {
			return "", errors.New("missing jwt in the cookie")
		}
		return cookie.Value, nil
	}
}
