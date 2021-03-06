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
 *     Initial: 2018/03/06        Lin Hao
 */

package activity

import (
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
		cname = "activity"
	)

	url := conf.GithubConfig.MongoURL + "/" + github.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(github.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"Title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, github.Database, cname)
	Service = &serviceProvider{}
}

// Activity represents the TechCat activity information.
type Activity struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Title   string        `bson:"Title"`
	Image   string        `bson:"Image"`
	Intro   string        `bson:"Intro"`
	Active  bool          `bson:"Active"`
	Created time.Time     `bson:"Created"`
}

// Create activity information.
func (sp *serviceProvider) Create(title, image, intro *string) (string, error) {
	activity := Activity{
		ID:      bson.NewObjectId(),
		Title:   *title,
		Image:   *image,
		Intro:   *intro,
		Active:  true,
		Created: time.Now(),
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&activity)
	if err != nil {
		return "", err
	}

	return activity.ID.Hex(), nil
}

// List get all the activity.
func (sp *serviceProvider) List() ([]Activity, error) {
	var (
		err  error
		list []Activity
	)

	conn := session.Connect()
	defer conn.Disconnect()

	sort := "-Created"
	err = conn.GetMany(nil, &list, sort)

	return list, err
}

// ActiveList get all the active activity.
func (sp *serviceProvider) ActiveList() ([]Activity, error) {
	var (
		err  error
		list []Activity
	)

	conn := session.Connect()
	defer conn.Disconnect()

	sort := "-Created"
	err = conn.GetMany(bson.M{"Active": true}, &list, sort)

	return list, err
}

// ModifyActive modify activity status.
func (sp *serviceProvider) ModifyActive(id *string, active bool) error {
	updater := bson.M{"$set": bson.M{
		"Active": active,
	}}

	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*id)}, updater)
}
