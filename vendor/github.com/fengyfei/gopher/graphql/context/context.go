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
 *     Initial: 2017/09/26        Yang Chenglong
 */

package main

import (
	"context"
	"log"
	"net/http"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	Schema graphql.Schema
	serve  *echo.Echo

	userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"me": &graphql.Field{
					Type: userType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return p.Context.Value("currentUser"), nil
					},
				},
			},
		})
)

func GraphQLHandler(c echo.Context) error {
	user := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{1, "cool user"}

	result := graphql.Do(graphql.Params{
		Schema:        Schema,
		RequestString: c.Request().URL.Query().Get("query"),
		Context:       context.WithValue(context.Background(), "currentUser", user),
	})
	if result.HasErrors() {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
		err := errors.New(fmt.Sprintf("%v", result.Errors))
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func init() {
	serve = echo.New()

	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		log.Fatalf("failed to create schema, error: %v", err)
	}
	Schema = s
}

func main() {
	serve.GET("/graphql", GraphQLHandler)
	serve.Use(middleware.Logger())

	serve.Start("localhost:8080")
}
