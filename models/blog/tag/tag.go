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
	"github.com/fengyfei/nuts/mgo/copy"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/common"
	"github.com/fengyfei/gu/pkg/mongo"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider
	Service *serviceProvider
	mdSess  *mongo.Session
)

// Prepare initializing database.
func Prepare() {
	url := beego.AppConfig.String("mongo::url") + "/" + common.MDBlogDName

	mdSess = mongo.InitMDSess(url, common.MDBlogDName, common.MDTagColl, nil)
	Service = &serviceProvider{}
}

// MDTag represents the tag information.
type MDTag struct {
	TagID  bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Tag    string        `bson:"Tag" json:"tag"`
	Active bool          `bson:"Active" json:"active"`
}

// MDCreateTag use to create article.
type MDCreateTag struct {
	Tag string
}

// MDModifyTag use to modify tag information.
type MDModifyTag struct {
	TagID  string
	Tag    string
	Active bool
}

// GetList get all the tags.
func (sp *serviceProvider) GetList() ([]MDTag, error) {
	var (
		tags []MDTag
		err  error
	)

	err = copy.GetMany(mdSess.CollInfo, nil, &tags)

	return tags, err
}

// GetActiveList get all the active tags.
func (sp *serviceProvider) GetActiveList() ([]MDTag, error) {
	var (
		tags []MDTag
		err  error
	)

	selector := bson.M{"Active": true}
	err = copy.GetMany(mdSess.CollInfo, selector, &tags)

	return tags, err
}

// GetByID get tag based on article id.
func (sp *serviceProvider) GetByID(id string) (MDTag, error) {
	var (
		tag MDTag
		err error
	)

	selector := bson.M{"_id": bson.ObjectIdHex(id)}
	err = copy.GetMany(mdSess.CollInfo, selector, &tag)

	return tag, err
}

// Create create tag.
func (sp *serviceProvider) Create(tag string) error {
	tagInfo := MDTag{
		TagID:  bson.NewObjectId(),
		Tag:    tag,
		Active: true,
	}

	return copy.Insert(mdSess.CollInfo, &tagInfo)
}

// Modify modify tag information.
func (sp *serviceProvider) Modify(update *MDModifyTag) error {
	selector := bson.M{"_id": bson.ObjectIdHex(update.TagID)}
	updater := bson.M{"$set": bson.M{
		"Tag":    update.Tag,
		"Active": update.Active,
	}}

	return copy.Update(mdSess.CollInfo, selector, updater)
}
