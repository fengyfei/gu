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
	"errors"
	"net"
	"sync"

	"github.com/fengyfei/gu/libs/logger"
)

var (
	errIDNotExists = errors.New("this id does not exist")
)

// Server is a general tcp server.
type Server struct {
	mux       sync.RWMutex
	listener  net.Listener
	onlineMap map[string]bool
}

// NewServer creates a Server.
func NewServer(address string) (*Server, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener:  l,
		onlineMap: make(map[string]bool),
	}

	return server, nil
}

// IsOnline checks if the specified user id is online.
func (s *Server) IsOnline(id string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if !s.onlineMap[id] {
		return false
	}

	return true
}

// Add sets the user of the specified id to be online.
func (s *Server) Add(id string) bool {
	if id == "" {
		logger.Debug("Add:", errIDNotExists)
		return false
	}

	s.mux.Lock()
	s.onlineMap[id] = true
	s.mux.Unlock()

	return true
}

// Remove removes the user with the specified id from the online users.
func (s *Server) Remove(id string) {
	s.mux.Lock()
	delete(s.onlineMap, id)
	s.mux.Unlock()
}

// Counts represents the number of online users.
func (s *Server) Counts() int32 {
	s.mux.RLock()
	c := len(s.onlineMap)
	s.mux.RUnlock()

	return int32(c)
}
