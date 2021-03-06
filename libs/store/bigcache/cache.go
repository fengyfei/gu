/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/12/22        Yang Chenglong
 */

package cache

import (
	"github.com/allegro/bigcache"
)

// CacheDB represent a BigCache instance.
type CacheDB struct {
	db     *bigcache.BigCache
}

// NewCacheDB creates a new bigcache database.
func NewCacheDB(customConfig bigcache.Config) (*CacheDB, error) {
	db, err := bigcache.NewBigCache(customConfig)

	return &CacheDB{
		db:     db,
	}, err
}

// Set saves entry under the key.
func (c *CacheDB) Set(id string, name []byte) error {
	return c.db.Set(id, []byte(name))
}

// Get reads entry for the key.
func (c *CacheDB) Get(id string) ([]byte, error) {
	nameByte, err := c.db.Get(id)
	if err != nil {
		return nil, err
	}

	return nameByte, nil
}

// Reset empties all cache shards.
func (c *CacheDB) Reset() error {
	return c.db.Reset()
}

// Len computes number of entries in cache.
func (c *CacheDB) Len() int {
	return c.db.Len()
}
