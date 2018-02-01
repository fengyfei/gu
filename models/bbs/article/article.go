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
	"github.com/fengyfei/gu/models/user"
	"github.com/fengyfei/gu/applications/bbs/initialize"
)

type articleServiceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	ArticleService *articleServiceProvider
	articleSession *mongo.Connection
)

// Article represents the article information.
type Article struct {
	Id          bson.ObjectId `bson:"_id,omitempty"  json:"id"`
	Title       string        `bson:"title"          json:"title"`
	UserId      uint64        `bson:"userId"         json:"userId"`
	Content     string        `bson:"content"        json:"content"`
	Module      string        `bson:"module"         json:"module"`
	Theme       string        `bson:"theme"          json:"theme"`
	ModuleId    bson.ObjectId `bson:"moduleId"       json:"moduleId"`
	ThemeId     bson.ObjectId `bson:"themeId"        json:"themeId"`
	CommentNum  int64         `bson:"commentNum"     json:"commentNum"`
	Times       int64         `bson:"times"          json:"times"`
	LastComment string        `bson:"lastComment"    json:"lastComment"`
	Created     time.Time     `bson:"created"        json:"created"`
	Image       string        `bson:"image"          json:"image"`
	Status      bool          `bson:"status"         json:"status"`
}

// CreateArticle represents the article information when created.
type CreateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Module  string `json:"module"`
	Theme   string `json:"theme"`
	Image   string `json:"image"`
}

// UserReply represents the information about someone's reply.
type UserReply struct {
	Title   string    `json:"title" validate:"required,min=8,max=32"`
	Creator string    `json:"creator"`
	Replier string    `json:"replier"`
	Module  string    `json:"module"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
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
		Key:        []string{"title", "userId", "moduleId"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	articleSession = mongo.NewConnection(s, bbs.Database, CollArticle)
}

// Insert - add article.
func (sp *articleServiceProvider) Insert(article CreateArticle, userId uint64) (string, error) {
	moduleId, err := ModuleService.GetModuleID(article.Module)
	if err != nil {
		return "", err
	}

	ThemeId, err := ModuleService.GetThemeID(article.Module, article.Theme)
	if err != nil {
		return "", err
	}

	c , _ := initialize.Pool.Get()

	userInfo, err := user.UserServer.GetUserByID(c, userId)
	if err != nil {
		return "", err
	}

	art := Article{
		Title:       article.Title,
		UserId:      userId,
		Content:     article.Content,
		Module:      article.Module,
		Theme:       article.Theme,
		ModuleId:    moduleId,
		ThemeId:     ThemeId,
		CommentNum:  0,
		Times:       0,
		LastComment: userInfo.Username,
		Created:     time.Now(),
		Image:       article.Image,
		Status:      true,
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	err = conn.Insert(&art)
	if err != nil {
		return "", err
	}

	artId, err := sp.GetId(art.Title)
	err = ModuleService.UpdateArtNum(article.Module, bbs.IncCount)
	if err != nil {
		return "", err
	}

	return artId.Hex(), err
}

// GetByModuleID gets articles by moduleId.
func (sp *articleServiceProvider) GetByModuleID(page int, module string) ([]Article, error) {
	var list []Article

	moduleId, err := ModuleService.GetModuleID(module)
	if err != nil {
		return list, err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"moduleId": moduleId, "status": true}
	err = conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// GetByThemeID get articles by themeId.
func (sp *articleServiceProvider) GetByThemeID(page int, module, theme string) ([]Article, error) {
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
	if err != nil {
		return nil, err
	}

	return list, err
}

// GetByTitle get articles by title.
func (sp *articleServiceProvider) GetByTitle(title string) ([]Article, error) {
	var list []Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"title": bson.M{"$regex": title, "$options": "$i"}, "status": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, err
}

// GetByArtId get article by artId.
func (sp *articleServiceProvider) GetByArtId(artId bson.ObjectId) (Article, error) {
	var list Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": artId, "status": true}
	err := conn.GetUniqueOne(query, &list)
	if err != nil {
		return Article{}, err
	}

	return list, err
}

// GetByUserId get articles by title.
func (sp *articleServiceProvider) GetByUserId(userId uint64) ([]Article, error) {
	var list []Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"userId": userId, "status": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, err
}

// GetId gets ArtId.
func (sp *articleServiceProvider) GetId(title string) (bson.ObjectId, error) {
	var art Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"title": title}

	err := conn.GetUniqueOne(query, &art)
	if err != nil {
		return "", err
	}

	return art.Id, err
}

// GetInfo gets article's information.
func (sp *articleServiceProvider) GetInfo(artId bson.ObjectId) (Article, error) {
	var article Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": artId}
	err := conn.GetUniqueOne(query, &article)
	if err != nil {
		return Article{}, err
	}

	return article, err
}

// Delete deletes article.
func (sp *articleServiceProvider) Delete(title string) error {
	artId, err := sp.GetId(title)
	if err != nil {
		return err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": false}}
	err = conn.Update(bson.M{"_id": artId}, updater)
	if err != nil {
		return err
	}

	art, err := sp.GetInfo(artId)
	if err != nil {
		return err
	}

	module, err := ModuleService.ListInfo(art.ModuleId.Hex())
	if err != nil {
		return err
	}

	return ModuleService.UpdateArtNum(module.Name, bbs.DecCount)
}

// UpdateCommentNum update the commentNum.
func (sp *articleServiceProvider) UpdateCommentNum(artId bson.ObjectId, sort int) error {
	var updater = bson.M{}

	updater = bson.M{"$inc": bson.M{"commentNum": sort}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": artId}, updater)
}

//  UpdateTimes update times.
func (sp *articleServiceProvider) UpdateTimes(num int64, artId string) error {
	updater := bson.M{"$set": bson.M{"times": num}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artId), "status": true}, updater)
}

// DeleteByModule deletes articles by deleting module.
func (sp *articleServiceProvider) DeleteByModule(moduleId string) error {
	updater := bson.M{"$set": bson.M{"status": false}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	_, err := conn.Collection().UpdateAll(bson.M{"moduleId": bson.ObjectIdHex(moduleId)}, updater)
	return err
}

// DeleteByTheme deletes articles by deleting themes.
func (sp *articleServiceProvider) DeleteByTheme(moduleId, themeId string) error {
	updater := bson.M{"$set": bson.M{"status": false}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	_, err := conn.Collection().UpdateAll(bson.M{"moduleId": bson.ObjectIdHex(moduleId), "themeId": bson.ObjectIdHex(themeId)}, updater)
	return err
}

// UpdateLastComment update lastComment.
func (sp *articleServiceProvider) UpdateLastComment(artId, user string) error {
	updater := bson.M{"$set": bson.M{"lastComment": user}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artId), "status": true}, updater)
}
