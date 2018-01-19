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
 *     Initial: 2018/01/08        Feng Yifei
 */

package tcp

import (
	"net"
	"sync"

	"github.com/fengyfei/gu/libs/logger"
)

// Server is a general tcp server.
type Server struct {
	options  *Options
	conns    *sync.Map
	listener net.Listener
	buffer   []*Message
	sender   chan *Message
	mux      sync.RWMutex
	stop     chan struct{}
}

// NewServer creates a Server.
func NewServer(opts *Options) *Server {
	return &Server{
		options: opts,
		conns:   &sync.Map{},
		buffer:  make([]*Message, opts.ReadBufSize),
		sender:  make(chan *Message),
		stop:    make(chan struct{}),
	}
}

// ListenAndServe listen on the TCP network and then calls Serve to handle
// requests on incoming connections.
func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.options.Address+":"+s.options.Port)
	if err != nil {
		return err
	}

	s.Serve(l)
	return nil
}

// Serve start the TCP server to accept.
func (s *Server) Serve(listener net.Listener) {
	defer func() {
		s.mux.Lock()
		lis := s.listener
		s.listener = nil
		s.mux.Unlock()

		if lis != nil {
			lis.Close()
		}
	}()

	s.mux.Lock()
	s.listener = listener
	s.mux.Unlock()

	logger.Debug("TCP server listen on:", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		select {
		case <-s.stop:
			return
		default:
		}

		go s.handlerConn(conn)
	}
}

// handlerConn handle net.Conn.
func (s *Server) handlerConn(conn net.Conn) {
	go s.receive(conn)

	ipStr := conn.RemoteAddr().String()
	s.mux.Lock()
	if s.conns == nil {
		s.mux.Unlock()
		conn.Close()
		return
	}
	s.mux.Unlock()

	if !s.Add(ipStr) {
		conn.Close()
		return
	}
	defer s.Remove(ipStr)

	for {
		select {
		case <-s.stop:
			conn.Close()
			return
		case msg := <-s.sender:
			err := msg.Write(conn)
			if err != nil {
				s.options.Handler.OnError(err)
			}
		default:
		}
	}
}

// receive receive message from conn.
func (s *Server) receive(conn net.Conn) {
	index := 0
	capacity := s.options.ReadBufSize

	for {
		msg := s.buffer[index]
		msg.Reset()

		err := msg.Read(conn)
		if err != nil {
			s.options.Handler.OnError(err)
		} else {
			s.options.Handler.OnMessage(msg)
		}

		index = (index + 1) % capacity
	}
}

// Send send a message to conn.
func (s *Server) Send(payload []byte, conn net.Conn) error {
	msg, err := decode(payload)
	if err != nil {
		return err
	}

	return msg.Write(conn)
}

// Listener returns the net.Listener of Server.
func (s *Server) Listener() net.Listener {
	return s.listener
}

// IsOnline checks if the specified user id is online.
func (s *Server) IsOnline(id string) bool {
	if _, ok := s.conns.Load(id); !ok {
		return false
	}

	return true
}

// Add sets the user of the specified id to be online.
func (s *Server) Add(id string) bool {
	if id == "" {
		return false
	}
	s.conns.Store(id, true)
	return true
}

// Remove removes the user with the specified id from the online users.
func (s *Server) Remove(id string) {
	s.conns.Delete(id)
}

// Counts represents the number of online users.
func (s *Server) Counts() int32 {
	var size int32
	s.conns.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return size
}
