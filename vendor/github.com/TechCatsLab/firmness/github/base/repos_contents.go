/*
 * Revision History:
 *     Initial: 2018/07/05        Li Zebang
 */

package base

import (
	"context"
	"sync"

	"github.com/TechCatsLab/firmness/github/pool"

	"github.com/google/go-github/github"
)

// RepositoriesGetReadmeReturn -
type RepositoriesGetReadmeReturn struct {
	RepositoryContent *github.RepositoryContent
	Response          *github.Response
	Err               error
}

// RepositoriesGetReadme -
func RepositoriesGetReadme(ctx context.Context, owner, repo, tag string, opt *github.RepositoryContentGetOptions, p pool.Pool, wg *sync.WaitGroup, ret *RepositoriesGetReadmeReturn) {
	defer wg.Done()
	client, err := p.Get(tag)
	defer p.Put(client)
	if err != nil {
		ret.Err = err
		return
	}

	repositoryContent, response, err := client.Repositories.GetReadme(ctx, owner, repo, opt)
	if err != nil {
		ret.Err = err
		return
	}

	if response != nil {
		err = client.HandleResponse(response)
		if err != nil {
			ret.Err = err
			return
		}
	}

	ret.RepositoryContent = repositoryContent
	ret.Response = response
}
