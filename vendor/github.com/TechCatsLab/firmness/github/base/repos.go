/*
 * Revision History:
 *     Initial: 2018/06/20        Li Zebang
 */

package base

import (
	"context"
	"sync"

	"github.com/google/go-github/github"

	"github.com/TechCatsLab/firmness/github/pool"
)

// RepositoriesGetReturn -
type RepositoriesGetReturn struct {
	Repository *github.Repository
	Response   *github.Response
	Err        error
}

// RepositoriesGet -
func RepositoriesGet(ctx context.Context, owner, repo, tag string, p pool.Pool, wg *sync.WaitGroup, ret *RepositoriesGetReturn) {
	defer wg.Done()
	client, err := p.Get(tag)
	defer p.Put(client)
	if err != nil {
		ret.Err = err
		return
	}

	repository, response, err := client.Repositories.Get(ctx, owner, repo)
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

	ret.Repository = repository
	ret.Response = response
}

// RepositoriesListAllTopicsReturn -
type RepositoriesListAllTopicsReturn struct {
	Topics   []string
	Response *github.Response
	Err      error
}

// RepositoriesListAllTopics -
func RepositoriesListAllTopics(ctx context.Context, owner, repo, tag string, p pool.Pool, wg *sync.WaitGroup, ret *RepositoriesListAllTopicsReturn) {
	defer wg.Done()
	client, err := p.Get(tag)
	defer p.Put(client)
	if err != nil {
		ret.Err = err
		return
	}

	topics, response, err := client.Repositories.ListAllTopics(ctx, owner, repo)
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

	ret.Topics = topics
	ret.Response = response
}

// RepositoriesListLanguagesReturn -
type RepositoriesListLanguagesReturn struct {
	Languages map[string]int
	Response  *github.Response
	Err       error
}

// RepositoriesListLanguages -
func RepositoriesListLanguages(ctx context.Context, owner, repo, tag string, p pool.Pool, wg *sync.WaitGroup, ret *RepositoriesListLanguagesReturn) {
	defer wg.Done()
	client, err := p.Get(tag)
	defer p.Put(client)
	if err != nil {
		ret.Err = err
		return
	}

	languages, response, err := client.Repositories.ListLanguages(ctx, owner, repo)
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

	ret.Languages = languages
	ret.Response = response
}
