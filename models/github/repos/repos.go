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
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/echo/github/conf"
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

	url := conf.Configuration.MongoURL + "/" + github.Database
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

// Repos represents the GitHub repository information.
type Repos struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Avatar  string        `bson:"Avatar"`
	Name    string        `bson:"Name"`
	Link    string        `bson:"Link"`
	Image   string        `bson:"Image"`
	Intro   string        `bson:"Intro"`
	Active  bool          `bson:"Active"`
	Created time.Time     `bson:"Created"`
}

// List get all the repos.
func (sp *serviceProvider) List() ([]Repos, error) {
	var (
		err  error
		list []Repos
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(nil, &list)

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

	err = conn.GetMany(bson.M{"Active": true}, &list)

	return list, err
}

// GetByID get repos based on repos id.
func (sp *serviceProvider) GetByID(id *string) (Repos, error) {
	var (
		err   error
		repos Repos
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetByID(bson.ObjectIdHex(*id), &repos)

	return repos, err
}

// Create create repos information.
func (sp *serviceProvider) Create(avatar, name, link, image, intro *string) (string, error) {
	repos := Repos{
		ID:      bson.NewObjectId(),
		Avatar:  *avatar,
		Name:    *name,
		Link:    *link,
		Image:   *image,
		Intro:   *intro,
		Active:  true,
		Created: time.Now(),
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&repos)
	if err != nil {
		return "", err
	}

	return repos.ID.Hex(), nil
}

// ModifyActive modify repos status.
func (sp *serviceProvider) ModifyActive(id *string, active bool) error {
	updater := bson.M{"$set": bson.M{
		"Active": active,
	}}

	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*id)}, updater)
}
