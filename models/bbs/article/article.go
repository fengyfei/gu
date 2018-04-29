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
 *     Initial: 2018/01/24        Tong Yuehong
 */

package article

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/bbs"
)

type articleServiceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	ArticleService *articleServiceProvider
	articleSession *mongo.Connection
)

type (
	// Article represents the article information.
	Article struct {
		Id         bson.ObjectId `bson:"_id,omitempty"`
		Title      string        `bson:"title"`
		Brief      string        `bson:"brief"`
		Content    string        `bson:"content"`
		AuthorID   uint32        `bson:"authorID"`
		CategoryID bson.ObjectId `bson:"categoryID"`
		TagID      bson.ObjectId `bson:"tagID"`
		VisitNum   int64         `bson:"visitNum"`
		Created    string        `bson:"created"`
		Image      string        `bson:"image"`
		Active     bool          `bson:"active"`
	}
)

func init() {
	const (
		cname = "article"
	)

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"title"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	})

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"categoryID", "tagID"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	})

	articleSession = mongo.NewConnection(initialize.S, bbs.Database, cname)
}

// Insert - add article.
func (sp *articleServiceProvider) Insert(con orm.Connection, article *Article) error {
	err := CategoryService.UpdateArtNum(article.CategoryID.Hex(), bbs.Increase)
	if err != nil {
		return err
	}

	art := &Article{
		Title:      article.Title,
		Brief:      article.Brief,
		Content:    article.Content,
		AuthorID:   article.AuthorID,
		CategoryID: article.CategoryID,
		TagID:      article.TagID,
		Created:    article.Created,
		Image:      article.Image,
		Active:     true,
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Insert(art)
}

// GetByCategoryID return articles by categoryID.
func (sp *articleServiceProvider) GetByCategoryID(page int, categoryID string) ([]Article, error) {
	var (
		list []Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"categoryID": bson.ObjectIdHex(categoryID), "active": true}
	err := conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByTagID return articles by tagID.
func (sp *articleServiceProvider) GetByTagID(page int, categoryID, tagID string) ([]Article, error) {
	var (
		list []Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"categoryID": bson.ObjectIdHex(categoryID), "tagID": bson.ObjectIdHex(tagID), "active": true}
	err := conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// SearchByTitle return articles by searching title.
func (sp *articleServiceProvider) SearchByTitle(title string) ([]Article, error) {
	var (
		list []Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-created"

	query := bson.M{"title": bson.M{"$regex": title, "$options": "$i"}, "active": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByArtID return article by artID.
func (sp *articleServiceProvider) GetByArtID(artID string) (*Article, error) {
	var (
		list Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(artID), "active": true}
	err := conn.GetUniqueOne(query, &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

// GetByUserID return articles by title.
func (sp *articleServiceProvider) GetByUserID(userID uint32) ([]Article, error) {
	var (
		list []Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "created"

	query := bson.M{"authorID": userID, "active": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Delete deletes article.
func (sp *articleServiceProvider) Delete(artID string) error {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	art, err := sp.GetByArtID(artID)
	if err != nil {
		return err
	}

	err = CategoryService.UpdateArtNum(art.CategoryID.Hex(), bbs.Decrease)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"active": false}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artID)}, updater)
}

// UpdateCommentNum update the commentNum.
func (sp *articleServiceProvider) UpdateCommentNum(artID bson.ObjectId, operation int) error {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$inc": bson.M{"commentNum": operation}}

	return conn.Update(bson.M{"_id": artID}, updater)
}

// Updatevisit update visitNum.
func (sp *articleServiceProvider) UpdateVisit(num int64, artID string) error {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"visitNum": num}}

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artID), "active": true}, updater)
}

// DeleteByCategory delete the articles that belong to the deleted category.
func (sp *articleServiceProvider) DeleteByCategory(categoryID string) error {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"active": false}}

	_, err := conn.Collection().UpdateAll(bson.M{"categoryID": bson.ObjectIdHex(categoryID)}, updater)
	return err
}

// DeleteByTag deletes articles that belong to the deleted tag.
func (sp *articleServiceProvider) DeleteByTag(categoryID, tagID string) error {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": false}}

	_, err := conn.Collection().UpdateAll(bson.M{"categoryID": bson.ObjectIdHex(categoryID), "tagID": bson.ObjectIdHex(tagID)}, updater)
	return err
}

// Recommend gets the popular articles.
func (sp *articleServiceProvider) Recommend(page int) ([]Article, error) {
	var (
		list []Article
	)

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"active": true}

	err := conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).Sort("-visitNum").All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// ArtNum gets the number of someone's articles.
func (sp *articleServiceProvider) ArtNum(userID uint32) (int, error) {
	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"authorID": userID, "active": true}

	artNum, err := conn.Collection().Find(query).Count()
	if err != nil {
		return 0, err
	}

	return artNum, err
}

// IsExist checks whether the id exists.
func (sp *articleServiceProvider) IfExist(id string) error {
	conn := articleSession.Connect()
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
