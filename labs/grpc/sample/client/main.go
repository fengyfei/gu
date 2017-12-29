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
 *     Initial: 2017/12/29        Feng Yifei
 */

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/fengyfei/gu/labs/grpc/sample/greeter"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	n       = 50
)

func main() {
	conns := make([]*grpc.ClientConn, n)
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithWriteBufferSize(1*1024*1024))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		conns[i] = conn
	}

	start := time.Now().UnixNano()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		name := fmt.Sprintf("Go %d", i+1)
		go func() {
			conn := conns[i%n]
			c := pb.NewGreeterClient(conn)
			r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
			if err != nil {
				log.Printf("could not greet: %v", err)
			} else {
				log.Printf("Greeting: %s", r.Message)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	end := time.Now().UnixNano()

	fmt.Printf("Execution time: %f\n", float64(end-start)/1000000000)
}
