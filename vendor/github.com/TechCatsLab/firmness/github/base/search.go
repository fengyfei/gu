package base

import (
	"context"
	"sync"

	"github.com/TechCatsLab/firmness/github/pool"
	"github.com/google/go-github/github"
)

// RepositoriesReturn -
type RepositoriesReturn struct {
	RepositoriesSearchResult *github.RepositoriesSearchResult
	Response                 *github.Response
	Err                      error
}

// SearchRepositories -
func SearchRepositories(ctx context.Context, query, tag string, opt *github.SearchOptions, p pool.Pool, wg *sync.WaitGroup, ret *RepositoriesReturn) {
	defer wg.Done()
	client, err := p.Get(tag)
	defer p.Put(client)
	if err != nil {
		ret.Err = err
		return
	}

	repositoriesSearchResult, response, err := client.Search.Repositories(ctx, query, opt)
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

	ret.RepositoriesSearchResult = repositoriesSearchResult
	ret.Response = response
}
