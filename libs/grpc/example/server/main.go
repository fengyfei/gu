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

package main

import (
	"context"
	"log"

	"github.com/fengyfei/gu/libs/grpc"
	pb "github.com/fengyfei/gu/libs/grpc/example/greeter"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

var (
	counter = 0
)

type server struct{}

// SayHello implements greeter.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	counter++

	if (counter % 2) == 0 {
		// panic("greeter server panic")
	}

	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	s, err := grpc.NewServerFromAddress(port)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	pb.RegisterGreeterServer(s.Server(), &server{})
	reflection.Register(s.Server())

	log.Printf("Serving on: %s", port)
	s.Start()
}
