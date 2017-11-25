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
 *     Initial: 2017/11/24        Jia Chenhui
 */

package trending

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
		cname = "trending"
	)

	url := conf.Configuration.MongoURL + "/" + github.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(github.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"Title"},
		Background: true,
		Sparse:     true,
	})
	s.DB(github.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"Lang"},
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, github.Database, cname)
	Service = &serviceProvider{}
}

// Trending represents the increasing trend of the number of GitHub library stars.
type Trending struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Title    string        `bson:"Title"`
	Abstract string        `bson:"Abstract"`
	Lang     string        `bson:"Lang"` //  today in the "20060102" format + lang
	Stars    int           `bson:"Stars"`
	Today    int           `bson:"Today"` // the increments of stars today
}

// CreateList insert multiple trending records.
func (sp *serviceProvider) CreateList(docs []*Trending) error {
	conn := session.Connect()
	defer conn.Disconnect()

	for _, d := range docs {
		info := Trending{
			Title:    d.Title,
			Abstract: d.Abstract,
			Lang:     d.Lang,
			Stars:    d.Stars,
			Today:    d.Today,
		}

		err := conn.Insert(&info)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByLang get library trending based on the language.
func (sp *serviceProvider) GetByLang(lang *string) ([]Trending, error) {
	var (
		err    error
		result []Trending
	)

	conn := session.Connect()
	defer conn.Disconnect()

	date := time.Now().Format("20060102")
	query := date + *lang

	err = conn.GetMany(bson.M{"Lang": query}, &result)

	return result, err
}
