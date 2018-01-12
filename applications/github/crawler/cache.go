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
	"sort"
	"sync"

	"github.com/fengyfei/gu/libs/crawler/github"
)

// reposCache use to store the different repos of trending data, the key is
// "title + date".
type reposCache map[string]*github.Trending

// trendingCache use to store the trending data of specified language.
type trendingCache struct {
	mux   sync.RWMutex
	cache map[string]reposCache
}

func newTrendingCache() *trendingCache {
	return &trendingCache{
		cache: make(map[string]reposCache),
	}
}

// Store store the trending data in TrendingCache.
func (tc *trendingCache) Store(lang string, trending *github.Trending) {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	reposKey := trending.Title + trending.Date

	if _, ok := tc.cache[lang]; ok {
		tc.cache[lang][reposKey] = trending
	} else {
		tc.cache[lang] = make(reposCache)
		tc.cache[lang][reposKey] = trending
	}
}

// Load getting the trending list of the specified language from TrendingCache.
// The results were descended in sequence according to the field "Today" of
// struct "github.Trending".
func (tc *trendingCache) Load(lang string) ([]*github.Trending, bool) {
	var list []*github.Trending

	tc.mux.RLock()
	defer tc.mux.RUnlock()

	if reposMap, ok := tc.cache[lang]; ok {
		for reposKey := range reposMap {
			list = append(list, reposMap[reposKey])
		}

		sortByStar(list)
		return list, true
	}

	return nil, false
}

// Flush clears the cache for the specified language.
func (tc *trendingCache) Flush(lang string) {
	tc.mux.Lock()
	tc.cache[lang] = make(reposCache)
	tc.mux.Unlock()
}

func sortByStar(list []*github.Trending) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Today > list[j].Today
	})
}
