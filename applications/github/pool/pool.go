/*
 * Revision History:
 *     Initial: 2018/07/06        Li Zebang
 */

package pool

import (
	"time"

	"github.com/TechCatsLab/firmness/github/pool"
)

var (
	// Tag -
	Tag = "default"
	// GithubTokens -
	GithubTokens = []*pool.Token{&pool.Token{Tag: Tag, Token: "3c8bf037bd6bfe3858ad8649099b0ee82f7d6cf0"}}
	// GithubPool -
	GithubPool pool.Pool
)

// InitGithubPool -
func InitGithubPool() (err error) {
	GithubPool, err = pool.NewPool(time.Second, GithubTokens...)
	return err
}
