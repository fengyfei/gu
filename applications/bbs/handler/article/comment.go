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
	jwtgo "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/bbs/article"
)

type (
	commentID struct {
		CommentID string `json:"commentID"`
	}
)

// AddComment create comment.
func AddComment(this *server.Context) error {
	var req article.CreateComment

	if err := this.JSONBody(&req); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	userID := this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64)
	req.CreatorID = uint32(userID)

	err := article.CommentService.Create(req)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// DeleteComment delete comment.
func DeleteComment(this *server.Context) error {
	var commentID commentID

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

// CommentInfo return comment's information
func CommentInfo(this *server.Context) error {
	var commentID commentID

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

	list, err := article.CommentService.UserReply(user.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// GetByArticle return comments by articleId.
func GetByArticle(this *server.Context) error {
	var artID struct {
		ArtID string `json:"artID"`
	}

	if err := this.JSONBody(&artID); err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(artID.ArtID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.CommentService.GetByArtID(artID.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}
