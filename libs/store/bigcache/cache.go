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
	"time"

	"github.com/allegro/bigcache"
)

type CacheServiceProvider struct{}

var (
	customConfig  bigcache.Config
	CacheInstance *bigcache.BigCache
	CacheServer   *CacheServiceProvider
)

func init() {
	customConfig = bigcache.Config{
		Shards:             1024,
		LifeWindow:         3 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
	}

	CacheInstance, _ = bigcache.NewBigCache(customConfig)
	CacheServer = &CacheServiceProvider{}
}

func (csp *CacheServiceProvider) SetOne(id string, name []byte) {
	CacheInstance.Set(id, []byte(name))
}

func (csp *CacheServiceProvider) GetOne(id string) (string, error) {
	nameByte, err := CacheInstance.Get(id)
	if err != nil {
		return "", err
	}

	return string(nameByte), nil
}
