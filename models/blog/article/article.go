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
	"github.com/fengyfei/gu/models/blog/tag"
)

type articleServiceProvider struct{}

var (
	// ArticleService expose articleServiceProvider
	ArticleService *articleServiceProvider
	session        *mongo.Connection
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
		Key:        []string{"Title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, blog.Database, cname)
	ArticleService = &articleServiceProvider{}
}

// Article represents the article information.
type Article struct {
	ID        bson.ObjectId   `bson:"_id,omitempty" json:"id" validate:"required"`
	AuthorID  int32           `bson:"AuthorID"      json:"authorID"`
	Title     string          `bson:"Title"         json:"title"`
	Content   string          `bson:"Content"       json:"content"`
	Abstract  string          `bson:"Abstract"      json:"abstract"`
	TagsID    []bson.ObjectId `bson:"Tag"           json:"tag"`
	AuditorID int32           `bson:"auditorID"     json:"auditorID"`
	View      int32           `bson:"view"          json:"view"`
	CreatedAt time.Time       `bson:"CreatedAt"     json:"created_at"`
	UpdatedAt time.Time       `bson:"UpdatedAt"     json:"updated_at"`
	Status    int8            `bson:"status"        json:"status"`
}

// CreateArticle represents the article information when created.
type CreateArticle struct {
	AuthorID int32    `json:"authorID"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Abstract string   `json:"abstract"`
	Tag      []string `json:"tag"`
}

// Create create article.
func (sp *articleServiceProvider) Create(article CreateArticle) (string, error) {
	var tagIDs = make([]bson.ObjectId, len(article.Tag))
	for i, tags := range article.Tag {
		tagID, err := tag.TagService.GetID(tags)
		if err != nil {
			return "", err
		}

		tagIDs[i] = tagID
	}

	articleInfo := Article{
		Title:     article.Title,
		AuthorID:  article.AuthorID,
		Content:   article.Content,
		Abstract:  article.Abstract,
		TagsID:    tagIDs,
		AuditorID: 0,
		View:      0,
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
func (sp *articleServiceProvider) ListApproval(page int) ([]Article, error) {
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
func (sp *articleServiceProvider) ListCreated() ([]Article, error) {
	var articles []Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.Created}
	err := conn.GetMany(query, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ModifyStatus modify the  article status.
func (sp *articleServiceProvider) ModifyStatus(articleID string, status int8, staffID int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": status, "AuditorID": staffID}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

// Delete deletes article.
func (sp *articleServiceProvider) Delete(articleID string, staffID int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"status": blog.Delete}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

//ListDenied return articles which are denied.
func (sp *articleServiceProvider) ListDenied() ([]Article, error) {
	var articles []Article
	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.NotApproval}
	err := conn.GetMany(query, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// GetByID return the article's information.
func (sp *articleServiceProvider) GetByID(articleID string) (*Article, error) {
	var article Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(articleID), "status": blog.Approval}
	err := conn.GetUniqueOne(query, &article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// AddTags add tags to specified article.
func (sp *articleServiceProvider) AddTags(articleID string, tags []string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pushAll": bson.M{"Tag": tags}})
}

// RemoveTags remove tags from specified article.
func (sp *articleServiceProvider) RemoveTags(articleID string, tags []string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pullAll": bson.M{"Tag": tags}})
}

// ModifyArticle update article.
func (sp *articleServiceProvider) ModifyArticle(articleID string, article CreateArticle) error {
	conn := session.Connect()
	defer conn.Disconnect()

	var tagIDs = make([]bson.ObjectId, len(article.Tag))
	for i, tags := range article.Tag {
		tagID, err := tag.TagService.GetID(tags)
		if err != nil {
			return err
		}

		tagIDs[i] = tagID
	}

	updater := bson.M{"$set": bson.M{
		"title":    article.Title,
		"content":  article.Content,
		"abstract": article.Abstract,
		"tag":      tagIDs,
	}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

// UpdateView update view of article.
func (sp *articleServiceProvider) UpdateView(articleID *string, num int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"view": num}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*articleID)}, updater)
}
