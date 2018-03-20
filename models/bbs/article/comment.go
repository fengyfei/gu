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

	"github.com/fengyfei/gu/applications/bbs/conf"
	"github.com/fengyfei/gu/libs/mongo"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/user"
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
		ArtID     bson.ObjectId `bson:"artID"          json:"artID"`
		CreatorID uint32        `bson:"creatorID"      json:"creatorID"`
		Creator   string        `bson:"creator"        json:"creator"`
		ReplierID uint32        `bson:"replierID"      json:"replierID"`
		Replier   string        `bson:"replier"        json:"replier"`
		ParentID  bson.ObjectId `bson:"parentID"       json:"parentID"`
		Content   string        `bson:"content"        json:"content"`
		Created   string        `bson:"created"        json:"created"`
		Status    int8          `bson:"status"         json:"status"`
	}

	// CreateComment represents the article information when created.
	CreateComment struct {
		Creator   string `json:"creator"`
		CreatorID uint32 `json:"creatorID"`
		ReplierID uint32 `json:"replierID"`
		ParentID  string `json:"parentID"`
		ArtID     string `json:"artID"`
		Content   string `json:"content"`
		Created   string `json:"created"`
	}

	// ShowComment return the comment's information which is showed to user.
	ShowComment struct {
		ID        bson.ObjectId `json:"ID"`
		CreatorID uint32        `json:"creatorID"`
		Creator   string        `json:"creator"`
		ReplierID uint32        `json:"replierID"`
		Replier   string        `json:"replier"`
		Content   string        `json:"content"`
		Created   string        `json:"created"`
		SubComms  []Comment     `json:"subComms"`
	}

	// CreateReply return the information when inserting comment.
	CreateReply struct {
		CreatorID uint32 `json:"creatorID"`
		Creator   string `json:"creator"`
		ReplierID uint32 `json:"replierID"`
		Replier   string `json:"replier"`
	}
)

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
func (sp *commentServiceProvider) Create(con orm.Connection, comment CreateComment) (*CreateReply, error) {
	var (
		comm Comment
		info *CreateReply
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	//creator, err := user.UserServer.GetUserByID(con, comment.CreatorID)
	//if err != nil {
	//	return nil, err
	//}

	replier, err := user.UserServer.GetUserByID(con, comment.ReplierID)
	if err != nil {
		return nil, err
	}

	comm = Comment{
		CreatorID: comment.CreatorID,
		//Creator:   creator.UserName,
		Creator:   comment.Creator,
		ArtID:     bson.ObjectIdHex(comment.ArtID),
		Content:   comment.Content,
		Created:   comment.Created,
		ReplierID: comment.ReplierID,
		ParentID:  bson.ObjectIdHex(comment.ParentID),
		Replier:   replier.UserName,
		Status:    bbs.CommentUnread,
	}

	err = conn.Insert(&comm)
	if err != nil {
		return nil, err
	}

	err = ArticleService.UpdateCommentNum(comm.ArtID, bbs.Increase)
	if err != nil {
		return nil, err
	}

	//err = ArticleService.UpdateLastComment(comment.ArtID, creator.UserName)
	err = ArticleService.UpdateLastComment(comment.ArtID, comment.Creator)
	if err != nil {
		return nil, err
	}

	info = &CreateReply{
		CreatorID: comm.CreatorID,
		ReplierID: comm.ReplierID,
		Creator:   comm.Creator,
		Replier:   comm.Replier,
	}

	return info, nil
}

// Delete delete comment.
func (sp *commentServiceProvider) Delete(commentID bson.ObjectId) error {
	var (
		last Comment
	)

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

	err = ArticleService.UpdateCommentNum(comment.ArtID, bbs.Decrease)
	if err != nil {
		return err
	}

	query := bson.M{"artID": comment.ArtID, "status": bson.M{"$ne": bbs.CommentDeleted}}
	err = conn.Collection().Find(query).Sort("-created").One(&last)
	if err != nil {
		return err
	}

	return ArticleService.UpdateLastComment(last.ArtID.Hex(), last.Creator)
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
func (sp *commentServiceProvider) GetByArtID(artID string) ([]*ShowComment, error) {
	var (
		comments []Comment
		list     = make([]*ShowComment, len(comments))
	)

	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"artID": bson.ObjectIdHex(artID), "status": bson.M{"$ne": bbs.CommentDeleted}, "parentID": bson.ObjectIdHex(artID)}
	sort := "-created"

	err := conn.GetMany(query, &comments, sort)
	if err != nil {
		return nil, err
	}

	for i, comment := range comments {
		subcomment, err := sp.SubComment(comment.Id)
		if err != nil {
			return nil, err
		}

		list[i] = &ShowComment{
			ID:        comment.Id,
			Creator:   comment.Creator,
			CreatorID: comment.CreatorID,
			Replier:   comment.Replier,
			ReplierID: comment.ReplierID,
			Content:   comment.Content,
			Created:   comment.Created,
			SubComms:  subcomment,
		}
	}

	return list, nil
}

// GetByUserID return comments by userID
func (sp *commentServiceProvider) GetByUserID(userID uint32) ([]Comment, error) {
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

// UserReply return the information about someone's reply.
func (sp *commentServiceProvider) UserReply(userID uint32) ([]UserReply, error) {
	comments, err := sp.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	list := make([]UserReply, len(comments))
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

// CommentNum return the number of someone's comment.
func (sp *commentServiceProvider) CommentNum(userID uint32) (int, error) {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"creatorID": userID, "status": bson.M{"$ne": bbs.CommentDeleted}}

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

	query := bson.M{"replierID": userID, "status": bbs.CommentRead}

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

	query := bson.M{"replierID": userID, "status": bbs.CommentUnread}

	err := conn.GetMany(query, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// MessageRead modify the status of messsage when read.
func (sp *commentServiceProvider) MessageRead(commentID string) error {
	conn := commentSession.Connect()
	defer conn.Disconnect()

	query := bson.M{"_id": bson.ObjectIdHex(commentID)}

	return conn.Update(query, bson.M{"$set": bson.M{"status": bbs.CommentRead}})
}
