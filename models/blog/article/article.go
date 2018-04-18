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
 *     Modify : 2018/03/25        Chen Yanchen
 */

package article

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/blog/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/blog"
	"regexp"
)

type articleServiceProvider struct{}

var (
	// ArticleService expose articleServiceProvider
	ArticleService *articleServiceProvider
	session        *mongo.Connection
)

func init() {
	const cname = "article"

	url := conf.Config.MongoURL + "/" + blog.Database

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
}

// Article represents the article information.
type Article struct {
	ID        bson.ObjectId   `bson:"_id,omitempty"`
	AuthorId  int32           `bson:"authorid"`
	AuditorId int32           `bson:"auditorid"`
	Title     string          `bson:"title"`
	Brief     string          `bson:"brief"`
	Content   string          `bson:"content"`
	Image     string          `bson:"image"`
	TagsID    []bson.ObjectId `bson:"tagsid"`
	Views     uint64          `bson:"views"`
	Created   time.Time       `bson:"created"`
	Updated   time.Time       `bson:"updated"`
	Status    int8            `bson:"status"`
}

// Create create article.
func (sp *articleServiceProvider) Create(art *Article) (string, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	art.Created = time.Now()
	art.Updated = art.Created

	err := conn.Insert(&art)
	if err != nil {
		return "", err
	}
	query := bson.M{"title": art.Title}
	conn.GetUniqueOne(query, &art)

	return art.ID.Hex(), nil
}

// ListApproval returns the articles which are passed.
func (sp *articleServiceProvider) ListApproval(page int) ([]Article, error) {
	var art []Article

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.Approval}
	err := conn.Collection().Find(query).Limit(blog.Skip).Skip(page * blog.Skip).All(&art)
	if err != nil {
		return nil, err
	}

	return art, nil
}

// ListCreated return articles which are waiting for checking.
func (sp *articleServiceProvider) ListCreated(page int) ([]Article, error) {
	var art []Article

	conn := session.Connect()
	defer conn.Disconnect()

	q := bson.M{"status": blog.Created}
	err := conn.Collection().Find(q).Limit(blog.Skip).Skip(page * blog.Skip).All(&art)
	if err != nil {
		return nil, err
	}

	return art, nil
}

// ModifyStatus modify the  article status.
func (sp *articleServiceProvider) ModifyStatus(articleID string, status int8, staffID int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": status, "auditorid": staffID, "updated": time.Now()}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

// Delete deletes article.
func (sp *articleServiceProvider) Delete(articleID string, staffID int32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"auditorid": staffID, "status": blog.Delete}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

//ListDenied return articles which are denied.
func (sp *articleServiceProvider) ListDenied() ([]Article, error) {
	var (
		articles []Article
	)

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
func (sp *articleServiceProvider) GetByID(articleID string) (Article, error) {
	var (
		article Article
	)

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(articleID), "status": blog.Approval}
	err := conn.GetUniqueOne(query, &article)
	if err != nil {
		return article, err
	}

	return article, nil
}

// AddTags add tags to specified article.
func (sp *articleServiceProvider) AddTags(articleID string, tags []string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pushAll": bson.M{"tagsid": tags}})
}

// RemoveTags remove tags from specified article.
func (sp *articleServiceProvider) RemoveTags(articleID string, tagsid []bson.ObjectId) error {
	conn := session.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, bson.M{"$pullAll": bson.M{"tagsid": tagsid}})
}

// ModifyArticle update article.
func (sp *articleServiceProvider) ModifyArticle(articleID string, article *Article) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{
		"title":   article.Title,
		"content": article.Content,
		"tagsid":  article.TagsID,
		"update":  time.Now(),
	}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

// UpdateView update view of article.
func (sp *articleServiceProvider) UpdateView(articleID *string, num uint64) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"view": num}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*articleID)}, updater)
}

// GetByTag get article by tag.
func (s *articleServiceProvider) GetByTagId(id string, page int) ([]Article, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var art []Article
	q := bson.M{"tagsid": bson.ObjectIdHex(id), "status": blog.Approval}
	err := conn.Collection().Find(q).Limit(blog.Skip).Skip(page * blog.Skip).All(&art)
	if err != nil {
		return nil, err
	}
	return art, nil
}

// GetByAuthorID get articles by author ID.
func (s *articleServiceProvider) GetByAuthorID(id int32) ([]Article, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var art []Article
	q := bson.M{"authorid": id, "status": blog.Approval}
	err := conn.GetMany(q, &art)
	if err != nil {
		return nil, err
	}
	return art, nil
}

// CountByTag
func (sp *articleServiceProvider) CountByTag(id bson.ObjectId) (int, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	q := bson.M{"tagsid": id}
	num, err := conn.Collection().Find(q).Count()
	if err != nil {
		return 0, err
	}
	return num, nil
}

// GetBrief get the first line from content.
func (sp *articleServiceProvider) GetBrief(content string) string {
	reg, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	brief := reg.ReplaceAllString(content, "")
	return brief
}
