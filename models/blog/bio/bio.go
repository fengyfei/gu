/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 *     Initial: 2017/10/27        Jia Chenhui
 */

package bio

import (
	"github.com/astaxie/beego"
	"github.com/fengyfei/nuts/mgo/copy"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/common"
	"github.com/fengyfei/gu/pkg/mongo"
)

type serviceProvider struct{}

var (
	// Service expose serviceProvider
	Service  *serviceProvider
	mdSess   *mongo.Session
	uniqueID bson.ObjectId
)

// Prepare initializing database.
func Prepare() {
	url := beego.AppConfig.String("mongo::url") + "/" + common.MDBlogDName

	mdSess = mongo.InitMDSess(url, common.MDBlogDName, common.MDBioColl, nil)
	Service = &serviceProvider{}
	uniqueID = bson.NewObjectId()
}

type MDBio struct {
	BioID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Title string        `bson:"Title" json:"title"`
	Bio   string        `bson:"Bio" json:"bio"`
}

// MDCreateBio use to create article.
type MDCreateBio struct {
	Title string
	Bio   string
}

// GetBio get bio based on the unique ID.
func (sp *serviceProvider) GetBio() (MDBio, error) {
	var (
		bioInfo MDBio
		err     error
	)

	err = copy.GetByID(mdSess.CollInfo, uniqueID, &bioInfo)

	return bioInfo, err
}

// Create create bio based on the unique ID.
func (sp *serviceProvider) Create(bio *MDCreateBio) error {
	bioInfo := MDBio{
		BioID: uniqueID,
		Title: bio.Title,
		Bio:   bio.Bio,
	}

	selector := bson.M{"_id": uniqueID}
	_, err := copy.Upsert(mdSess.CollInfo, selector, &bioInfo)

	return err
}
