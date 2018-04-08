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
	"gopkg.in/mgo.v2/bson"

	"fmt"
	mysql "github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/bbs/article"
	"github.com/fengyfei/gu/models/user"
)

type (
	commentID struct {
		CommentID string `json:"commentid"`
	}

	userid struct {
		UserID uint32 `json:"userid"`
	}

	// createComment represents the article information when created.
	createComment struct {
		//CreatorID uint32 `json:"creatorid"`
		RepliedID uint32 `json:"repliedid"`
		ParentID  string `json:"parentid"`
		ArtID     string `json:"artid"`
		Content   string `json:"content"     validate:"required"`
		Created   string `json:"created"     validate:"required"`
	}

	// showComment return the comment's information which is showed to user.
	showComment struct {
		ID       bson.ObjectId     `json:"id"`
		Creator  string            `json:"creator"`
		Replier  string            `json:"replier"`
		Content  string            `json:"content"`
		Created  string            `json:"created"`
		SubComms []article.Comment `json:"subcomms"`
	}

	// createReply return the information when inserting comment.
	createReply struct {
		CreatorID uint32 `json:"creatorid"`
		Creator   string `json:"creator"`
		RepliedID uint32 `json:"repliedid"`
		Replier   string `json:"replier"`
	}

	// UserReply represents the information about someone's reply.
	userReply struct {
		Title    string `json:"title"`
		Creator  string `json:"creator"`
		Replier  string `json:"replier"`
		Category string `json:"category"`
		Content  string `json:"content"`
		Created  string `json:"created"`
	}
)

// AddComment create comment.
func AddComment(this *server.Context) error {
	var (
		req createComment
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	fmt.Println("%+v", req)

	if !bson.IsObjectIdHex(req.ArtID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	//userID := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64)
	userID := uint32(1000)

	creator, err := user.UserService.GetUserByID(conn, userID)
	if err != nil {
		return err
	}

	replier, err := user.UserService.GetUserByID(conn, req.RepliedID)
	if err != nil {
		return err
	}

	err = article.ArticleService.IfExist(req.ArtID)
	if err != nil {
		return err
	}

	addcomment := &article.Comment{
		CreatorID: userID,
		Creator:   creator.UserName,
		RepliedID: req.RepliedID,
		Replier:   replier.UserName,
		ParentID:  bson.ObjectIdHex(req.ParentID),
		ArtID:     bson.ObjectIdHex(req.ArtID),
		Content:   req.Content,
		Created:   req.Created,
	}

	err = article.CommentService.Create(conn, addcomment)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	fmt.Println("5555555")
	info := &createReply{
		CreatorID: userID,
		Creator:   creator.UserName,
		RepliedID: req.RepliedID,
		Replier:   replier.UserName,
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, info)
}

// DeleteComment delete comment.
func DeleteComment(this *server.Context) error {
	var (
		commentID commentID
	)

	if err := this.JSONBody(&commentID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(commentID.CommentID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CommentService.Delete(bson.ObjectIdHex(commentID.CommentID))
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// CommentInfo return comment's information.
func CommentInfo(this *server.Context) error {
	var (
		commentID commentID
	)

	if err := this.JSONBody(&commentID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(commentID.CommentID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.CommentService.ListInfo(bson.ObjectIdHex(commentID.CommentID))
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// UserReply return the information about someone's reply.
func UserReply(this *server.Context) error {
	var user struct {
		UserID uint32 `json:"userID"`
	}

	if err := this.JSONBody(&user); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	comments, err := article.CommentService.UserReply(conn, user.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	list := make([]userReply, len(comments))
	for i, comment := range comments {
		art, err := article.ArticleService.GetByArtID(comment.ArtID.Hex())
		if err != nil {
			return err
		}

		category, err := article.CategoryService.ListInfo(art.CategoryID.Hex())
		list[i] = userReply{
			Title:    art.Title,
			Creator:  comment.Creator,
			Replier:  comment.Replier,
			Category: category.Name,
			Content:  comment.Content,
			Created:  comment.Created,
		}
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByArticle return comments by articleId.
func GetByArticle(this *server.Context) error {
	var artID struct {
		ArtID string `json:"artid"`
	}

	if err := this.JSONBody(&artID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(artID.ArtID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	comments, err := article.CommentService.GetByArtID(conn, artID.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	list := make([]showComment, len(comments))
	for i, comment := range comments {
		subcomment, err := article.CommentService.SubComment(comment.Id)
		if err != nil {
			return err
		}

		list[i] = showComment{
			ID:       comment.Id,
			Creator:  comment.Creator,
			Replier:  comment.Replier,
			Content:  comment.Content,
			Created:  comment.Created,
			SubComms: subcomment,
		}
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// HistoryMessage return the message which is read by userid.
func HistoryMessage(this *server.Context) error {
	var (
		user userid
	)

	if err := this.JSONBody(&user); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.CommentService.HistoryMessage(user.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// UnreadMessage return the unread message by userid.
func UnreadMessage(this *server.Context) error {
	var (
		user userid
	)

	if err := this.JSONBody(&user); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.CommentService.UnreadMessage(user.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// MessageRead change the status of the message which is read.
func MessageRead(this *server.Context) error {
	var (
		comment commentID
	)

	if err := this.JSONBody(&comment); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.CommentService.MessageRead(comment.CommentID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}
