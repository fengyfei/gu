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
	errEmptyParam   = errors.New("no id or value")
	errEmptyConnMap = errors.New("empty connections map")
)

// Session represents a set of connections.
type Session struct {
	mux     sync.RWMutex
	ID      string
	connMap map[string]interface{}
}

// NewSession creates a Session.
func NewSession(id string) *Session {
	return &Session{
		ID:      id,
		connMap: make(map[string]interface{}),
	}
}

// SessionID returns the session ID.
func (s *Session) SessionID() string {
	return s.ID
}

// Contains checks whether the specified id is included in this session.
func (s *Session) Contains(id string) bool {
	if s.connMap == nil {
		logger.Error(errEmptyConnMap)
		return false
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if _, ok := s.connMap[id]; !ok {
		return false
	}

	return true
}

// Add store id and v in this session.
func (s *Session) Add(id string, v interface{}) error {
	if id == "" || v == nil {
		return errEmptyParam
	}

	if s.connMap == nil {
		return errEmptyConnMap
	}

	s.mux.Lock()
	s.connMap[id] = v
	s.mux.Unlock()

	return nil
}

// Get gets the value of the specified id.
func (s *Session) Get(id string) interface{} {
	if s.connMap == nil {
		logger.Error(errEmptyConnMap)
		return nil
	}

	s.mux.RLock()
	v := s.connMap[id]
	s.mux.RUnlock()

	return v
}

// Remove removes the specified id from this session.
func (s *Session) Remove(id string) {
	if s.connMap == nil {
		logger.Error(errEmptyConnMap)
		return
	}

	s.mux.Lock()
	delete(s.connMap, id)
	s.mux.Unlock()
}

// Broadcast send information to the connection in this session.
func (s *Session) Broadcast(message []byte) error {
	if s.connMap == nil {
		return errEmptyConnMap
	}

	for _, v := range s.connMap {
		conn, ok := v.(net.Conn)
		if !ok {
			continue
		}

		conn.Write(message)
	}

	return nil
}

// Counts returns the number of connections included in this session.
func (s *Session) Counts() int32 {
	if s.connMap == nil {
		logger.Error(errEmptyConnMap)
		return int32(0)
	}

	s.mux.RLock()
	c := len(s.connMap)
	s.mux.RUnlock()

	return int32(c)
}

// Destroy empty this session.
func (s *Session) Destroy() {
	if s.connMap == nil {
		logger.Error(errEmptyConnMap)
		return
	}

	s.mux.Lock()
	s.connMap = make(map[string]interface{})
	s.mux.Unlock()
}
