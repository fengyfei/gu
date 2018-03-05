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
	mysql "github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/user"
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
	ArtID     bson.ObjectId `bson:"artID"          json:"artID"`
	CreatorID uint32        `bson:"creatorID"      json:"creatorID"`
	Creator   string        `bson:"creator"        json:"creator"`
	ReplierID uint32        `bson:"replierID"      json:"replierID"`
	Replier   string        `bson:"replier"        json:"replier"`
	ParentID  bson.ObjectId `bson:"parentID"       json:"parentID"`
	Content   string        `bson:"content"        json:"content"`
	Created   time.Time     `bson:"created"        json:"created"`
	IsActive  bool          `bson:"isActive"       json:"isActive"`
}

// CreateComment represents the article information when created.
type CreateComment struct {
	CreatorID uint32 `json:"creatorID"`
	ReplierID uint32 `json:"replierID"`
	ParentID  string `json:"parentID"`
	ArtID     string `json:"artID"`
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
}

// Create insert comment.
func (sp *commentServiceProvider) Create(comment CreateComment) error {
	con, err := mysql.Pool.Get()
	defer mysql.Pool.Release(con)

	creator, err := user.UserServer.GetUserByID(con, comment.CreatorID)
	if err != nil {
		return err
	}

	replier, err := user.UserServer.GetUserByID(con, comment.CreatorID)
	if err != nil {
		return err
	}

	comm := Comment{
		CreatorID: comment.CreatorID,
		Creator:   creator.UserName,
		ReplierID: comment.ReplierID,
		Replier:   replier.UserName,
		ParentID:  bson.ObjectIdHex(comment.ParentID),
		ArtID:     bson.ObjectIdHex(comment.ArtID),
		Content:   comment.Content,
		Created:   time.Now(),
		IsActive:  true,
	}

	conn := commentSession.Connect()
	defer conn.Disconnect()

	err = conn.Insert(&comm)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(comm.ArtID, bbs.Increase)
	if err != nil {
		return err
	}

	return ArticleService.UpdateLastComment(comment.ArtID, creator.UserName)
}

// Delete delete comment.
func (sp *commentServiceProvider) Delete(commentID bson.ObjectId) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	comment, err := sp.ListInfo(commentID)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"isActive": false}}
	err = conn.Update(bson.M{"_id": commentID}, updater)
	if err != nil {
		return err
	}

	err = ArticleService.UpdateCommentNum(comment.ArtID, bbs.Decrease)
	if err != nil {
		return err
	}

	var last Comment
	query := bson.M{"artID": comment.ArtID, "isActive": true}
	err = conn.Collection().Find(query).Sort("-created").One(&last)
	if err != nil {
		return err
	}

	return ArticleService.UpdateLastComment(last.ArtID.Hex(), last.Creator)
}

// ListInfo return comment's information.
func (sp *commentServiceProvider) ListInfo(commentID bson.ObjectId) (*Comment, error) {
	var comment Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": commentID, "isActive": true}
	err := conn.GetUniqueOne(query, &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetByArtID return comments by artID
func (sp *commentServiceProvider) GetByArtID(artID string) ([]Comment, error) {
	var comments []Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"artID": artID, "isActive": true}
	sort := "-Created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// GetByUserID return comments by userID
func (sp *commentServiceProvider) GetByUserID(userID uint32) ([]Comment, error) {
	var comments []Comment

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"creatorID": userID, "isActive": true}
	sort := "-Created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// UserReply return the information about someone's reply.
func (sp *commentServiceProvider) UserReply(userID uint32) ([]UserReply, error) {
	comments, err := sp.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var list = make([]UserReply, len(comments))
	for i, comment := range comments {
		art, err := ArticleService.GetByArtID(comment.ArtID)
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
