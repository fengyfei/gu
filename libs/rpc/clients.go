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

package rpc

import (
	"errors"
)

var (
	// ErrRPCNoClientAvailable represents no client left in collection.
	ErrRPCNoClientAvailable = errors.New("no rpc client available now")
)

// Clients represents rpc client collections.
type Clients struct {
	// rpc.Client is thread-safe, so Locker isn't required here
	clients []*Client
}

// Dials dial to rpc servers.
func Dials(opts []Options) *Clients {
	cs := new(Clients)

	for _, opt := range opts {
		cs.clients = append(cs.clients, Dial(opt))
	}

	return cs
}

// Get get a usable rpc client based on specified addr.
func (cs *Clients) Get(addr string) (*Client, error) {
	for {
		c, err := cs.get()
		if err != nil {
			return nil, err
		} else if c.options.Addr == addr {
			return c, nil
		}

		continue
	}
}

// get a usable rpc client.
func (cs *Clients) get() (*Client, error) {
	for _, c := range cs.clients {

		if c != nil && c.Client != nil && c.Error() == nil {
			return c, nil
		}
	}

	return nil, ErrRPCNoClientAvailable
}

// IsAvailable checks if exists a available client.
func (cs *Clients) IsAvailable() bool {
	_, err := cs.get()
	if err != nil {
		return true
	}

	return false
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (cs *Clients) Call(serviceMethod string, args interface{}, reply interface{}) error {
	var (
		err error
		c   *Client
	)

	if c, err = cs.get(); err == nil {
		err = c.Call(serviceMethod, args, reply)
	}

	return err
}

// Ping ping the rpc connect and reconnect when has an error.
func (cs *Clients) Ping(ping string) {
	for _, c := range cs.clients {
		go c.Ping(ping)
	}
}
