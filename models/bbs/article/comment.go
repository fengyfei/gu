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
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/user"
	"github.com/fengyfei/gu/applications/bbs/initialize"
)

type commentServiceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	CommentService *commentServiceProvider
	commentSession *mongo.Connection
	conn           orm.Connection
)

// Comment represents the comment information.
type Comment struct {
	Id        bson.ObjectId `bson:"_id,omitempty"  json:"id"`
	ArtId     bson.ObjectId `bson:"artId"          json:"artId"`
	CreatorId uint64        `bson:"creatorId"      json:"creatorId"`
	Creator   string        `bson:"creator"        json:"creator"`
	ReplierId uint64        `bson:"replierId"      json:"replierId"`
	Replier   string        `bson:"replier"        json:"replier"`
	ParentId  bson.ObjectId `bson:"parentId"       json:"parentId"`
	Content   string        `bson:"content"        json:"content"`
	Created   time.Time     `bson:"created"        json:"created"`
	Status    bool          `bson:"status"         json:"status"`
}

// CreateComment represents the article information when created.
type CreateComment struct {
	CreatorId uint64 `json:"creatorId"`
	ReplierId uint64 `json:"replierId"`
	ParentId  string `json:"parentId"`
	ArtId     string `json:"artId"`
	Content   string `json:"content"`
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

	if err != nil {
		panic(err)
	}
}

// Create - insert comment.
func (sp *commentServiceProvider) Create(comment CreateComment) error {
	c , _ := initialize.Pool.Get()
	creator, err := user.UserServer.GetUserByID(c, comment.CreatorId)
	if err != nil {
		return err
	}

	replier, err := user.UserServer.GetUserByID(c, comment.CreatorId)
	if err != nil {
		return err
	}

	comm := Comment{
		CreatorId: comment.CreatorId,
		Creator:   creator.Username,
		ReplierId: comment.ReplierId,
		Replier:   replier.Username,
		ParentId:  bson.ObjectIdHex(comment.ParentId),
		ArtId:     bson.ObjectIdHex(comment.ArtId),
		Content:   comment.Content,
		Created:   time.Now(),
		Status:    true,
	}

	conn := commentSession.Connect()
	defer conn.Disconnect()

	err = conn.Insert(&comm)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(comm.ArtId, bbs.IncCount)
	if err != nil {
		return err
	}

	return ArticleService.UpdateLastComment(comment.ArtId, creator.Username)
}

// Delete - delete comment.
func (sp *commentServiceProvider) Delete(commentId bson.ObjectId) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	comment, err := sp.ListInfo(commentId)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"status": false}}
	err = conn.Update(bson.M{"_id": commentId}, updater)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(comment.ArtId, bbs.DecCount)
	if err != nil {
		return err
	}

	var last Comment
	query := bson.M{"artId": comment.ArtId, "status": true}
	err = conn.Collection().Find(query).Sort("-created").One(&last)
	if err != nil {
		return err
	}

	return ArticleService.UpdateLastComment(last.ArtId.Hex(), last.Creator)
}

// ListInfo get comment's information.
func (sp *commentServiceProvider) ListInfo(commentId bson.ObjectId) (Comment, error) {
	var comment Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": commentId, "status": true}
	err := conn.GetUniqueOne(query, &comment)
	if err != nil {
		return Comment{}, err
	}

	return comment, err
}

// GetByArtId get comments by artId
func (sp *commentServiceProvider) GetByArtId(artId string) ([]Comment, error) {
	var comments []Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"artId": artId, "status": true}
	sort := "-Created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, err
}

// GetByUserId get comments by userId
func (sp *commentServiceProvider) GetByUserId(userId uint64) ([]Comment, error) {
	var comments []Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"creatorId": userId, "status": true}
	sort := "-Created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, err
}

// UserReply shows the information about someone's reply.
func (sp *commentServiceProvider) UserReply(userId uint64) ([]UserReply, error) {
	comments, err := sp.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var list = make([]UserReply, len(comments))
	for i, comment := range comments {
		art, err := ArticleService.GetByArtId(comment.ArtId)
		if err != nil {
			return nil, err
		}

		list[i] = UserReply{
			Title:   art.Title,
			Creator: comment.Creator,
			Replier: comment.Replier,
			Module:  art.Module,
			Content: comment.Content,
			Created: comment.Created,
		}
	}

	return list, nil
}
