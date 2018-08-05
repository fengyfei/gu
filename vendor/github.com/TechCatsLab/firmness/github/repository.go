/*
 * Revision History:
 *     Initial: 2018/06/20        Li Zebang
 */

package github

import (
	"context"
	"sync"

	"github.com/TechCatsLab/firmness/github/base"
	"github.com/TechCatsLab/firmness/github/pool"
)

// RepositoryInformation -
type RepositoryInformation struct {
	Owner       *string    `json:"owner"`
	Avatar      *string    `json:"avatar"`
	Repo        *string    `json:"repo"`
	Description *string    `json:"description"`
	Topics      []string   `json:"topics"`
	Languages   []language `json:"languages"`
	Readme      *string    `json:"readme"`
}

type language struct {
	Language   string  `json:"language"`
	Proportion float32 `json:"proportion"`
}

// GetRepositoryInformation -
func GetRepositoryInformation(owner, repo, tag string, p pool.Pool) (*RepositoryInformation, error) {
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var (
		getReturn           = &base.RepositoriesGetReturn{}
		getReadmeReturn     = &base.RepositoriesGetReadmeReturn{}
		listAllTopicsReturn = &base.RepositoriesListAllTopicsReturn{}
		listLanguagesReturn = &base.RepositoriesListLanguagesReturn{}
	)
	go base.RepositoriesGet(ctx, owner, repo, tag, p, wg, getReturn)
	go base.RepositoriesGetReadme(ctx, owner, repo, tag, nil, p, wg, getReadmeReturn)
	go base.RepositoriesListAllTopics(ctx, owner, repo, tag, p, wg, listAllTopicsReturn)
	go base.RepositoriesListLanguages(ctx, owner, repo, tag, p, wg, listLanguagesReturn)

	wg.Wait()

	if getReturn.Err != nil {
		return nil, getReturn.Err
	}
	if getReadmeReturn.Err != nil {
		return nil, getReadmeReturn.Err
	}
	if listAllTopicsReturn.Err != nil {
		return nil, listAllTopicsReturn.Err
	}
	if listLanguagesReturn.Err != nil {
		return nil, listLanguagesReturn.Err
	}

	var (
		ls    = make([]language, len(listLanguagesReturn.Languages))
		sum   float32
		index int
	)
	for key, val := range listLanguagesReturn.Languages {
		ls[index] = language{
			Language:   key,
			Proportion: float32(val),
		}
		sum += ls[index].Proportion
		index++
	}
	for index := range ls {
		ls[index].Proportion /= sum
	}

	return &RepositoryInformation{
		Owner:       &owner,
		Avatar:      getReturn.Repository.Owner.AvatarURL,
		Repo:        &repo,
		Description: getReturn.Repository.Description,
		Topics:      listAllTopicsReturn.Topics,
		Languages:   ls,
		Readme:      getReadmeReturn.RepositoryContent.DownloadURL,
	}, nil
}
