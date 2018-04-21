/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/01/25        Tong Yuehong
 */

package article

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/bbs"
)

type categoryServiceProvider struct{}

var (
	// ErrNotFound - No result found
	ErrNotFound = errors.New("No result found")
	ErrExist = errors.New("Already exist")
	// CategoryService expose serviceProvider
	CategoryService *categoryServiceProvider
	categorySession *mongo.Connection
)

type (
	// Tag represents the second category.
	Tag struct {
		Id     bson.ObjectId `bson:"id"    json:"id"`
		Name   string        `bson:"name"            json:"name"  validate:"required"`
		Active bool          `bson:"active"          json:"active"`
	}

	// Category represents the category information.
	Category struct {
		Id        bson.ObjectId `bson:"_id,omitempty"   json:"id"`
		Name      string        `bson:"name"            json:"name"`
		ArtNum    int64         `bson:"artNum"          json:"artNum"`
		VisitNum  int64         `bson:"visitNum"        json:"visitnum"`
		Tags      []Tag         `bson:"tags"            json:"tags"`
		Recommend int32         `bson:"recommend"       json:"recommend"`
		Active    bool          `bson:"active"          json:"active"`
	}
)

func init() {
	const (
		cname = "category"
	)

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"tags.id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	categorySession = mongo.NewConnection(initialize.S, bbs.Database, cname)
}

//CreateCategory add category.
func (sp *categoryServiceProvider) CreateCategory(category *Category) error {
	cate := &Category{
		Name:   category.Name,
		Active: true,
	}

	conn := categorySession.Connect()
	defer conn.Disconnect()

	num, err := conn.Collection().Find(bson.M{"name": cate.Name}).Count()
	if err != nil {
		return err
	}

	if num != 0 {
		return ErrExist
	}

	return conn.Insert(cate)
}

// CreateTag add tag.
func (sp *categoryServiceProvider) CreateTag(categoryID, tagName string) error {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	tag := Tag{
		Id:     bson.NewObjectId(),
		Name:   tagName,
		Active: true,
	}

	num, err := conn.Collection().Find(bson.M{"name": tag.Name}).Count()
	if err != nil {
		return err
	}

	if num != 0 {
		return ErrExist
	}

	updater := bson.M{"$addToSet": bson.M{"tags": tag}}

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(categoryID)}, updater)
}

// UpdateArtNum update the artNum of the category.
func (sp *categoryServiceProvider) UpdateArtNum(categoryID string, operation int) error {
	updater := bson.M{"$inc": bson.M{"artNum": operation}}

	conn := categorySession.Connect()
	defer conn.Disconnect()

	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(categoryID)}, updater)
	return err
}

// UpdateCategoryView update CategoryView.
func (sp *categoryServiceProvider) UpdateCategoryVisit(num int64, categoryID string) error {
	updater := bson.M{"$set": bson.M{"visitNum": num}}

	conn := categorySession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(categoryID), "active": true}, updater)
}

// ListTags return category's tags.
func (sp *categoryServiceProvider) ListInfo(categoryID string) (*Category, error) {
	var (
		category Category
	)

	conn := categorySession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(categoryID), "active": true}
	err := conn.GetUniqueOne(query, &category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// AllCategory return all categories.
func (sp *categoryServiceProvider) AllCategories() ([]Category, error) {
	var (
		category []Category
	)

	conn := categorySession.Connect()
	defer conn.Disconnect()

	sort := "-created"

	query := bson.M{"active": true}
	err := conn.GetMany(query, &category, sort)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory delete category.
func (sp *categoryServiceProvider) DeleteCategory(categoryID string) error {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"active": false}}
	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(categoryID)}, updater)
	if err != nil {
		return err
	}

	return ArticleService.DeleteByCategory(categoryID)
}

// DeleteTag delete tag.
func (sp *categoryServiceProvider) DeleteTag(categoryID, tagID string) error {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"tags.$.active": false}}
	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(categoryID), "tags.id": bson.ObjectIdHex(tagID)}, updater)
	if err != nil {
		return err
	}

	return ArticleService.DeleteByTag(categoryID, tagID)
}

// UpdateRecommend update the recommend
func UpdateRecommend() {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	categories, err := CategoryService.AllCategories()
	if err != nil {
		logger.Error(err)
	}

	for _, category := range categories {
		if category.ArtNum == 0 {
			updater := bson.M{"$set": bson.M{"recommend": 0}}
			err = conn.Update(bson.M{"_id": category.Id}, updater)
			if err != nil {
				logger.Error(err)
			}
			continue
		}

		recommend := category.VisitNum / category.ArtNum
		updater := bson.M{"$set": bson.M{"recommend": recommend}}
		err = conn.Update(bson.M{"_id": category.Id}, updater)
		if err != nil {
			logger.Error(err)
		}
	}
}

// ListRecommend return categories which are recommended.
func (sp *categoryServiceProvider) ListRecommend() ([]Category, error) {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	var list []Category
	query := bson.M{"active": true}
	err := conn.Collection().Find(query).Sort("-recommend").Limit(5).All(&list)
	return list, err
}

// IsExist checks whether the id exists.
func (sp *categoryServiceProvider) IsExist(id string) error {
	conn := categorySession.Connect()
	defer conn.Disconnect()

	num, err := conn.Collection().Find(bson.M{"_id": bson.ObjectIdHex(id)}).Count()
	if err != nil {
		return err
	}

	if num == 0 {
		return ErrNotFound
	}

	return nil
}

func (sp *categoryServiceProvider) GetCategoryByID(id string) (*Category, error) {
	var (
		category Category
	)

	conn := categorySession.Connect()
	defer conn.Disconnect()

	err := conn.GetUniqueOne(bson.M{"_id": bson.ObjectIdHex(id), "active": true}, &category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (sp *categoryServiceProvider) GetTagByID(cid string, tid string) (*Tag, error) {
	var (
		category Category
	)

	conn := categorySession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(cid), "tags.id": bson.ObjectIdHex(tid), "active": true, "tags.active": true}

	err := conn.Collection().Find(query).Select(bson.M{"tags.$": 1}).One(&category)
	if err != nil {
		return nil, err
	}

	return &category.Tags[0], nil
}
