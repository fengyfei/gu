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
 *     Initial: 2017/09/27        Jia Chenhui
 */

package schema

import (
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"

	"github.com/fengyfei/gopher/graphql/user/module"
	"github.com/fengyfei/gopher/graphql/user/util"
)

var UserSchema graphql.Schema

func init() {
	s, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	UserSchema = s
}

// ExecQuery execute the query.
func ExecQuery(query string, schema graphql.Schema) (*graphql.Result, error) {
	p := graphql.Params{
		Schema:        schema,
		RequestString: query,
	}

	result := graphql.Do(p)
	if result.HasErrors() {
		log.Printf("Wrong result, unexpected errors: %v", result.Errors)
		err := errors.New(fmt.Sprintf("%v", result.Errors))
		return nil, err
	}

	return result, nil
}

var (
	rootQuery        = graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation     = graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}
	rootSubscription = graphql.ObjectConfig{Name: "RootSubscription", Fields: subscription}

	schemaConfig = graphql.SchemaConfig{
		Query:        graphql.NewObject(rootQuery),
		Mutation:     graphql.NewObject(rootMutation),
		Subscription: graphql.NewObject(rootSubscription),
	}
)

var (
	// User data structure
	userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: graphql.String,
			},
			"admin": &graphql.Field{
				Type: graphql.String,
			},
			"active": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
)

var (
	// query data
	// get: curl -g 'http://localhost:8989/graphql?query={getUser(login:"jch"){login,admin,active}}'
	// An example GraphQL query might look like:
	/*
		{
		  getUser(login: "leon") {
			login, admin, active
		  }
		}
	*/
	fields = graphql.Fields{
		"getUser": &graphql.Field{
			Type:        userType,
			Description: "Get single user info",
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: Get,
		},
	}

	// mutation data
	// create: curl -g 'http://localhost:8989/graphql?query=mutation+_{createUser(login:"jch",admin:"yes",active:"yes"){login,admin,active}}'
	// An example GraphQL mutation might look like:
	/*
		mutation {
		  createUser(login: "leon", admin: "true", active: "true") {
		    login, admin, active
		  }
		}
	*/
	mutations = graphql.Fields{
		"createUser": &graphql.Field{
			Type:        userType,
			Description: "Create a new user",
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"admin": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"active": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: Create,
		},
	}

	// subscription to user information
	/*
		subscription {
		  subscribeUser(postId:"a"){
		  	admin, active
		  }
		}
	*/
	subscription = graphql.Fields{
		"subscribeUser": &graphql.Field{
			Type:        userType,
			Description: "Subscribe to user information",
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return util.RandSeq(10), nil
			},
		},
	}
)

// GetSingle get single user information.
func Get(p graphql.ResolveParams) (interface{}, error) {
	login := p.Args["login"].(string)

	user, err := module.UserService.Get(login)
	if err != nil {
		log.Printf("Get user returned error: %v", err)
		return nil, err
	}

	return user, nil
}

// Create create single user.
func Create(p graphql.ResolveParams) (interface{}, error) {
	user := module.User{
		Login:  p.Args["login"].(string),
		Admin:  p.Args["admin"].(string),
		Active: p.Args["active"].(string),
	}

	ok, err := module.UserService.Create(&user)
	if err != nil {
		log.Printf("Create user returned error: %v", err)
		return nil, err
	}

	return ok, nil
}
