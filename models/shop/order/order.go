/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
	// "github.com/fengyfei/gu/applications/beego/shop/util/wechatPay"
	//User "github.com/fengyfei/gu/models/shop/user"
	Cart "github.com/fengyfei/gu/models/shop/cart"
	"fmt"
	"github.com/fengyfei/gu/applications/shop/util/wechatPay"
)

type serviceProvider struct{}

var (
	Service         *serviceProvider
	defaultParentId uint  = 0x0
	unPay           uint  = 0x0
	StatusUnpay     uint  = 0x0
	StatusPaid      uint  = 0x1
	StatusConfirmed uint  = 0x2
	PayWayOnline    int32 = 0x0
)

const AMonth = 30 * 24 * 60 * 60 * 1e9

type Order struct {
	ID         uint    `gorm:"primary_key;auto_increment"`
	BillID     string  `gorm:"not null"`
	UserID     uint    `gorm:"not null"`
	ParentID   uint    `gorm:"not null"`
	Status     uint    `gorm:"not null";default:0`
	WareId     uint    `gorm:"not null"`
	Count      uint    `gorm:"not null";default:0`
	Price      float64 `gorm:"not null";default:0`
	CreatedAt  *time.Time
}

type OrderItem struct {
	WareId uint    `json:"wareId" validate:"required"`
	Count  uint    `json:"count" validate:"required"`
	Price  float64 `json:"price" validate:"required"`
}

func (this *serviceProvider) OrderByWechat(conn orm.Connection, userId uint, IP string, receiveWay int8, orders []OrderItem) (string, error) {
	var (
		parentOrder Order
		err         error
		totalPrice  float64
		childOrders []OrderItem
		wareIdList  = make([]uint, len(orders))
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

func (this *serviceProvider) ChangeStateByOne(conn orm.Connection, ID, status uint) error {
	db := conn.(*gorm.DB).Exec("USE shop")
	return db.Model(&Order{}).Where("id = ?", ID).Update("status", status).Error
}

func (this *serviceProvider) ChangeStateByGroup(conn orm.Connection, IdList []uint, status uint) error {
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

func (this *serviceProvider) GetUserOrder(conn orm.Connection, userId uint) (*[]Order, error) {
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

func (this *serviceProvider) GetByStatus(conn orm.Connection, status uint) ([]Order, error) {
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

func (this *serviceProvider) GetByParentId(conn orm.Connection, parentId uint) ([]Order, error) {
	var parents []Order
	db := conn.(*gorm.DB)
	err := db.Where("parent_id = ", parentId).Find(&parents).Error
	return parents, err
}

func (this *serviceProvider) ClearUnpaidOrder(conn orm.Connection) error {
	db := conn.(*gorm.DB)
	tx := db.Begin()
	aMonthAgo := time.Now().Add(-AMonth)

	err := tx.Where("status = ? && created_at < ?", StatusUnpay, aMonthAgo).Delete(&Order{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
