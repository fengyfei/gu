/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd.
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
 *     Initial: 2018/01/08        Jia Chenhui
 */

package grpc

import (
	"errors"

	gorpc "google.golang.org/grpc"
)

// Client defines an interface representing a gRPC-based client.
type Client interface {
	// Target returns the RPC server address.
	Target() string
	// Conn returns the client connection to an RPC server.
	Conn() *gorpc.ClientConn
	// Options returns a list of DialOption.
	Options() []gorpc.DialOption
	// Stop tears down the ClientConn and all underlying connections.
	Stop()
}

type clientImpl struct {
	target  string
	options []gorpc.DialOption
	conn    *gorpc.ClientConn
}

// Target - Client interface implementation.
func (c *clientImpl) Target() string {
	return c.target
}

// Conn - Client interface implementation.
func (c *clientImpl) Conn() *gorpc.ClientConn {
	return c.conn
}

// Options - Client interface implementation.
func (c *clientImpl) Options() []gorpc.DialOption {
	return c.options
}

// Stop - Client interface implementation.
func (c *clientImpl) Stop() {
	c.conn.Close()
}

// NewClientToTarget creates a new Client based on a target server address.
func NewClientToTarget(target string) (Client, error) {
	if target == "" {
		return nil, errors.New("empty target server")
	}

	var dialOpts []gorpc.DialOption

	dialOpts = append(dialOpts, gorpc.WithInsecure())
	// buffer size options
	dialOpts = append(dialOpts, gorpc.WithWriteBufferSize(serverMaxRecvMsgSize))
	// keepalive options
	dialOpts = append(dialOpts, clientKeepaliveOptions()...)

	conn, err := gorpc.Dial(target, dialOpts...)
	if err != nil {
		return nil, err
	}

	client := &clientImpl{
		target:  target,
		options: dialOpts,
		conn:    conn,
	}

	return client, nil
}
