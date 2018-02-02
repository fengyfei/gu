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
	UserID      uint64        `bson:"userID"         json:"userID"`
	Content     string        `bson:"content"        json:"content"`
	Module      string        `bson:"module"         json:"module"`
	Theme       string        `bson:"theme"          json:"theme"`
	ModuleID    bson.ObjectId `bson:"moduleID"       json:"moduleID"`
	ThemeID     bson.ObjectId `bson:"themeID"        json:"themeID"`
	CommentNum  int64         `bson:"commentNum"     json:"commentNum"`
	Times       int64         `bson:"times"          json:"times"`
	LastComment string        `bson:"lastComment"    json:"lastComment"`
	Created     time.Time     `bson:"created"        json:"created"`
	Image       string        `bson:"image"          json:"image"`
	IsActive    bool          `bson:"isActive"       json:"isActive"`
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
		Key:        []string{"title", "userID", "moduleID"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	articleSession = mongo.NewConnection(s, bbs.Database, CollArticle)
}

// Insert - add article.
func (sp *articleServiceProvider) Insert(article CreateArticle, userID uint64) (string, error) {
	moduleID, err := ModuleService.GetModuleID(article.Module)
	if err != nil {
		return "", err
	}

	ThemeID, err := ModuleService.GetThemeID(article.Module, article.Theme)
	if err != nil {
		return "", err
	}

	userInfo, err := user.UserServer.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	art := Article{
		Title:       article.Title,
		UserID:      userID,
		Content:     article.Content,
		Module:      article.Module,
		Theme:       article.Theme,
		ModuleID:    moduleID,
		ThemeID:     ThemeID,
		CommentNum:  0,
		Times:       0,
		LastComment: userInfo.UserName,
		Created:     time.Now(),
		Image:       article.Image,
		IsActive:      true,
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	err = conn.Insert(&art)
	if err != nil {
		return "", err
	}

	artId, err := sp.GetID(art.Title)
	err = ModuleService.UpdateArtNum(article.Module, bbs.Increase)
	if err != nil {
		return "", err
	}

	return artId.Hex(), nil
}

// GetByModuleID return articles by moduleID.
func (sp *articleServiceProvider) GetByModuleID(page int, module string) ([]Article, error) {
	var list []Article

	moduleID, err := ModuleService.GetModuleID(module)
	if err != nil {
		return list, err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"moduleID": moduleID, "isActive": true}
	err = conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByThemeID return articles by themeID.
func (sp *articleServiceProvider) GetByThemeID(page int, module, theme string) ([]Article, error) {
	var list []Article

	moduleID, err := ModuleService.GetModuleID(module)
	if err != nil {
		return list, err
	}

	themeID, err := ModuleService.GetThemeID(module, theme)
	if err != nil {
		return list, err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"moduleID": moduleID, "themeID": themeID, "isActive": true}
	err = conn.Collection().Find(query).Limit(conf.BBSConfig.Pages).Skip(page * conf.BBSConfig.Pages).All(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByTitle return articles by title.
func (sp *articleServiceProvider) GetByTitle(title string) ([]Article, error) {
	var list []Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"title": bson.M{"$regex": title, "$options": "$i"}, "isActive": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetByArtID return article by artID.
func (sp *articleServiceProvider) GetByArtID(artID bson.ObjectId) (*Article, error) {
	var list Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": artID, "isActive": true}
	err := conn.GetUniqueOne(query, &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

// GetByUserID return articles by title.
func (sp *articleServiceProvider) GetByUserID(userID uint64) ([]Article, error) {
	var list []Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"userID": userID, "isActive": true}
	err := conn.GetMany(query, &list, sort)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetID return ArtID.
func (sp *articleServiceProvider) GetID(title string) (bson.ObjectId, error) {
	var art Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"title": title}

	err := conn.GetUniqueOne(query, &art)
	if err != nil {
		return "", err
	}

	return art.Id, nil
}

// GetInfo return article's information.
func (sp *articleServiceProvider) GetInfo(artID bson.ObjectId) (*Article, error) {
	var article Article

	conn := articleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": artID}
	err := conn.GetUniqueOne(query, &article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// Delete deletes article.
func (sp *articleServiceProvider) Delete(title string) error {
	artID, err := sp.GetID(title)
	if err != nil {
		return err
	}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"isActive": false}}
	err = conn.Update(bson.M{"_id": artID}, updater)
	if err != nil {
		return err
	}

	art, err := sp.GetInfo(artID)
	if err != nil {
		return err
	}

	module, err := ModuleService.ListInfo(art.ModuleID.Hex())
	if err != nil {
		return err
	}

	return ModuleService.UpdateArtNum(module.Name, bbs.Decrease)
}

// UpdateCommentNum update the commentNum.
func (sp *articleServiceProvider) UpdateCommentNum(artID bson.ObjectId, sort int) error {
	var updater = bson.M{}

	updater = bson.M{"$inc": bson.M{"commentNum": sort}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": artID}, updater)
}

//  UpdateTimes update times.
func (sp *articleServiceProvider) UpdateTimes(num int64, artID string) error {
	updater := bson.M{"$set": bson.M{"times": num}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artID), "isActive": true}, updater)
}

// DeleteByModule deletes articles by deleting module.
func (sp *articleServiceProvider) DeleteByModule(moduleID string) error {
	updater := bson.M{"$set": bson.M{"isActive": false}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	_, err := conn.Collection().UpdateAll(bson.M{"moduleID": bson.ObjectIdHex(moduleID)}, updater)
	return err
}

// DeleteByTheme deletes articles by deleting themes.
func (sp *articleServiceProvider) DeleteByTheme(moduleID, themeID string) error {
	updater := bson.M{"$set": bson.M{"isActive": false}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	_, err := conn.Collection().UpdateAll(bson.M{"moduleID": bson.ObjectIdHex(moduleID), "themeID": bson.ObjectIdHex(themeID)}, updater)
	return err
}

// UpdateLastComment update lastComment.
func (sp *articleServiceProvider) UpdateLastComment(artID, user string) error {
	updater := bson.M{"$set": bson.M{"lastComment": user}}

	conn := articleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": bson.ObjectIdHex(artID), "isActive": true}, updater)
}
