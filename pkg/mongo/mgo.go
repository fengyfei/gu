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
 *     Initial: 2017/10/24        Jia Chenhui
 */

package mongo

import (
	"github.com/fengyfei/nuts/mgo/copy"
	"gopkg.in/mgo.v2"

	"github.com/fengyfei/gu/pkg/log"
)

// Session represents a communication session with the database.
type Session struct {
	CollInfo *copy.CollectionInfo
}

// InitMDSess establishes a new session to the cluster.
func InitMDSess(url, db, coll string, index *mgo.Index) *Session {
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	log.GlobalLogReporter.Debug("The MongoDB of blog server connected.")

	s.SetMode(mgo.Monotonic, true)

	collInfo := &copy.CollectionInfo{
		Session:    s,
		Database:   db,
		Collection: coll,
		Index:      index,
	}

	return &Session{
		CollInfo: collInfo,
	}
}
