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
 *     Initial: 2018/01/24        Tong Yuehong
 */

package article

import (
	"strings"

	"gopkg.in/mgo.v2/bson"

	mysql "github.com/fengyfei/gu/applications/bbs/initialize"
	"github.com/fengyfei/gu/applications/bbs/util"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/bbs"
	"github.com/fengyfei/gu/models/bbs/article"
	"github.com/fengyfei/gu/models/user"
)

type (
	title struct {
		Title string `json:"title" validate:"required"`
	}

	createArticle struct {
		Title      string `json:"title"       validate:"required,max=50,min=20"`
		Content    string `json:"content"     validate:"required,min=20"`
		AuthorID   uint32 `json:"authorid"`
		CategoryID string `json:"categoryid"  validate:"required"`
		TagID      string `json:"tagid"       validate:"required"`
		Image      string `json:"image"`
		Created    string `json:"created"     validate:"required"`
	}

	replyInfo struct {
		Id         bson.ObjectId `json:"id"`
		Title      string        `json:"title"`
		Brief      string        `json:"brief"`
		Content    string        `json:"content"`
		Author     string        `json:"author"`
		AuthorID   uint32        `json:"authorid"`
		Category   string        `json:"category"`
		Tag        string        `json:"tag"`
		Image      string        `json:"image"`
		CommentNum int           `json:"commentnum"`
		VisitNum   int64         `json:"visitnum"`
		Created    string        `json:"created"`
	}
)

// AddArticle - add article.
func AddArticle(this *server.Context) error {
	var (
		reqAdd  createArticle
		imgpath string
		ip      string
	)

	if err := this.JSONBody(&reqAdd); err != nil {
		logger.Error("AddArticle json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&reqAdd); err != nil {
		logger.Error("AddArticle Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("AddArticle Can't get mysql connection:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMysql, nil)
	}

	//userID := uint32(this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64))
	userID := uint32(1000)

	err = article.CategoryService.IsExist(reqAdd.CategoryID)
	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	conpath, err := util.Save(userID, reqAdd.Content, util.Content)
	if err != nil {
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if reqAdd.Image != "" {
		imgpath, err = util.SaveImg(userID, reqAdd.Image, util.Image)
		if err != nil {
			return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
		}
	} else {
		imgpath = ""
	}

	brief := util.GetBrief(reqAdd.Content)
	ip = "http://192.168.0.107:8080"
	path := strings.Replace(conpath, ".", ip, 1)
	imgpath = strings.Replace(imgpath, ".", ip, 1)

	addArticle := &article.Article{
		Title:      reqAdd.Title,
		Brief:      brief,
		Content:    path,
		AuthorID:   userID,
		CategoryID: bson.ObjectIdHex(reqAdd.CategoryID),
		TagID:      bson.ObjectIdHex(reqAdd.TagID),
		Image:      imgpath,
		Created:    reqAdd.Created,
	}

	err = article.ArticleService.Insert(conn, addArticle)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// GetByArtID gets articles by ArtID.
func GetByArtID(this *server.Context) error {
	var (
		req struct {
			ArtID string `json:"artid"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error("GetByArtID json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	art, err := article.ArticleService.GetByArtID(req.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	reply, err := Reply(1, *art)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, reply)
}

// GetByCategoryID gets articles by CategoryID.
func GetByCategoryID(this *server.Context) error {
	var (
		req struct {
			Page       int    `json:"page"`
			CategoryID string `json:"categoryid"    validate:"required"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error("GetByCategoryID json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("GetByCategoryID Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	info, err := article.ArticleService.GetByCategoryID(req.Page, req.CategoryID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	reply, err := Reply(len(info), info...)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, reply)
}

// GetByTagID - gets articles by TagID.
func GetByTagID(this *server.Context) error {
	var (
		req struct {
			Page       int    `json:"page"`
			CategoryID string `json:"categoryid"    validate:"required"`
			TagID      string `json:"tagid"       validate:"required"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error("GetByTagID", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("GetByTagID Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	info, err := article.ArticleService.GetByTagID(req.Page, req.CategoryID, req.TagID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	reply, err := Reply(len(info), info...)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, reply)
}

// SearchByTitle - gets articles by title.
func SearchByTitle(this *server.Context) error {
	var (
		title title
	)

	if err := this.JSONBody(&title); err != nil {
		logger.Error("SearchByTitle json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	info, err := article.ArticleService.SearchByTitle(title.Title)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	reply, err := Reply(len(info), info...)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, reply)
}

// GetByUserID - gets articles by userID.
func GetByUserID(this *server.Context) error {
	var (
		userid struct {
			UserID uint32 `json:"userid"`
		}
	)

	if err := this.JSONBody(&userid); err != nil {
		logger.Error("GetByUserID json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	list, err := article.ArticleService.GetByUserID(userid.UserID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, list)
}

// DeleteArt deletes article.
func DeleteArt(this *server.Context) error {
	var (
		artid struct {
			ArtID string `json:"artid"`
		}
	)

	if err := this.JSONBody(&artid); err != nil {
		logger.Error("DeleteArt json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	//userID := uint32(this.Request().Context().Value("user").(jwtgo.MapClaims)["userid"].(float64))
	//userID := uint32(1000)
	err := article.ArticleService.Delete(artid.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	//if err := util.DeleteFile(userID, )

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// UpdateVisit update visitNum.
func UpdateVisit(this *server.Context) error {
	var (
		visit struct {
			Num   int64  `json:"num"`
			ArtID string `json:"artid"`
		}
	)

	if err := this.JSONBody(&visit); err != nil {
		logger.Error("UpdateVisit json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&visit); err != nil {
		logger.Error("UpdateVisit Validate():", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if !bson.IsObjectIdHex(visit.ArtID) {
		logger.Error(bbs.InvalidObjectId)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.UpdateVisit(visit.Num, visit.ArtID)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Recommend return popular articles.
func Recommend(this *server.Context) error {
	var (
		req struct {
			Page int `json:"page"`
		}
	)

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Recommend json", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	info, err := article.ArticleService.Recommend(req.Page)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	reply, err := Reply(len(info), info...)
	if err != nil {
		logger.Error(err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, reply)
}

// Reply return a struct of replyInfo.
func Reply(len int, info ...article.Article) ([]*replyInfo, error) {
	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't get mysql connection:", err)
		return nil, err
	}

	var reply = make([]*replyInfo, len)
	for i, art := range info {
		author, err := user.UserService.GetUserByID(conn, art.AuthorID)
		if err != nil {
			logger.Error("GetUserByID error", err)
			return nil, err
		}

		category, err := article.CategoryService.GetCategoryByID(art.CategoryID.Hex())
		if err != nil {
			logger.Error("GetCategoryByID error", err)
			return nil, err
		}

		tag, err := article.CategoryService.GetTagByID(art.CategoryID.Hex(), art.TagID.Hex())
		if err != nil {
			logger.Error("GetTagByID error", err)
			return nil, err
		}

		commentNum, err := article.CommentService.NumByArt(art.Id.Hex())
		if err != nil {
			logger.Error("NumByArt error", err)
			return nil, err
		}

		reply[i] = &replyInfo{
			Id:         art.Id,
			Title:      art.Title,
			Brief:      art.Brief,
			Content:    art.Content,
			Author:     author.UserName,
			AuthorID:   art.AuthorID,
			Category:   category.Name,
			Tag:        tag.Name,
			CommentNum: commentNum,
			VisitNum:   art.VisitNum,
			Image:      art.Image,
			Created:    art.Created,
		}

	}

	return reply, nil
}
