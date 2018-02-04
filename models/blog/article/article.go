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
 *     Initial: 2017/10/24        Jia Chenhui
 *     Modify : 2018/02/04        Tong Yuehong
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
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id" validate:"required"`
	AuthorID  bson.ObjectId `bson:"Author"        json:"authorID"`
	Title     string        `bson:"Title"         json:"title"`
	Content   string        `bson:"Content"       json:"content"`
	Abstract  string        `bson:"Abstract"      json:"abstract"`
	Tag       []string      `bson:"Tag"           json:"tag"`
	AuditorID int32         `bson:"auditorID"     json:"auditorID"`
	CreatedAt time.Time     `bson:"CreatedAt"     json:"created_at"`
	UpdatedAt time.Time     `bson:"UpdatedAt"     json:"updated_at"`
	Status    int8          `bson:"status"        json:"status"`
}

// CreateArticle represents the article information when created.
type CreateArticle struct {
	AuthorID string   `json:"authorID"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Abstract string   `json:"abstract"`
	Tag      []string `json:"tag"`
}

// Create create article.
func (sp *serviceProvider) Create(article CreateArticle) (string, error) {
	articleInfo := Article{
		Title:     article.Title,
		AuthorID:  bson.ObjectIdHex(article.AuthorID),
		Content:   article.Content,
		Abstract:  article.Abstract,
		Tag:       article.Tag,
		AuditorID: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    blog.Created,
	}

	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&articleInfo)
	if err != nil {
		return "", err
	}

	return articleInfo.ID.Hex(), nil
}

// ListApproval returns the articles which are passed.
func (sp *serviceProvider) ListApproval(page int) ([]Article, error) {
	var articles []Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.Approval}
	err := conn.Collection().Find(query).Limit(blog.Skip).Skip(page * blog.Skip).All(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ListCreated return articles which are waiting for checking.
func (sp *serviceProvider) ListCreated() ([]Article, error) {
	var articles []Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.Created}
	err := conn.GetMany(query, &articles)
	if err != nil {
		return nil,  err
	}

	return articles, nil
}

// ModifyStatus modify the  article status.
func (sp *serviceProvider) ModifyStatus(articleID string, status int8, staffID int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"status": status, "AuditorID": staffID}
	return conn.Connect().Update(bson.M{"_id": articleID}, updater)
}

//ListDenied return articles which are denied.
func (sp *serviceProvider) ListDenied() ([]Article, error) {
	var articles []Article
	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.NotApproval}
	err := conn.Connect().GetMany(query, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// Delete delete article.
func (sp *serviceProvider) Delete(articleID string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"status": blog.Delete}
	return conn.Connect().Update(bson.M{"_id": articleID}, updater)
}

// GetByID return the article's information.
func (sp *serviceProvider) GetByID(articleID string) (*Article, error) {
	var article Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(articleID), "status": blog.Approval}
	err := conn.Connect().GetUniqueOne(query, &article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

