/*
 * Revision History:
 *     Initial: 2017/12/28        Jia Chenhui
 *     Modify:  2018/06/29        Li Zebang
 */

package repos

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/firmness/github/base"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/applications/github/pool"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/github/repos"
)

const (
	githubURL = "https://github.com/"
)

type (
	// infoRepos -
	infoRepos struct {
		Owner     *string          `json:"owner"`
		Avatar    *string          `json:"avatar"`
		Name      *string          `json:"name"`
		Stars     *int             `json:"stars"`
		Forks     *int             `json:"forks"`
		Intro     *string          `json:"intro"`
		Readme    *string          `json:"readme"`
		Topics    []string         `json:"topics"`
		Languages []repos.Language `json:"languages"`
	}

	// createReq - The request struct that create repos information.
	createReq struct {
		URL   *string `json:"url" validate:"required,url"`
		Image *string `json:"image"`
	}

	// activateReq - The request struct that modify repos status.
	activateReq struct {
		ID     *string `json:"id" validate:"required,alphanum,len=24"`
		Active bool    `json:"active"`
	}

	// infoReq - The request struct that get a list of repos detail information.
	infoReq struct {
		ID string `json:"id"`
	}

	// readmeReq - The request struct that gets the content of the repository README.md file.
	readmeReq struct {
		RepoOwner *string `json:"repoowner" validate:"required"`
		RepoName  *string `json:"reponame" validate:"required"`
	}

	// readmeResp - The response struct that gets the content of the repository README.md file.
	readmeResp struct {
		Content *string `json:"content"`
	}

	// infoResp - The more detail of repos.
	infoResp struct {
		ID        string           `json:"id"`
		Owner     *string          `json:"owner"`
		Avatar    *string          `json:"avatar"`
		Image     *string          `json:"image"`
		Name      *string          `json:"name"`
		Stars     *int             `json:"stars"`
		Forks     *int             `json:"forks"`
		Intro     *string          `json:"intro"`
		Readme    *string          `json:"readme"`
		Topics    []string         `json:"topics"`
		Languages []repos.Language `json:"languages"`
		Active    bool             `json:"active"`
		Created   time.Time        `json:"created"`
	}
)

func getRepositoryInformation(owner, repo, tag string) (*infoRepos, error) {
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var (
		getReturn           = &base.RepositoriesGetReturn{}
		getReadmeReturn     = &base.RepositoriesGetReadmeReturn{}
		listAllTopicsReturn = &base.RepositoriesListAllTopicsReturn{}
		listLanguagesReturn = &base.RepositoriesListLanguagesReturn{}
	)
	go base.RepositoriesGet(ctx, owner, repo, pool.Tag, pool.GithubPool, wg, getReturn)
	go base.RepositoriesGetReadme(ctx, owner, repo, pool.Tag, nil, pool.GithubPool, wg, getReadmeReturn)
	go base.RepositoriesListAllTopics(ctx, owner, repo, pool.Tag, pool.GithubPool, wg, listAllTopicsReturn)
	go base.RepositoriesListLanguages(ctx, owner, repo, pool.Tag, pool.GithubPool, wg, listLanguagesReturn)

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
		ls    = make([]repos.Language, len(listLanguagesReturn.Languages))
		sum   float32
		index int
	)
	for key, val := range listLanguagesReturn.Languages {
		ls[index] = repos.Language{
			Language:   key,
			Proportion: float32(val),
		}
		sum += ls[index].Proportion
		index++
	}
	for index := range ls {
		ls[index].Proportion /= sum
	}
	for out := 0; out < len(ls)-1; out++ {
		for in := out + 1; in < len(ls); in++ {
			if ls[out].Proportion < ls[in].Proportion {
				ls[out], ls[in] = ls[in], ls[out]
			}
		}
	}

	return &infoRepos{
		Owner:     &owner,
		Avatar:    getReturn.Repository.Owner.AvatarURL,
		Name:      &repo,
		Stars:     getReturn.Repository.StargazersCount,
		Forks:     getReturn.Repository.ForksCount,
		Intro:     getReturn.Repository.Description,
		Readme:    getReadmeReturn.RepositoryContent.DownloadURL,
		Topics:    listAllTopicsReturn.Topics,
		Languages: ls,
	}, nil
}

