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
 *     Initial: 2018/01/25        Tong Yuehong
 */

package article

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/bbs"
	"fmt"
)

type moduleServiceProvider struct{}

var (
	// ErrMDNotFound - No result found
	ErrMDNotFound = errors.New("No result found")
	// ModuleService expose serviceProvider
	ModuleService *moduleServiceProvider
	moduleSession *mongo.Connection
)

// Theme represents the second category.
type Theme struct {
	Id     bson.ObjectId `bson:"id"         json:"id"`
	Name   string        `bson:"name"       json:"name"  validate:"required"`
	Status bool          `bson:"status"     json:"status"`
}

// Module represents the module information.
type Module struct {
	Id         bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Name       string        `bson:"name"            json:"name"`
	ArtNum     int64         `bson:"artNum"          json:"artNum"`
	ModuleView int64         `bson:"moduleView"      json:"moduleView"`
	Recommend  int           `bson:"recommend"       json:"recommend"`
	Themes     []Theme       `bson:"themes"          json:"themes"`
	Status     bool          `bson:"status"          json:"status"`
}

// CreateModule represents the module information when created.
type CreateModule struct {
	Name string `json:"name" validate:"required"`
}

func init() {
	const (
		cname = "module"
	)

	url := conf.BBSConfig.MongoURL + "/" + bbs.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	moduleSession = mongo.NewConnection(s, bbs.Database, cname)
}

// GetModuleID gets moduleId by name.
func (sp *moduleServiceProvider) GetModuleID(name string) (bson.ObjectId, error) {
	var module Module

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"name": name}

	err := conn.GetUniqueOne(query, &module)

	return module.Id, err
}

// GetThemeID gets moduleId by name.
func (sp *moduleServiceProvider) GetThemeID(moduleName, themeName string) (bson.ObjectId, error) {
	var module Module

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"name": moduleName, "themes.name": themeName}

	err := conn.Collection().Find(query).Select(bson.M{"themes.$": 1}).One(&module)
	if err != nil {
		return "", err
	} else if len(module.Themes) == 0 {
		return "", ErrMDNotFound
	}

	return module.Themes[0].Id, nil
}

// CreateModule add module.
func (sp *moduleServiceProvider) CreateModule(module CreateModule) error {
	mod := Module{
		Name:       module.Name,
		ArtNum:     0,
		ModuleView: 0,
		Status:     true,
	}

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	return conn.Insert(&mod)
}

// CreateTheme add theme.
func (sp *moduleServiceProvider) CreateTheme(module, theme string) error {
	moduleId, err := sp.GetModuleID(module)
	if err != nil {
		return err
	}

	t := Theme{
		Id:     bson.NewObjectId(),
		Name:   theme,
		Status: true,
	}

	updater := bson.M{"$addToSet": bson.M{"themes": t}}

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": moduleId}, updater)
}

// UpdateArtNum update the artNum of the module.
func (sp *moduleServiceProvider) UpdateArtNum(module string, sort int) error {
	var updater interface{}

	moduleId, err := sp.GetModuleID(module)
	if err != nil {
		return err
	}

	updater = bson.M{"$inc": bson.M{"artNum": sort}}

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	err = conn.Update(bson.M{"_id": moduleId}, updater)
	return err
}

//  UpdateModuleView update ModuleView.
func (sp *moduleServiceProvider) UpdateModuleView(num int64, module string) error {
	moduleId, err := sp.GetModuleID(module)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"moduleView": num}}

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	return conn.Update(bson.M{"_id": moduleId, "status": true}, updater)
}

// ListInfo gets module's information.
func (sp *moduleServiceProvider) ListInfo(moduleId string) (Module, error) {
	var module Module

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	fmt.Println("aaaaa", moduleId)
	query := bson.M{"_id": bson.ObjectIdHex(moduleId), "status": true}
	err := conn.GetUniqueOne(query, &module)
	if err != nil {
		return Module{}, err
	}

	return module, err
}

// AllModules gets all modules.
func (sp *moduleServiceProvider) AllModules() ([]Module, error) {
	var module []Module

	conn := moduleSession.Connect()
	defer conn.Disconnect()

	sort := "-Created"

	query := bson.M{"status": true}
	err := conn.GetMany(query, &module, sort)
	if err != nil {
		return nil, err
	}

	return module, err
}

// DeleteModule delete module.
func (sp *moduleServiceProvider) DeleteModule(moduleId string) error {
	conn := moduleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"status": false}}
	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(moduleId)}, updater)
	if err != nil {
		return err
	}

	return ArticleService.DeleteByModule(moduleId)
}

// DeleteTheme delete theme.
func (sp *moduleServiceProvider) DeleteTheme(moduleId, themeId string) error {
	conn := moduleSession.Connect()
	defer conn.Disconnect()

	updater := bson.M{"$set": bson.M{"themes.$.status": false}}
	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(moduleId), "themes.id": bson.ObjectIdHex(themeId)}, updater)
	if err != nil {
		return err
	}

	return ArticleService.DeleteByTheme(moduleId, themeId)
}
