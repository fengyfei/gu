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
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/bbs"
)

type articleserviceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	ArticleService *articleserviceProvider
	articleSession *mongo.Connection
)

// Article represents the article information.
type Article struct {
	Id         bson.ObjectId `bson:"_id,omitempty"         json:"id"`
	Title      string        `bson:"title"       json:"title" validate:"required,max=12"`
	UserId     uint64        `bson:"userId"      json:"userId"`
	Content    string        `bson:"content"     json:"content"`
	ModuleId   bson.ObjectId `bson:"moduleId"    json:"moduleId"`
	ThemeId    bson.ObjectId `bson:"themeId"     json:"themeId"`
	CommentNum int64         `bson:"commentNum"  json:"commentNum"`
	Times      int64         `bson:"times"       json:"times"`
	Created    time.Time     `bson:"created"     json:"created"`
	Image      string        `bson:"image"       json:"image"`
	Status     bool          `bson:"status"      json:"status"`
}

// CreateArticle represents the article information when created.
type CreateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Module  string `json:"module"`
	Theme   string `json:"theme"`
	Image   string `json:"image"`
}

func init() {
	const (
		CollArticle = "article"
	)

	url := conf.BBSConfig.MongoURL + "/" + bbs.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(bbs.Database).C(CollArticle).EnsureIndex(mgo.Index{
		Key:        []string{"title"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	articleSession = mongo.NewConnection(s, bbs.Database, CollArticle)
}

// Insert - add article.
func (sp *articleserviceProvider) Insert(article CreateArticle, userId uint64) (string, error) {
	moduleId, err := ModuleService.GetModuleID(article.Module)
	if err != nil {
		return "", err
	}

	ThemeId, err := ModuleService.GetThemeID(article.Module, article.Theme)
	if err != nil {
		return "", err
	}

	art := Article{
		Title:      article.Title,
		UserId:     userId,
		Content:    article.Content,
		ModuleId:   moduleId,
		ThemeId:    ThemeId,
		CommentNum: 0,
		Times:      0,
		Created:    time.Now(),
		Image:      article.Image,
		Status:     true,
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	err = conn.Insert(&art)
	if err != nil {
		return "", err
	}

	artId, err := sp.GetId(art.Title)
	err = ModuleService.UpdateArtNum(article.Module, "add")

	return artId.Hex(), err
}

// GetByModuleID gets articles by moduleId.
func (sp *articleserviceProvider) GetByModuleID(page int, module string) ([]Article, error) {
	var list []Article

	moduleId, err := ModuleService.GetModuleID(module)
	if err != nil {
		return list, err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"moduleId": moduleId, "status": true}
	err = conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)

	return list, err
}

// GetByThemeID get articles by themeId.
func (sp *articleserviceProvider) GetByThemeID(page int, module, theme string) ([]Article, error) {
	var list []Article

	moduleId, err := ModuleService.GetModuleID(module)
	if err != nil {
		return list, err
	}

	themeId, err := ModuleService.GetThemeID(module, theme)
	if err != nil {
		return list, err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"moduleId": moduleId, "themeId": themeId, "status": true}
	err = conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)

	return list, err
}

// GetByTitle get articles by title.
func (sp *articleserviceProvider) GetByTitle(title string) ([]Article, error) {
	var list []Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"title": bson.M{"$regex": title, "$options": "$i"}, "status": true}
	err := conn.GetMany(query, &list, sort)

	return list, err
}

// GetId gets ArtId.
func (sp *articleserviceProvider) GetId(title string) (bson.ObjectId, error) {
	var art Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"title": title}

	err := conn.GetUniqueOne(query, &art)

	return art.Id, err
}

// GetInfo gets article's information.
func (sp *articleserviceProvider) GetInfo(artId bson.ObjectId) (Article, error) {
	var article Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": artId}
	err := conn.GetUniqueOne(query, &article)

	return article, err
}

// Delete deletes article.
func (sp *articleserviceProvider) Delete(title string) error {
	artId, err := sp.GetId(title)
	if err != nil {
		return err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": false}}
	err = conn.Update(bson.M{"_id": artId}, updater)
	if err != nil {
		return nil
	}

	art, err := sp.GetInfo(artId)
	if err != nil {
		return err
	}

	module, err := ModuleService.GetInfo(art.ModuleId)
	if err != nil {
		return err
	}

	err = ModuleService.UpdateArtNum(module.Name, "sub")
	return err
}
