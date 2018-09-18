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
	// account  inxchina@hotmail.com
	// password db8a7d70f823543be01411588ad6e2e5ad5b62df
	GithubTokens = []*pool.Token{&pool.Token{Tag: Tag, Token: "db8a7d70f823543be01411588ad6e2e5ad5b62df"}}
	// GithubPool -
	GithubPool pool.Pool
)

// InitGithubPool -
func InitGithubPool() (err error) {
	GithubPool, err = pool.NewPool(time.Second, GithubTokens...)
	return err
}
