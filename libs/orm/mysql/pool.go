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
 *     Initial: 2017/11/02        Jia Chenhui
 *     Modify : 2017/11/04        Yang Chenglong
 */

package mysql

import (
	"container/ring"
	"errors"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
)

const (
	poolMaxSize = 200
	dialect     = "mysql"
)

var (
	ErrNoConnection = errors.New("MySQL Connection expired")
)

// Pool represents the database connection pool.
type Pool struct {
	db      string
	maxSize int
	lock    sync.Mutex
	pool    *ring.Ring
	size    int
}

// NewPool create a Pool according to the specified db and pool size.
func NewPool(db string, size int) *Pool {
	var (
		err  error
		conn *ring.Ring
	)

	if size <= 0 {
		size = 20
	}

	if size > poolMaxSize {
		size = poolMaxSize
	}

	pool := &Pool{
		db:      db,
		maxSize: size,
	}

	for i := 0; i < size; i++ {
		conn = ring.New(1)
		conn.Value, err = gorm.Open(dialect, db)

		if err != nil {
			logger.Error(err)
			continue
		}

		if pool.pool == nil {
			pool.pool = conn
		}

		pool.pool.Link(conn)
	}

	pool.size = pool.pool.Len()
	if pool.size != size {
		logger.Debug("New pool not enough!")
	}

	logger.Debug("Pool size:", pool.size)

	return pool
}

// Get get a connection from the pool.
func (p *Pool) Get() (orm.Connection, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.size == 0 {
		conn, err := gorm.Open(dialect, p.db)

		if err != nil {
			return nil, ErrNoConnection
		}
		return conn, nil
	}

	p.size -= 1

	conn := p.pool.Unlink(1)
	return conn.Value.(orm.Connection), nil
}

// Release put the connection back into the pool.
func (p *Pool) Release(v orm.Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.size == p.maxSize {
		v.(*gorm.DB).Close()
		return
	}

	conn := ring.New(1)
	conn.Value = v

	p.size += 1
	p.pool.Prev().Link(conn)
}

// Close close the pool.
func (p *Pool) Close() {
	f := func(v interface{}) {
		if v == nil {
			return
		}

		conn := v.(*gorm.DB)
		conn.Close()
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.size = 0
	p.pool.Do(f)
	p.pool = nil
}
