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
 *     Initial: 2017/12/05        Jia Chenhui
 */

package server

import (
	"errors"

	"github.com/fengyfei/gu/libs/orm/mysql"
	"github.com/fengyfei/gu/libs/permission"
	"github.com/fengyfei/gu/libs/rpc"
	"github.com/fengyfei/gu/models/staff"
)

const (
	poolSize = 20
)

var (
	errPermissionNotMatch = errors.New("user permissions and url permissions do not match")
	errNoRole             = errors.New("this user or url has no role")
)

type AuthRPC struct {
	pool *mysql.Pool
}

func newAuthRPC(db string) *AuthRPC {
	ar := &AuthRPC{}
	ar.pool = mysql.NewPool(db, poolSize)

	if ar.pool == nil {
		panic("MySQL DB connection error.")
	}

	return ar
}

// Ping is general rpc keepalive interface.
func (ar *AuthRPC) Ping(req *rpc.ReqKeepAlive, resp *rpc.RespKeepAlive) error {
	return nil
}

// Verify check whether the user permissions match the URL permissions.
func (ar *AuthRPC) Verify(args permission.Args, reply *bool) error {
	err := ar.verifier(&(args.URL), args.UId)
	if err != nil {
		*reply = false
		return err
	}

	*reply = true
	return nil
}

func (ar *AuthRPC) verifier(url *string, uid int32) error {
	conn, err := ar.pool.Get()
	if err != nil {
		return err
	}
	defer ar.pool.Release(conn)

	urlRoles, err := staff.Service.URLPermissions(conn, url)
	if err != nil {
		return err
	}

	// If there is no permission record, pass the validation.
	if len(urlRoles) == 0 {
		return nil
	}

	userRoles, err := staff.Service.AssociatedRoleMap(conn, uid)
	if err != nil {
		return err
	}

	// If the user does not have a role, return the error directly.
	if len(userRoles) == 0 {
		return errNoRole
	}

	for urlR := range urlRoles {
		if userRoles[urlR] {
			return nil
		}
	}

	return errPermissionNotMatch
}
