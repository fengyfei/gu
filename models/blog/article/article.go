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
 *     Initial: 2017/10/24        Jia Chenhui
 */

package article

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/fengyfei/nuts/mgo/copy"
	"gopkg.in/mgo.v2"
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

// Prepare initializing database and create index.
func Prepare() {
	url := beego.AppConfig.String("mongo::url") + "/" + common.MDBlogDName

	titleIndex := &mgo.Index{
		Key:        []string{"title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	mdSess = mongo.InitSession(url, common.MDBlogDName, common.ArticleColl, titleIndex)
	Service = &serviceProvider{}
}

// Article represents the article information.
type Article struct {
	ArticleID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author    string        `bson:"Author" json:"author"`
	Title     string        `bson:"Title" json:"title"`
	Content   string        `bson:"Content" json:"content"`
	Abstract  string        `bson:"Abstract" json:"abstract"`
	Tag       []string      `bson:"Tag" json:"tag"`
	CreatedAt time.Time     `bson:"CreatedAt" json:"created_at"`
	UpdatedAt time.Time     `bson:"UpdatedAt" json:"updated_at"`
	Active    bool          `bson:"Active" json:"active"`
}

// GetList get all the articles.
func (sp *serviceProvider) GetList() ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	err = copy.GetMany(mdSess.CollInfo, nil, &articles)

	return articles, err
}

// GetActiveList get all the active articles.
func (sp *serviceProvider) GetActiveList() ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	selector := bson.M{"Active": true}
	err = copy.GetMany(mdSess.CollInfo, selector, &articles)

	return articles, err
}

// GetByID get article based on article id.
func (sp *serviceProvider) GetByID(id string) (Article, error) {
	var (
		article Article
		err     error
	)

	objID := bson.ObjectIdHex(id)
	err = copy.GetByID(mdSess.CollInfo, objID, &article)

	return article, err
}

// GetByTags get articles based on tag id.
func (sp *serviceProvider) GetByTags(tags []string) ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	selector := bson.M{"Tag": bson.M{"$all": tags}}
	err = copy.GetMany(mdSess.CollInfo, selector, &articles)

	return articles, err
}

// Create create article.
func (sp *serviceProvider) Create(article *Article) (string, error) {
	articleInfo := Article{
		ArticleID: bson.NewObjectId(),
		Author:    article.Author,
		Title:     article.Title,
		Content:   article.Content,
		Abstract:  article.Abstract,
		Tag:       article.Tag,
		CreatedAt: time.Now(),
		Active:    true,
	}

	err := copy.Insert(mdSess.CollInfo, &articleInfo)
	if err != nil {
		return "", err
	}

	return articleInfo.ArticleID.Hex(), nil
}

// Modify modify article information.
func (sp *serviceProvider) Modify(update *Article) error {
	updater := bson.M{"$set": bson.M{
		"Title":     update.Title,
		"Content":   update.Content,
		"Abstract":  update.Abstract,
		"Active":    update.Active,
		"UpdatedAt": time.Now(),
	}}

	return copy.Update(mdSess.CollInfo, bson.M{"_id": bson.ObjectId(update.ArticleID)}, updater)
}

// AddTags add tags to specified article.
func (sp *serviceProvider) AddTags(articleID string, tags []string) error {
	selector := bson.M{"_id": bson.ObjectIdHex(articleID)}
	updater := bson.M{"$pushAll": bson.M{"Tag": tags}}

	return copy.Update(mdSess.CollInfo, selector, updater)
}

// RemoveTags remove tags from specified article.
func (sp *serviceProvider) RemoveTags(articleID string, tags []string) error {
	selector := bson.M{"_id": bson.ObjectIdHex(articleID)}
	updater := bson.M{"$pullAll": bson.M{"Tag": tags}}

	return copy.Update(mdSess.CollInfo, selector, updater)
}
