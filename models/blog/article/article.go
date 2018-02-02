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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/blog"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider
	Service *serviceProvider
	session *mongo.Connection
)

func init() {
	const (
		cname = "article"
	)

	url := beego.AppConfig.String("mongo::url") + "/" + blog.Database

	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)

	s.DB(blog.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, blog.Database, cname)
	Service = &serviceProvider{}
}

// Article represents the article information.
type Article struct {
	ArticleID bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"required"`
	Author    string        `bson:"Author" json:"author"`
	Title     string        `bson:"Title" json:"title"`
	Content   string        `bson:"Content" json:"content"`
	Abstract  string        `bson:"Abstract" json:"abstract"`
	Tag       []string      `bson:"Tag" json:"tag"`
	CreatedAt time.Time     `bson:"CreatedAt" json:"created_at"`
	UpdatedAt time.Time     `bson:"UpdatedAt" json:"updated_at"`
	Active    bool          `bson:"Active" json:"active"`
}

// List get all the articles.
func (sp *serviceProvider) List() ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(nil, &articles)

	return articles, err
}

// ActiveList get all the active articles.
func (sp *serviceProvider) ActiveList() ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(bson.M{"Active": true}, &articles)

	return articles, err
}

// GetByID get article based on article id.
func (sp *serviceProvider) GetByID(id string) (Article, error) {
	var (
		article Article
		err     error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetByID(bson.ObjectIdHex(id), &article)

	return article, err
}

// GetByTags get articles based on tag id.
func (sp *serviceProvider) GetByTags(tags *[]string) ([]Article, error) {
	var (
		articles []Article
		err      error
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(bson.M{"Tag": bson.M{"$all": *tags}}, &articles)

	return articles, err
}

// Create create article.
func (sp *serviceProvider) Create(author, title, abstract, content *string, tag *[]string) (string, error) {
	articleInfo := Article{
		ArticleID: bson.NewObjectId(),
		Author:    *author,
		Title:     *title,
		Content:   *content,
		Tag:       *tag,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&articleInfo)
	if err != nil {
		return "", err
	}

	return articleInfo.ArticleID.Hex(), nil
}

// Modify modify article information.
func (sp *serviceProvider) Modify(id, title, content, abstract *string, active *bool) error {
	updater := bson.M{"$set": bson.M{
		"Title":     *title,
		"Content":   *content,
		"Abstract":  *abstract,
		"Active":    *active,
		"UpdatedAt": time.Now(),
	}}

	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*id)}, updater)
}

// AddTags add tags to specified article.
func (sp *serviceProvider) AddTags(articleID string, tags []string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pushAll": bson.M{"Tag": tags}})
}

// RemoveTags remove tags from specified article.
func (sp *serviceProvider) RemoveTags(articleID string, tags []string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pullAll": bson.M{"Tag": tags}})
}
