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
 *     Initial: 2018/03/17        Tong Yuehong
 */

package article

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/bbs"
)

type collectionServiceProvider struct{}

var (
	//CollectionService expose serviceProvider.
	CollectionService *collectionServiceProvider
	collectionSession *mongo.Connection
)

type (
	Collection struct {
		Id     bson.ObjectId   `bson:"_id,omitempty"  json:"id"`
		UserID uint32          `bson:"userID"         json:"userID"`
		ArtID  []bson.ObjectId `bson:"artID"          json:"artID"`
	}

	CreateColl struct {
		UserID uint32   `bson:"userID"         json:"userID"`
		ArtID  []string `bson:"artID"          json:"artID"`
	}


)

func init() {
	const (
		Collection = "collection"
	)

	url := conf.BBSConfig.MongoURL + "/" + bbs.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB(bbs.Database).C(Collection).EnsureIndex(mgo.Index{
		Key:        []string{"userID"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	collectionSession = mongo.NewConnection(s, bbs.Database, Collection)
}

func (sp *collectionServiceProvider) Insert(created CreateColl) error {
	var artID = make([]bson.ObjectId, len(created.ArtID))
	for i, artid := range created.ArtID {
		artID[i] = bson.ObjectIdHex(artid)
	}

	conn := collectionSession.Connect()
	defer conn.Disconnect()

	collection := &Collection{
		UserID: created.UserID,
		ArtID:  artID,
	}

	query := bson.M{"userID": created.UserID}
	num, err := conn.Collection().Find(query).Count()
	if err != nil {
		return err
	}

	if num != 0 {
		err = conn.Update(query, bson.M{"$push": bson.M{"artID": bson.M{"$each": artID}}})
		if err != nil {
			return err
		}
	} else {
		err = conn.Insert(collection)
	}

	return nil
}

// UnCollect cancel collection.
func (sp *collectionServiceProvider) UnCollect(userID uint32, artID string) error {
	conn := collectionSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"userID": userID}
	return conn.Update(query, bson.M{"$pull": bson.M{"artID": bson.ObjectIdHex(artID)}})
}
