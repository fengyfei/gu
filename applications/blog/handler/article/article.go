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
 *     Initial: 2017/10/27        ShiChao
 *     Modify : 2018/02/02        Tong Yuehong
 *     Modify : 2018/03/25        Chen Yanchen
 */

package article

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"strings"

	"github.com/fengyfei/gu/applications/blog/mysql"
	"github.com/fengyfei/gu/applications/blog/util"
	"github.com/fengyfei/gu/applications/core"
	"github.com/fengyfei/gu/libs/constants"
	"github.com/fengyfei/gu/libs/http/server"
	"github.com/fengyfei/gu/libs/logger"
	"github.com/fengyfei/gu/models/blog/article"
	"github.com/fengyfei/gu/models/blog/tag"
	"github.com/fengyfei/gu/models/staff"
)

// CreateArticle - insert article.
func CreateArticle(this *server.Context) error {
	var req struct {
		Title   string          `json:"title" validate:"required,max=32"`
		Content string          `json:"content" validate:"required"`
		TagsID  []bson.ObjectId `json:"tagsid"`
		Image   string          `json:"image"`
	}

	err := this.JSONBody(&req)
	if err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	AuthorID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	if req.Image != "" {
		req.Image, err = util.SaveImage(req.Image, "article/")
		if err != nil {
			logger.Error("Save image failed:", err)
			return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
		}
	}
	/*
		addr, err := net.InterfaceAddrs()
		for _, address := range addr {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
				}
			}
		}
	*/
	ip := "http://192.168.0.102:21002"
	imgpath := strings.Replace(req.Image, "./file", ip, 1)
	a := &article.Article{
		AuthorId: AuthorID,
		Title:    req.Title,
		Brief:    article.ArticleService.GetBrief(req.Content),
		Content:  req.Content,
		TagsID:   req.TagsID,
		Image:    imgpath,
	}

	id, err := article.ArticleService.Create(a)
	if err != nil {
		logger.Error("Create false:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	for i, _ := range req.TagsID {
		num, err := article.ArticleService.CountByTag(req.TagsID[i])
		if err != nil {
			return err
		}
		tag.TagService.UpdateCount(req.TagsID[i], num)
	}

	return core.WriteStatusAndIDJSON(this, constants.ErrSucceed, id)
}

// ArticleByID return article by articleID.
func ArticleByID(this *server.Context) error {
	var req struct {
		ID string `json:"id" validate:"required,len=24"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	a, err := article.ArticleService.GetByID(req.ID)
	if err != nil {
		logger.Error("Not found article:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}
	resp, err := replyArticle(&a)
	if err != nil {
		logger.Error(err)
		return err
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, &resp)
}

// ListCreated return articles which are waiting for checking.
func ListCreated(this *server.Context) error {
	var req struct {
		Page int `json:"page"`
	}
	var resp []completeArticle

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
	}

	articles, err := article.ArticleService.ListCreated(req.Page)
	if err != nil {
		logger.Error("Not found articles:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	for _, v := range articles {
		a, err := replyArticle(&v)
		if err != nil {
			return err
		}
		resp = append(resp, *a)
	}
	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, resp)
}

// ListApproval returns the articles which are passed.
func ListApproval(this *server.Context) error {
	var page struct {
		Page int `json:"page"`
	}

	if err := this.JSONBody(&page); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.ListApproval(page.Page)
	if err != nil {
		logger.Error("Not found articles:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	resp := make([]*simpleArticle, len(articles))
	for i, _ := range articles {
		a, err := replyArt(&articles[i])
		if err != nil {
			logger.Error("Error in ListApproval:", err)
			return err
		}
		resp[i] = a
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, resp)
}

// ModifyStatus modify the article status.
func ModifyStatus(this *server.Context) error {
	var req struct {
		ArticleID string `json:"aid" validate:"required"`
		Status    int8   `json:"status"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}
	StaffId := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.ModifyStatus(req.ArticleID, req.Status, StaffId)
	if err != nil {
		logger.Error("Modify status false:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// Delete delete article.
func Delete(this *server.Context) error {
	var req struct {
		ArticleID string `json:"aid" validate:"required,len=24"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	StaffID := int32(this.Request().Context().Value("staff").(jwtgo.MapClaims)["staffid"].(float64))

	err := article.ArticleService.Delete(req.ArticleID, StaffID)
	if err != nil {
		logger.Error("Delete article false:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// UpdateView update article's view.
func UpdateView(this *server.Context) error {
	var req struct {
		ArticleID string `json:"aid"`
		Views     uint64 `json:"views"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	if err := this.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := article.ArticleService.UpdateView(&req.ArticleID, req.Views)
	if err != nil {
		logger.Error("Update views false:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// ModifyArticle modify article.
func ModifyArticle(this *server.Context) error {
	var req struct {
		ArticleID string          `json:"aid" validate:"required"`
		Title     string          `json:"title" validate:"required,max=32"`
		Content   string          `json:"content" validate:"required"`
		TagsID    []bson.ObjectId `json:"tagsid"`
		Image     string          `json:"image"`
	}

	if err := this.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	err := this.Validate(&req)
	if err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	req.Image, err = util.SaveImage(req.Image, "article/")
	if err != nil {
		logger.Error("Save image failed:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrInvalidParam, nil)
	}

	a := &article.Article{
		Title:   req.Title,
		Brief:   article.ArticleService.GetBrief(req.Content),
		Content: req.Content,
		TagsID:  req.TagsID,
		Image:   req.Image,
	}

	err = article.ArticleService.ModifyArticle(req.ArticleID, a)
	if err != nil {
		logger.Error("Modify article false:", err)
		return core.WriteStatusAndDataJSON(this, constants.ErrMongoDB, nil)
	}

	return core.WriteStatusAndDataJSON(this, constants.ErrSucceed, nil)
}

// GetByTag get article by tag.
func GetByTag(c *server.Context) error {
	var req struct {
		TagId string `json:"id" validate:"required"`
		Page int `json:"page"`
	}

	if err := c.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.GetByTagId(req.TagId, req.Page)
	if err != nil {
		logger.Error("Get articles false:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp := make([]simpleArticle, len(articles))
	for i, _ := range articles {
		a, err := replyArt(&articles[i])
		if err != nil {
			logger.Error("Can't reply article:", err)
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}
		resp[i] = *a
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

// GetByAuthorID
func GetByAuthorID(c *server.Context) error {
	var req struct {
		AuthorID int32 `json:"id"`
	}

	if err := c.JSONBody(&req); err != nil {
		logger.Error("Request error:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if err := c.Validate(&req); err != nil {
		logger.Error("Validate error:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	articles, err := article.ArticleService.GetByAuthorID(req.AuthorID)
	if err != nil {
		logger.Error("Get articles false:", err)
		return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
	}

	resp := make([]simpleArticle, len(articles))
	for i, _ := range articles {
		a, err := replyArt(&articles[i])
		if err != nil {
			logger.Error("Can't reply article:", err)
			return core.WriteStatusAndDataJSON(c, constants.ErrMongoDB, nil)
		}
		resp[i] = *a
	}

	return core.WriteStatusAndDataJSON(c, constants.ErrSucceed, resp)
}

type completeArticle struct {
	ID       string
	AuthorId int32
	Author   string
	Title    string
	Content  string
	Image    string
	TagsId   []bson.ObjectId
	Tags     []string
	Views    uint64
	Created  string
	Updated  string
	Status   int8
}

// replyArticle return the article details.
func replyArticle(a *article.Article) (*completeArticle, error) {
	var tags []string

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't connect mysql:", err)
		return nil, err
	}

	author, err := staff.Service.GetByID(conn, a.AuthorId)
	if err != nil {
		logger.Error("Can't find author name:", err)
	}

	for _, v := range a.TagsID {
		tid := v.Hex()
		t, err := tag.TagService.GetByID(&tid)
		if err != nil {
			logger.Error("Can't find tag:", err)
			return nil, err
		}
		tags = append(tags, t.Tag)
	}
	art := &completeArticle{
		ID:       a.ID.Hex(),
		AuthorId: a.AuthorId,
		Author:   author.Name,
		Title:    a.Title,
		Content:  a.Content,
		Image:    a.Image,
		TagsId:   a.TagsID,
		Tags:     tags,
		Views:    a.Views,
		Created:  a.Created.String(),
	}
	return art, nil
}

type simpleArticle struct {
	Id      string
	Title   string
	Author  string
	Brief   string
	Tags    []string
	Views   uint64
	Created string
	Image   string
}

func replyArt(a *article.Article) (*simpleArticle, error) {
	var tags []string

	conn, err := mysql.Pool.Get()
	defer mysql.Pool.Release(conn)
	if err != nil {
		logger.Error("Can't connect mysql:", err)
		return nil, err
	}

	author, err := staff.Service.GetByID(conn, a.AuthorId)
	if err != nil {
		logger.Error("Can't find author name:", err)
	}

	for i, _ := range a.TagsID {
		tid := a.TagsID[i].Hex()
		t, err := tag.TagService.GetByID(&tid)
		if err != nil {
			logger.Error("Can't find tag:", err)
			return nil, err
		}
		tags = append(tags, t.Tag)
	}
	art := &simpleArticle{
		Id:      a.ID.Hex(),
		Title:   a.Title,
		Author:  author.Name,
		Brief:   a.Brief,
		Image:   a.Image,
		Tags:    tags,
		Views:   a.Views,
		Created: a.Created.String(),
	}
	return art, nil
}
