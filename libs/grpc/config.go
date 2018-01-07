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
	"time"

	gorpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// KeepaliveOptions is the options for set the gRPC servers or clients.
type KeepaliveOptions struct {
	ClientKeepaliveTime    int
	ClientKeepaliveTimeout int
	ServerKeepaliveTime    int
	ServerKeepaliveTimeout int
}

func serverKeepaliveOptions() []gorpc.ServerOption {
	var serverOpts []gorpc.ServerOption

	kap := keepalive.ServerParameters{
		Time:    time.Duration(keepaliveOptions.ServerKeepaliveTime) * time.Second,
		Timeout: time.Duration(keepaliveOptions.ServerKeepaliveTimeout) * time.Second,
	}
	serverOpts = append(serverOpts, gorpc.KeepaliveParams(kap))

	kep := keepalive.EnforcementPolicy{
		MinTime:             time.Duration(keepaliveOptions.ClientKeepaliveTime) * time.Second,
		PermitWithoutStream: true,
	}
	serverOpts = append(serverOpts, gorpc.KeepaliveEnforcementPolicy(kep))
	return serverOpts
}
