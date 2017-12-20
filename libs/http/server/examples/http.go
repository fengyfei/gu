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

package main

import (
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/http/server/middleware"
	"github.com/fengyfei/gu/libs/logger"
)

type echo struct {
	Message *string `json:"message" validate:"required,min=6"`
}

func echoHandler(c *server.Context) error {
	var (
		err  error
		req  echo
		user string
	)

	header := c.GetHeader("auth")
	logger.Debug("header:", header)
	c.SetHeader("auth", "vbnmfs")

	user = c.FormValue("user")
	logger.Debug("user:", user)

	c.SetCookie("id", "123456")
	id, err := c.GetCookie("id")
	if err != nil {
		return err
	}
	logger.Debug("id:", id.Value)

	if err = c.JSONBody(&req); err != nil {
		logger.Error("Parses the JSON request body error:", err)
		return err
	}
	err = c.Validate(&req)
	if err != nil {
		logger.Error(err)
		return err
	}

	return c.ServeJSON(&req)
}

func postHandler(c *server.Context) error {
	// w.Write([]byte("Post\n"))
	return nil
}

func panicHandler(c *server.Context) error {
	panic("Panic testing")
	return nil
}

func main() {
	configuration := &server.Configuration{
		Address: "127.0.0.1:9573",
	}

	router := server.NewRouter()
	router.Post("/", echoHandler)
	router.Post("/post", postHandler)
	router.Get("/panic", panicHandler)

	ep := server.NewEntrypoint(configuration, nil)

	// add middlewares
	ep.AttachMiddleware(middleware.NegroniLoggerHandler())

	if err := ep.Start(router.Handler()); err != nil {
		logger.Error(err)
		return
	}

	ep.Wait()
}
