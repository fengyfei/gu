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
 *     Initial: 2017/12/30        Jia Chenhui
 */

package crawler

import (
	"sync"

	"github.com/fengyfei/gu/libs/crawler/github"
)

// trendingCache use to store trending in specified language.
type trendingCache struct {
	mux   sync.RWMutex
	cache map[string][]*github.Trending
}

func newCache() *trendingCache {
	return &trendingCache{
		cache: make(map[string][]*github.Trending),
	}
}

// Store store the trending data in TrendingCache.
func (tc *trendingCache) Store(lang string, trending []*github.Trending) {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	tc.cache[lang] = trending
}

// Load getting the trending data of the specified language from TrendingCache.
func (tc *trendingCache) Load(lang string) ([]*github.Trending, bool) {
	tc.mux.RLock()
	defer tc.mux.RUnlock()

	t, ok := tc.cache[lang]
	return t, ok
}
