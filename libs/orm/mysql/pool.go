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
	_ "github.com/go-sql-driver/mysql"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
)

const (
	poolMaxSize = 100
	dialect     = "mysql"
)

// Pool represents the database connection pool.
type Pool struct {
	db *gorm.DB
}

// NewPool create a Pool according to the specified db and pool size.
func NewPool(db string, size int) *Pool {
	d, err := gorm.Open(dialect, db)

	if err != nil {
		panic(err)
	}

	pool := &Pool{
		db: d,
	}

	if size <= 0 {
		size = 20
	} else if size > poolMaxSize {
		size = poolMaxSize
	}

	pool.db.DB().SetMaxIdleConns(size)
	pool.db.DB().SetMaxOpenConns(size << 1)

	return pool
}

// Get get a connection from the pool.
func (p *Pool) Get() (orm.Connection, error) {
	return p.db, nil
}

// Release put the connection back into the pool.
func (p *Pool) Release(v orm.Connection) {
	return
}

// Close close the pool.
func (p *Pool) Close() {
	p.db.Close()
}
