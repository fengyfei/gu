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
	GithubTokens = []*pool.Token{&pool.Token{Tag: Tag, Token: "e32d654e6fd756a8f5a5df156f3313f6e3a6ff21"}}
	// GithubPool -
	GithubPool pool.Pool
)

// InitGithubPool -
func InitGithubPool() (err error) {
	GithubPool, err = pool.NewPool(time.Second, GithubTokens...)
	return err
}
