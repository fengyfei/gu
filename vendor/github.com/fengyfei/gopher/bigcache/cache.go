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
 *     Initial: 2017/09/16        Jia Chenhui
 *     Modify : 2017/12/22        Yang Chenglong
 */

package cache

import (
	"github.com/allegro/bigcache"
)

type CacheDB struct {
	db     *bigcache.BigCache
}

func NewCacheDB(customConfig bigcache.Config) (*CacheDB, error) {
	db, err := bigcache.NewBigCache(customConfig)

	return &CacheDB{
		db:     db,
	}, err
}

func (c *CacheDB) Set(id string, name []byte) error {
	return c.db.Set(id, []byte(name))
}

func (c *CacheDB) Get(id string) ([]byte, error) {
	nameByte, err := c.db.Get(id)
	if err != nil {
		return nil, err
	}

	return nameByte, nil
}

func (c *CacheDB) Reset() error {
	return c.db.Reset()
}

func (c *CacheDB) Len() int {
	return c.db.Len()
}
