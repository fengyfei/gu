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
 *     Initial: 2018/01/28        Tong Yuehong
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

type commentServiceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	CommentService *commentServiceProvider
	commentSession *mongo.Connection
)

// Comment represents the comment information.
type Comment struct {
	Id        bson.ObjectId `bson:"_id,omitempty"  json:"id"`
	CreatedId uint64        `bson:"createdId"      json:"createdId"`
	ReplyId   uint64        `bson:"replyId"        json:"replyId"`
	ParentId  bson.ObjectId `bson:"parentId"       json:"parentId"`
	ArtId     bson.ObjectId `bson:"artId"          json:"artId"`
	Content   string        `bson:"content"        json:"content"`
	Created   time.Time     `bson:"created"        json:"created"`
	Status    bool          `bson:"status"         json:"status"`
}

// CreateComment represents the article information when created.
type CreateComment struct {
	CreatedId uint64        `json:"createdId"`
	ReplyId   uint64        `json:"replyId"`
	ParentId  string        `json:"parentId"`
	ArtId     string        `json:"artId"`
	Content   string        `json:"content"`
}

func init() {
	const (
		CollComment = "comment"
	)

	url := conf.BBSConfig.MongoURL + "/" + bbs.Database
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	commentSession = mongo.NewConnection(s, bbs.Database, CollComment)
}

// Create - insert comment.
func (sp *commentServiceProvider) Create(comment CreateComment) error {
	comm := Comment{
		CreatedId: comment.CreatedId,
		ReplyId:   comment.ReplyId,
		ParentId:  bson.ObjectIdHex(comment.ParentId),
		ArtId:     bson.ObjectIdHex(comment.ArtId),
		Content:   comment.Content,
		Created:   time.Now(),
		Status:    true,
	}

	conn := commentSession.Connect()
	defer conn.Disconnect()

	err := conn.Insert(&comm)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(comm.ArtId, "add")
	return err
}

// Delete - delete comment.
func (sp *commentServiceProvider) Delete(commentId bson.ObjectId) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	c ,err := sp.GetInfo(commentId)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"status": false}}
	err = conn.Update(bson.M{"_id": commentId}, updater)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(c.ArtId, "sub")
	return err
}

// GetInfo get comment's information.
func (sp *commentServiceProvider) GetInfo(commentId bson.ObjectId) (Comment, error) {
	var comment Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": commentId, "status":true}
	err := conn.GetUniqueOne(query, &comment)

	return comment, err
}
