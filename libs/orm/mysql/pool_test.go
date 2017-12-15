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
 *     Initial: 2017/12/15        Feng Yifei
 */

package mysql

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

// TestNewPool -
func TestNewPool(t *testing.T) {
	db := "root:111111@tcp(127.0.0.1:3306)/mysql"

	pool := NewPool(db, 10)
	defer pool.Close()

	for i := 0; i < 100; i++ {
		go func() {
			conn, err := pool.Get()

			if conn == nil {
				t.Error(err)
				return
			}

			defer pool.Release(conn)

			result := conn.(*gorm.DB).Exec("SELECT * FROM user")
			if result.Error != nil {
				t.Error(result.Error)
			}
		}()
	}

	time.Sleep(2 * time.Second)

	for i := 0; i < 100; i++ {
		go func() {
			conn, err := pool.Get()

			if conn == nil {
				t.Error(err)
				return
			}

			defer pool.Release(conn)

			result := conn.(*gorm.DB).Exec("SELECT * FROM user")
			if result.Error != nil {
				t.Error(result.Error)
			}
		}()
	}

	time.Sleep(4 * time.Second)
}
