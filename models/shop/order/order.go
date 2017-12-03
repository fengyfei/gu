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
 *     Initial: 2017/11/25        ShiChao
 */

package order

import (
	"time"
	"github.com/fengyfei/gu/libs/orm"
	"github.com/jinzhu/gorm"
	"github.com/fengyfei/gu/applications/beego/shop/util/wechatPay"
	//User "github.com/fengyfei/gu/models/shop/user"
	Cart "github.com/fengyfei/gu/models/shop/cart"
	"fmt"
)

type serviceProvider struct{}

var (
	Service         *serviceProvider
	defaultParentId int32 = 0x0
	unPay           int32 = 0x0
	StatusUnpay     int32 = 0x0
	StatusPaid      int32 = 0x1
	StatusConfirmed int32 = 0x2
)

type Order struct {
	ID         int32 `gorm:"primary_key;auto_increment"`
	BillID     string
	UserID     int32
	ParentID   int32
	Status     int32
	WareId     int32
	Count      int32
	Price      float64
	ReceiveWay int8
	CreatedAt  *time.Time
}

type OrderItem struct {
	WareId int32   `json:"wareId" validate:"required"`
	Count  int32   `json:"count" validate:"required"`
	Price  float64 `json:"price" validate:"required"`
}

func (this *serviceProvider) OrderByWechat(conn orm.Connection, userId int32, IP string, receiveWay int8, orders []OrderItem) (string, error) {
	var (
		parentOrder Order
		err         error
		totalPrice  float64
		childOrders []OrderItem
		wareIdList  = make([]int32, len(orders))
		/*user        *User.User
		totalFee    int64
		paySign     string*/
	)

	childOrders = orders
	parentOrder.BillID = wechatPay.GenerateBillID()
	fmt.Println(parentOrder.BillID)
	parentOrder.UserID = userId
	parentOrder.ReceiveWay = receiveWay
	parentOrder.ParentID = defaultParentId
	parentOrder.Status = unPay

	db := conn.(*gorm.DB).Exec("USE shop")
	tx := db.Begin()

	err = tx.Create(&parentOrder).Error
	if err != nil {
		goto errFinish
	}

	for i := 0; i < len(childOrders); i++ {
		wareIdList[i] = childOrders[i].WareId
		child := &Order{}
		curOrder := childOrders[i]
		child.Price = curOrder.Price * float64(curOrder.Count)
		child.Count = curOrder.Count
		child.WareId = curOrder.WareId
		child.ParentID = parentOrder.ID

		totalPrice += child.Price
		err = tx.Create(&child).Error
		if err != nil {
			goto errFinish
		}
	}

	err = Cart.Service.RemoveWhenOrder(tx, userId, wareIdList)
	if err != nil {
		goto errFinish
	}

	/*user, err = User.Service.GetUserByID(conn, userId)
	if err != nil {
		goto errFinish
	}

	totalFee = int64(totalPrice * 100)
	paySign, err = wechatPay.OnPay(user.UserName, "desc", parentOrder.BillID, IP, totalFee)
	if err != nil {
		goto errFinish
	}

	return paySign, nil*/
	tx.Commit()
	return "", nil

errFinish:
	tx.Rollback()
	return "", err
}

func (this *serviceProvider) ChangeStateByOne(conn orm.Connection, ID, status int32) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	return db.Model(&Order{}).Where("id = ?", ID).Update("status", status).Error
}

func (this *serviceProvider) ChangeStateByGroup(conn orm.Connection, IdList []int32, status int32) error {
	var err error
	db := conn.(*gorm.DB).Exec("USE shop")
	tx := db.Begin()
	for _, id := range IdList {
		err = tx.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
		if err != nil {
			goto onErr
		}
	}

	tx.Commit()
	return nil
onErr:
	tx.Rollback()
	return err
}

func (this *serviceProvider) GetUserOrder(conn orm.Connection, userId int32) (*[]Order, error) {
	var (
		orders []Order
	)
	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func (this *serviceProvider) GetByStatus(conn orm.Connection, status int32) ([]Order, error) {
	var (
		orders []Order
	)
	db := conn.(*gorm.DB).Exec("USE shop")
	err := db.Where("status = ?", status).Find(&orders).Error
	return orders, err
}

func (this *serviceProvider) GetParents(conn orm.Connection) ([]Order, error) {
	var parents []Order
	db := conn.(*gorm.DB)
	err := db.Where("parent_id = 0").Find(&parents).Error
	return parents, err
}

func (this *serviceProvider) GetByParentId(conn orm.Connection, parentId int32) ([]Order, error) {
	var parents []Order
	db := conn.(*gorm.DB)
	err := db.Where("parent_id = ", parentId).Find(&parents).Error
	return parents, err
}
