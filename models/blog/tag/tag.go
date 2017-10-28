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
 */

package tag

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/blog"
	"github.com/fengyfei/nuts/mgo/copy"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider
	Service *serviceProvider
	session *mongo.Session
)

// Prepare initializing database.
func Prepare() {
	url := beego.AppConfig.String("mongo::url") + "/" + blog.Database

	session = mongo.InitSession(url, blog.Database, blog.TagIndex, nil)
	Service = &serviceProvider{}
}

// Tag represents the tag information.
type Tag struct {
	TagID  bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Tag    string        `bson:"Tag" json:"tag"`
	Active bool          `bson:"Active" json:"active"`
}

// GetList get all the tags.
func (sp *serviceProvider) GetList() ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	err = copy.GetMany(session.CollInfo, nil, &tags)

	return tags, err
}

// GetActiveList get all the active tags.
func (sp *serviceProvider) GetActiveList() ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	selector := bson.M{"Active": true}
	err = copy.GetMany(session.CollInfo, selector, &tags)

	return tags, err
}

// GetByID get tag based on article id.
func (sp *serviceProvider) GetByID(id string) (Tag, error) {
	var (
		tag Tag
		err error
	)

	objID := bson.ObjectIdHex(id)
	err = copy.GetByID(session.CollInfo, objID, &tag)

	return tag, err
}

// Create create tag.
func (sp *serviceProvider) Create(tag string) (string, error) {
	tagInfo := Tag{
		TagID:  bson.NewObjectId(),
		Tag:    tag,
		Active: true,
	}

	err := copy.Insert(session.CollInfo, &tagInfo)
	if err != nil {
		return "", err
	}

	return tagInfo.TagID.Hex(), nil
}

// Modify modify tag information.
func (sp *serviceProvider) Modify(update *Tag) error {
	selector := bson.M{"_id": bson.ObjectId(update.TagID)}
	updater := bson.M{"$set": bson.M{
		"Tag":    update.Tag,
		"Active": update.Active,
	}}

	return copy.Update(session.CollInfo, selector, updater)
}
