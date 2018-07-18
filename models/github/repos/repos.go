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
 *     Initial: 2017/11/17        Jia Chenhui
 */

package repos

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/github/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/github"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider
	Service *serviceProvider
	session *mongo.Connection
)

func init() {
	const (
		cname = "repos"
	)

	url := conf.GithubConfig.MongoURL + "/" + github.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(github.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, github.Database, cname)
	Service = &serviceProvider{}
}

const (
	listSize = 10
)

type (
	// Repos - represents the GitHub repository information.
	Repos struct {
		ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Owner     *string       `bson:"Owner" json:"owner"`
		Avatar    *string       `bson:"Avatar" json:"avatar"`
		Name      *string       `bson:"Name" json:"name"`
		Image     *string       `bson:"Image" json:"image"`
		Intro     *string       `bson:"Intro" json:"intro"`
		Readme    *string       `bson:"Readme" json:"readme"`
		Stars     *int          `bson:"Stars" json:"stars"`
		Forks     *int          `bson:"Forks" json:"forks"`
		Topics    []string      `bson:"Topics" json:"topics"`
		Languages []Language    `bson:"Languages" json:"languages"`
		Active    bool          `bson:"Active" json:"active"`
		Created   time.Time     `bson:"Created" json:"created"`
	}

	// Language - represents the GitHub repository program language information.
	Language struct {
		Language   string  `bson:"Language" json:"language"`
		Proportion float32 `bson:"Proportion" json:"proportion"`
	}
)

// Create create repos information.
func (sp *serviceProvider) Create(owner, avatar, name, image, intro, readme *string, stars, forks *int, topics []string, languages []Language) (string, error) {
	repos := Repos{
		ID:        bson.NewObjectId(),
		Owner:     owner,
		Avatar:    avatar,
		Name:      name,
		Image:     image,
		Intro:     intro,
		Readme:    readme,
		Stars:     stars,
		Forks:     forks,
		Topics:    topics,
		Languages: languages,
		Active:    true,
		Created:   time.Now(),
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&repos)
	if err != nil {
		return "", err
	}

	return repos.ID.Hex(), nil
}

// List get all the repos.
func (sp *serviceProvider) List() ([]Repos, error) {
	var (
		err  error
		list []Repos
	)

	conn := session.Connect()
	defer conn.Disconnect()

	sort := "-Created"
	err = conn.GetMany(nil, &list, sort)

	return list, err
}

// ActiveList get all the active repos.
func (sp *serviceProvider) ActiveList() ([]Repos, error) {
	var (
		err  error
		list []Repos
	)

	conn := session.Connect()
	defer conn.Disconnect()

	sort := "-Created"
	err = conn.GetMany(bson.M{"Active": true}, &list, sort)

	return list, err
}

// GetByID get a list of records that are greater than the specified ID.
func (sp *serviceProvider) GetByID(id *string) ([]Repos, error) {
	var (
		err   error
		list  []Repos
		query bson.M
		sort  = "-Created"
	)

	conn := session.Connect()
	defer conn.Disconnect()

	if id == nil || *id == "" {
		query = nil
	} else {
		query = bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex(*id)}}
	}

	err = conn.GetLimitedRecords(query, listSize, &list, sort)

	return list, err
}

// GetByName get a record by name.
func (sp *serviceProvider) GetByName(name *string) (*Repos, error) {
	var (
		err error
		doc Repos
	)

	conn := session.Connect()
	defer conn.Disconnect()

	if name == nil || *name == "" {
		return nil, errors.New("name cann't be empty")
	}

	err = conn.GetUniqueOne(bson.M{"Name": *name}, &doc)

	return &doc, err
}

// ModifyActive modify repos status.
func (sp *serviceProvider) ModifyActive(id *string, active bool) error {
	updater := bson.M{"$set": bson.M{"Active": active}}

	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*id)}, updater)
}
