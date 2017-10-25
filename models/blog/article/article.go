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

type ArticleServiceProvider struct{}

var (
	ArticleService *ArticleServiceProvider
	mdSess         *mongo.Session
)

// PrepareArticle initializing database and create index.
func PrepareArticle() {
	url := beego.AppConfig.String("mongo::url") + "/" + common.MDBlogDName

	titleIndex := &mgo.Index{
		Key:        []string{"title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	mdSess = mongo.InitMDSess(url, common.MDBlogDName, common.MDArticleColl, titleIndex)
	ArticleService = &ArticleServiceProvider{}
}

// MDArticle represents the article information.
type MDArticle struct {
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

// MDCreateArticle use to create article.
type MDCreateArticle struct {
	Author   string
	Title    string
	Content  string
	Abstract string
	Tag      []string
}

// MDModifyArticle use to modify article information.
type MDModifyArticle struct {
	ArticleID string
	Title     string
	Content   string
	Abstract  string
	Tag       []string
	Active    bool
}

// GetList get all the articles.
func (asp *ArticleServiceProvider) GetList() ([]MDArticle, error) {
	var (
		articles []MDArticle
		err      error
	)

	err = copy.GetMany(mdSess.CollInfo, nil, &articles)

	return articles, err
}

// GetByID get article based on article id.
func (asp *ArticleServiceProvider) GetByID(id string) (MDArticle, error) {
	var (
		article MDArticle
		err     error
	)

	bsonID := bson.M{"_id": bson.ObjectIdHex(id)}

	err = copy.GetByID(mdSess.CollInfo, bsonID, &article)

	return article, err
}

// Create create article.
func (asp *ArticleServiceProvider) Create(article *MDCreateArticle) error {
	articleInfo := MDArticle{
		ArticleID: bson.NewObjectId(),
		Author:    article.Author,
		Title:     article.Title,
		Content:   article.Content,
		CreatedAt: time.Now(),
		Active:    true,
	}

	return copy.Insert(mdSess.CollInfo, &articleInfo)
}

// Modify modify article information.
func (asp *ArticleServiceProvider) Modify(update *MDModifyArticle) error {
	selector := bson.M{"_id": bson.ObjectIdHex(update.ArticleID)}
	updater := bson.M{"$set": bson.M{
		"Title":    update.Title,
		"Content":  update.Content,
		"Abstract": update.Abstract,
		"Tag":      update.Tag,
		"Active":   update.Active,
	}}

	return copy.Update(mdSess.CollInfo, selector, updater)
}
