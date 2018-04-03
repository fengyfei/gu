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
		Key:        []string{"Title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, blog.Database, cname)
}

type (
	// Article represents the article information.
	Article struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		AuthorID  int32         `bson:"AuthorID"`
		AuditorID int32         `bson:"AuditorID"`
		Title     string        `bson:"Title"`
		Abstract  string        `bson:"Abstract"`
		Content   string        `bson:"Content"`
		Image     string        `bson:"Image"`
		Tags      []string      `bson:"Tags"`
		View      uint32        `bson:"view"`
		CreatedAt time.Time     `bson:"CreatedAt"`
		UpdatedAt time.Time     `bson:"UpdatedAt"`
		Status    int8          `bson:"status"`
	}
	// Art is Article response struct.
	Art struct {
		ID       bson.ObjectId `bson:"_id,omitempty"`
		Title    string        `bson:"Title"`
		Abstract string        `bson:"Abstract"`
		Image    string        `bson:"Image"`
	}
)

// Create create article.
func (sp *articleServiceProvider) Create(art *Article) (string, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	art.CreatedAt = time.Now()
	art.UpdatedAt = art.CreatedAt

	err := conn.Insert(art)
	if err != nil {
		return "", err
	}
	query := bson.M{"Title": art.Title}
	conn.GetUniqueOne(query, &art)

	return art.ID.Hex(), nil
}

// ListApproval returns the articles which are passed.
func (sp *articleServiceProvider) ListApproval(page int) ([]Art, error) {
	var art []Art

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
func (sp *articleServiceProvider) ListCreated() ([]Art, error) {
	var art []Art

	conn := session.Connect()
	defer conn.Disconnect()

	query := bson.M{"status": blog.Created}
	err := conn.GetMany(query, &art)
	if err != nil {
		return nil, err
	}

	return art, nil
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
func (sp *articleServiceProvider) GetByID(articleID string) (*Article, error) {
	var (
		article Article
	)

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
func (sp *articleServiceProvider) ModifyArticle(articleID string, article Article) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{
		"Title":    article.Title,
		"Content":  article.Content,
		"Abstract": article.Abstract,
		"Tag":      article.Tags,
	}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(articleID)}, updater)
}

// UpdateView update view of article.
func (sp *articleServiceProvider) UpdateView(articleID *string, num uint32) error {
	conn := session.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"view": num}}
	return conn.Update(bson.M{"_id": bson.ObjectIdHex(*articleID)}, updater)
}

// GetByTag get article by tag.
func (s *articleServiceProvider) GetByTag(tag string) ([]Art, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var art []Art
	q := bson.M{"Tags": tag, "status": blog.Approval}
	err := conn.GetMany(q, &art)
	if err != nil {
		return nil, err
	}
	return art, nil
}

// GetByAuthorID get articles by author ID.
func (s *articleServiceProvider) GetByAuthorID(id int32) ([]Art, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var art []Art
	q := bson.M{"AuthorID": id, "status": blog.Approval}
	err := conn.GetMany(q, &art)
	if err != nil {
		return nil, err
	}
	return art, nil
}
