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
)

type moduleserviceProvider struct{}

var (
	// ErrNotFound - No result found
	ErrMDNotFound = errors.New("No result found")
	// Service expose serviceProvider
	ModuleService *moduleserviceProvider
	modulesession *mongo.Connection
)

// Theme represents the second category.
type Theme struct {
	Id   bson.ObjectId `bson:"id"         json:"id"`
	Name string        `bson:"name"        json:"name"`
}

// Module represents the module information.
type Module struct {
	Id         bson.ObjectId `bson:"_id"         json:"id"`
	Name       string        `bson:"name"        json:"name"`
	ArtNum     int64         `bson:"artNum"      json:"artNum"`
	ModuleView int64         `bson:"moduleView"  json:"moduleView"`
	Recommand  int           `bson:"recommend"   json:"recommend"`
	Themes     []Theme       `bson:"themes"      json:"themes"`
	Status     bool          `bson:"status"      json:"status"`
}

// CreateModule represents the module information when created.
type CreateModule struct {
	Name string `json:"name"`
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

	modulesession = mongo.NewConnection(s, bbs.Database, cname)
}

// GetModuleId gets moduleId by name
func (sp *moduleserviceProvider) GetModuleID(name string) (bson.ObjectId, error) {
	var module Module

	conn := modulesession.Connect()
	defer conn.Disconnect()

	query := bson.M{"name": name}

	err := conn.GetUniqueOne(query, &module)

	return module.Id, err
}

// GetModuleId gets moduleId by name
func (sp *moduleserviceProvider) GetThemeID(moduleName, themeName string) (string, error) {
	var module Module

	conn := modulesession.Connect()
	defer conn.Disconnect()

	query := bson.M{"name": moduleName, "themes.name": themeName}

	err := conn.Collection().Find(query).Select(bson.M{"themes.$": 1}).One(&module)
	if len(module.Themes) == 0 {
		return "", ErrMDNotFound
	}

	return module.Themes[0].Id.Hex(), err
}

// add article
func (sp *moduleserviceProvider) CreateModule(module CreateModule) error {
	mod := Module{
		Id:         bson.NewObjectId(),
		Name:       module.Name,
		ArtNum:     0,
		ModuleView: 0,
		Status:     true,
	}
	conn := modulesession.Connect()
	err := conn.Insert(&mod)
	return err
}

// add module
func (sp *moduleserviceProvider) CreateTheme(module, theme string) error {
	moduleId, err := sp.GetModuleID(module)

	t := Theme{
		Id:   bson.NewObjectId(),
		Name: theme,
	}

	updater := bson.M{"$addToSet": bson.M{"themes": t}}

	conn := modulesession.Connect()
	defer conn.Disconnect()

	err = conn.Update(bson.M{"_id": moduleId}, updater)
	return err
}

// Update the artNum of the module
func (sp *moduleserviceProvider) UpdateArtNum(module string) error {
	moduleId, err := sp.GetModuleID(module)

	if err != nil {
		return err
	}

	updater := bson.M{"$inc": bson.M{"ArtNum": 1}}

	conn := modulesession.Connect()
	defer conn.Disconnect()

	err = conn.Update(bson.M{"_id": moduleId}, updater)
	return err
}

// Update ModuleView
func (sp *moduleserviceProvider) UpdateModuleView(num int64, module string) error {
	moduleId, err := sp.GetModuleID(module)
	if err != nil {
		return err
	}
	updater := bson.M{"$set": bson.M{"ModuleView": num}}

	conn := modulesession.Connect()
	defer conn.Disconnect()

	err = conn.Update(bson.M{"_id": moduleId}, updater)
	return err
}
