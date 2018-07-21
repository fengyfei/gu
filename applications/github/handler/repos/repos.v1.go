package repos

import (
	"time"

	"github.com/TechCatsLab/apix/http/server"
	mgo "gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/github/repos"
)

type (
	// infoRespV1 - for Wechat Mini Programs
	infoRespV1 struct {
		ID      string    `json:"id"`
		Avatar  string    `json:"avatar"`
		Name    string    `json:"name"`
		Image   string    `json:"image"`
		Intro   string    `json:"intro"`
		Lang    []string  `json:"lang"`
		Created time.Time `json:"created"`
		Active  bool      `json:"active"`
	}
)

// ActiveListV1 - for Wechat Mini Programs
func ActiveListV1(c *server.Context) error {
	var resp = make([]infoRespV1, 0)

	rlist, err := repos.Service.ActiveList()
	if err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}

		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	for _, r := range rlist {
		lang := make([]string, 0)
		for index := range r.Languages {
			if r.Languages[index].Proportion < 0.3 {
				break
			}
			lang = append(lang, r.Languages[index].Language)
		}
		info := infoRespV1{
			ID:      r.ID.Hex(),
			Avatar:  *r.Avatar,
			Name:    *r.Name,
			Image:   *r.Image,
			Intro:   *r.Intro,
			Lang:    lang,
			Created: r.Created,
			Active:  r.Active,
		}

		resp = append(resp, info)
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}
