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
 *     Initial: 2018/03/26        Chen Yanchen
 */

package project

import (
	"time"

	"github.com/fengyfei/gu/libs/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/blog/conf"
	"github.com/fengyfei/gu/models/blog"
)

type projectProviceServer struct{}

var ProjectServer *projectProviceServer

type (
	Project struct {
		ID      bson.ObjectId `bson:"_id,omitempty"`
		Title   string        `bson:"title"`
		Author  string        `bson:"author"`
		Detail  string        `bson:"detail"`
		Link    string        `bson:"link"`
		Image   string        `bson:"image"`
		Created time.Time     `bson:"created"`
		Status  int8          `bson:"status"`
	}
)

var session *mongo.Connection

func init() {
	const cname = "project"

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

// Creat a project's document.
func (sp *projectProviceServer) Creat(p *Project) error {
	conn := session.Connect()
	defer conn.Disconnect()

	p.Created = time.Now()

	err := conn.Insert(&p)
	return err
}

// Delete modify the status.
func (sp *projectProviceServer) Delete(id string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"status": blog.Delete}})
	return err
}

func (sp *projectProviceServer) Modify(p *Project) error {
	conn := session.Connect()
	defer conn.Disconnect()

	err := conn.Update(bson.M{"_id": p.ID}, &p)
	return err
}

// GetID use 'title' find out ID.
func (sp *projectProviceServer) GetID(title string) (bson.ObjectId, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var p Project
	err := conn.GetUniqueOne(bson.M{"title": title}, &p)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

// GetByID get project by ID.
func (sp *projectProviceServer) GetByID(id bson.ObjectId) (*Project, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var p Project
	err := conn.GetByID(id, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// AbstractList get all approval project.
func (sp *projectProviceServer) AbstractList() ([]Project, error) {
	conn := session.Connect()
	defer conn.Disconnect()

	var res []Project
	err := conn.GetMany(bson.M{"status": blog.Approval}, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