// Create - Create repos information.
func Create(c *server.Context) error {
	var (
		err error
		req createReq
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if !strings.HasPrefix(*req.URL, githubURL) {
		logger.Error(*req.URL, "is not on github")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	slice := strings.Split(strings.TrimPrefix(*req.URL, githubURL), "/")
	if len(slice) < 2 || slice[0] == "" || slice[1] == "" {
		logger.Error(*req.URL, "is invalid")
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	_, err = repos.Service.GetByOwnerAndName(&slice[0], &slice[1])
	if err == nil {
		return core.WriteStatusAndDataJSON(c, constants.ErrDuplicate, nil)
	} else if err != mgo.ErrNotFound {
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	ri, err := getRepositoryInformation(slice[0], slice[1], pool.Tag)
	if err != nil {
		logger.Error("cann't get infomation from github", err)
		if strings.Contains(err.Error(), "404 Not Found") {
			return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
		}
		return core.WriteStatusAndDataJSON(c, constants.ErrInternalServerError, nil)
	}

	id, err := repos.Service.Create(ri.Owner, ri.Avatar, ri.Name, req.Image, ri.Intro, ri.Readme, ri.Stars, ri.Forks, ri.Topics, ri.Languages)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// ModifyActive - Modify repos status.
func ModifyActive(c *server.Context) error {
	var (
		err error
		req activateReq
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = repos.Service.ModifyActive(req.ID, req.Active); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// List - Get all the repos.
func List(c *server.Context) error {
	list, err := repos.Service.List()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrNotFound, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp := make([]infoResp, 0, len(list))
	for _, r := range list {
		info := infoResp{
			ID:        r.ID.Hex(),
			Owner:     r.Owner,
			Avatar:    r.Avatar,
			Image:     r.Image,
			Name:      r.Name,
			Stars:     r.Stars,
			Forks:     r.Forks,
			Intro:     r.Intro,
			Readme:    r.Readme,
			Topics:    r.Topics,
			Languages: r.Languages,
			Active:    r.Active,
			Created:   r.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// ActiveList - Get all the active repos.
func ActiveList(c *server.Context) error {
	list, err := repos.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrNotFound, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp := make([]infoResp, 0, len(list))
	for _, r := range list {
		info := infoResp{
			ID:        r.ID.Hex(),
			Owner:     r.Owner,
			Avatar:    r.Avatar,
			Image:     r.Image,
			Name:      r.Name,
			Stars:     r.Stars,
			Forks:     r.Forks,
			Intro:     r.Intro,
			Readme:    r.Readme,
			Topics:    r.Topics,
			Languages: r.Languages,
			Active:    r.Active,
			Created:   r.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// Info - Get ten records that are greater than the specified ID.
func Info(c *server.Context) error {
	var (
		err error
		req infoReq
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	list, err := repos.Service.GetByID(&req.ID)
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrNotFound, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp := make([]infoResp, 0, len(list))
	for _, r := range list {
		info := infoResp{
			ID:        r.ID.Hex(),
			Owner:     r.Owner,
			Avatar:    r.Avatar,
			Image:     r.Image,
			Name:      r.Name,
			Stars:     r.Stars,
			Forks:     r.Forks,
			Intro:     r.Intro,
			Readme:    r.Readme,
			Topics:    r.Topics,
			Languages: r.Languages,
			Active:    r.Active,
			Created:   r.Created,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// ReadmeURL - Gets the URL of the repository README.md file.
func ReadmeURL(c *server.Context) error {
	var (
		err  error
		req  readmeReq
		resp readmeResp
		r    *repos.Repos
	)

	if err = c.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err = c.Validate(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	r, err = repos.Service.GetByOwnerAndName(req.RepoOwner, req.RepoName)
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrNotFound, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp.Content = r.Readme
	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
