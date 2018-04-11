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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/bbs"
)

type commentServiceProvider struct{}

var (
	// ArticleService expose serviceProvider.
	CommentService *commentServiceProvider
	commentSession *mongo.Connection
)

type (
	// Comment represents the comment information.
	Comment struct {
		Id        bson.ObjectId `bson:"_id,omitempty"  json:"id"`
		ArtID     bson.ObjectId `bson:"artID"          json:"artid"`
		Content   string        `bson:"content"        json:"content"`
		CreatorID uint32        `bson:"creatorID"      json:"creatorid"`
		Creator   string        `bson:"creator"        json:"creator"`
		RepliedID uint32        `bson:"repliedID"      json:"repliedid"`
		Replier   string        `bson:"replier"        json:"replier"`
		ParentID  bson.ObjectId `bson:"parentID"       json:"parentid"`
		Created   string        `bson:"created"        json:"created"`
		Status    int8          `bson:"status"         json:"status"`
	}
)

func init() {
	const (
		cname = "comment"
	)

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"creatorID"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	})

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"artID"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	})

	initialize.S.DB(bbs.Database).C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"repliedID"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	})

	commentSession = mongo.NewConnection(initialize.S, bbs.Database, cname)
}

// Create insert comment.
func (sp *commentServiceProvider) Create(con orm.Connection, comment *Comment) error {
	var (
		comm Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	//creator, err := user.UserService.GetUserByID(con, comment.CreatorID)
	//if err != nil {
	//	return err
	//}
	//
	//replier, err := user.UserService.GetUserByID(con, comment.RepliedID)
	//if err != nil {
	//	return err
	//}

	comm = Comment{
		ArtID:     comment.ArtID,
		Content:   comment.Content,
		CreatorID: comment.CreatorID,
		Creator:   comment.Creator,
		RepliedID: comment.RepliedID,
		Replier:   comment.Replier,
		ParentID:  comment.ParentID,
		Created:   comment.Created,
		Status:    bbs.CommentUnread,
	}

	err := conn.Insert(&comm)
	if err != nil {
		return err
	}

	return ArticleService.UpdateCommentNum(comm.ArtID, bbs.Increase)
}

// Delete delete comment.
func (sp *commentServiceProvider) Delete(commentID bson.ObjectId) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	comment, err := sp.ListInfo(commentID)
	if err != nil {
		return err
	}

	updater := bson.M{"$set": bson.M{"status": bbs.CommentDeleted}}
	err = conn.Update(bson.M{"_id": commentID}, updater)
	if err != nil {
		return err
	}

	return ArticleService.UpdateCommentNum(comment.ArtID, bbs.Decrease)
}

// ListInfo return comment's information.
func (sp *commentServiceProvider) ListInfo(commentID bson.ObjectId) (*Comment, error) {
	var (
		comment Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": commentID, "status": bson.M{"$ne": bbs.CommentDeleted}}
	err := conn.GetUniqueOne(query, &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetByArtID return comments by artID
func (sp *commentServiceProvider) GetByArtID(con orm.Connection, artID string) ([]Comment, error) {
	var (
		comments []Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"artID": bson.ObjectIdHex(artID), "status": bson.M{"$ne": bbs.CommentDeleted}, "parentID": bson.ObjectIdHex(artID)}
	sort := "-created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// UserReply return the information about someone's reply.
func (sp *commentServiceProvider) UserReply(con orm.Connection, userID uint32) ([]Comment, error) {
	var (
		comments []Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"creatorID": userID, "status": bson.M{"$ne": bbs.CommentDeleted}}
	sort := "created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// CommentNum return the number of someone's comment.
func (sp *commentServiceProvider) CommentNum(userID uint32) (int, error) {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"repliedID": userID, "status": bson.M{"$ne": bbs.CommentDeleted}}

	userComment, err := conn.Collection().Find(query).Count()
	if err != nil {
		return 0, err
	}

	return userComment, err
}

// SubComment returns the subComments of the mainComment.
func (sp *commentServiceProvider) SubComment(mainCommentID bson.ObjectId) ([]Comment, error) {
	var (
		comments []Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"parentID": mainCommentID, "status": bson.M{"$ne": bbs.CommentDeleted}}
	sort := "-created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	return comments, err
}

// HistoryMessage returns the message which is read.
func (sp *commentServiceProvider) HistoryMessage(userID uint32) ([]Comment, error) {
	var (
		list []Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"repliedID": userID, "status": bbs.CommentRead}

	err := conn.GetMany(query, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// UnreadMessage  returns the messages which are not read.
func (sp *commentServiceProvider) UnreadMessage(userID uint32) ([]Comment, error) {
	var (
		list []Comment
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"repliedID": userID, "status": bbs.CommentUnread}

	err := conn.GetMany(query, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// MessageRead modify the status of message when read.
func (sp *commentServiceProvider) MessageRead(commentID string) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(commentID)}

	return conn.Update(query, bson.M{"$set": bson.M{"status": bbs.CommentRead}})
}

// LastComment returns the last comment of the article.
func (sp *commentServiceProvider) LastComment(artID string) (*Comment, error) {
	var (
		lastComment Comment
	)
	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"artID": bson.ObjectIdHex(artID), "status": bson.M{"$ne": bbs.CommentDeleted}}
	err := conn.Collection().Find(query).Sort("created").One(&lastComment)
	if err != nil {
		return nil, err
	}

	return &lastComment, nil
}

// IsExist checks whether the id exists.
func (sp *commentServiceProvider) IfExist(id string) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	num, err := conn.Collection().Find(bson.M{"_id": bson.ObjectIdHex(id)}).Count()
	if err != nil {
		return err
	}

	if num == 0 {
		return ErrNotFound
	}

	return nil
}

func (sp *commentServiceProvider) NumByArt(artid string) (int, error) {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	num, err := conn.Collection().Find(bson.M{"artID": bson.ObjectIdHex(artid), "active": true}).Count()
	if err != nil {
		return 0, err
	}

	return num, err
}
