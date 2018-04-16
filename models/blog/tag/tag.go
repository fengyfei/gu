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
 *     Initial: 2017/10/25        Jia Chenhui
 *     Modify : 2018/02/05        Tong Yuehong
 *     Modify : 2018/03/25        Chen Yanchen
 */

package tag

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/blog/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/blog"
)

type tagServiceProvider struct{}

var (
	// TagService expose tagServiceProvider
	TagService *tagServiceProvider
	session    *mongo.Connection
)

func init() {
	const cname = "tag"

	url := conf.Config.MongoURL + "/" + blog.Database

	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)

	s.DB(blog.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"tag"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, blog.Database, cname)
	TagService = &tagServiceProvider{}
}

type (
	// Tag represents the tag information.
	Tag struct {
		TagID  bson.ObjectId `bson:"_id,omitempty"`
		Tag    string        `bson:"tag"`
		Count  int           `bson:"count"`
		Active bool          `bson:"active"`
	}
)

// GetList get all the tags.
func (sp *tagServiceProvider) GetList() ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(nil, &tags)

	return tags, err
}

// GetActiveList get all the active tags.
func (sp *tagServiceProvider) GetActiveTags() ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(bson.M{"active": true}, &tags)

	return tags, err
}

// GetByID get tag based on article id.
func (sp *tagServiceProvider) GetByID(id *string) (Tag, error) {
	var (
		tag Tag
		err error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetByID(bson.ObjectIdHex(*id), &tag)

	return tag, err
}

// Create create tag.
func (sp *tagServiceProvider) Create(tag *string) (string, error) {
	tagInfo := &Tag{
		Tag:    *tag,
		Active: true,
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(tagInfo)
	if err != nil {
		return "", err
	}
	conn.GetUniqueOne(bson.M{"tag": tag}, tagInfo)
	return tagInfo.TagID.Hex(), nil
}

// Modify modify tag information.
func (sp *tagServiceProvider) Modify(id, tag *string, active *bool) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{}
	if tag != nil {
		updater["tag"] = *tag
	}

	if active != nil {
		updater["active"] = *active
	}

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*id)}, updater)
}

// GetID return tag's id.
func (sp *tagServiceProvider) GetID(tag []string) (bson.ObjectId, error) {
	var tagInfo Tag

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"tag": tag, "active": true}

	err := conn.GetUniqueOne(query, &tagInfo)
	if err != nil {
		return "", err
	}

	return tagInfo.TagID, nil
}

// UpdateCount update Tag.Count
func (sp *tagServiceProvider) UpdateCount(id bson.ObjectId, num int) error {
	conn := session.Connect()
	defer conn.Disconnect()

	q := bson.M{"_id": id}
	err := conn.Update(q, bson.M{"$set": bson.M{"count": num}})
	return err
}
