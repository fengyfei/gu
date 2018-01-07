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
 *     Initial: 2018/01/07        Feng Yifei
 */

package grpc

import (
	"errors"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	gorpc "google.golang.org/grpc"
)

// Server defines an interface representing a GRPC-based server.
type Server interface {
	// Address returns the listener address
	Address() string
	// Server returns the underlying grpc.Server instance
	Server() *gorpc.Server
	// Listener returns the underlying net.Listener instance
	Listener() net.Listener
	// Start the underlying grpc.Server
	Start() error
	// Stop the underlying grpc.Server
	Stop()
}

type serverImpl struct {
	address  string
	listener net.Listener
	server   *gorpc.Server
}

// Address - Server interface implementation.
func (s *serverImpl) Address() string {
	return s.address
}

// Listener - Server interface implementation.
func (s *serverImpl) Listener() net.Listener {
	return s.listener
}

// Server - Server interface implementation.
func (s *serverImpl) Server() *gorpc.Server {
	return s.server
}

// Start - Server interface implementation.
func (s *serverImpl) Start() error {
	return s.server.Serve(s.listener)
}

// Stop - Server interface implementation.
func (s *serverImpl) Stop() {
	s.server.Stop()
}

// NewServerFromAddress creates a new server listening on address.
func NewServerFromAddress(addr string) (Server, error) {
	if addr == "" {
		return nil, errors.New("empty server address")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewServerFromListener(lis)
}

// NewServerFromListener creates a new server based on a listener.
func NewServerFromListener(listener net.Listener) (Server, error) {
	var serverOpts []gorpc.ServerOption

	server := &serverImpl{
		address:  listener.Addr().String(),
		listener: listener,
	}

	// recovery middleware options
	serverOpts = append(serverOpts, gorpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		),
	))
	serverOpts = append(serverOpts, gorpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		),
	))

	// message size options
	serverOpts = append(serverOpts, gorpc.MaxRecvMsgSize(serverMaxRecvMsgSize))
	serverOpts = append(serverOpts, gorpc.MaxSendMsgSize(serverMaxSendMsgSize))

	// keepalive options
	serverOpts = append(serverOpts, serverKeepaliveOptions()...)

	server.server = gorpc.NewServer(serverOpts...)

	return server, nil
}
