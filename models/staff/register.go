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
 *     Initial: 2017/10/31        Jia Chenhui
 */

package staff

import (
	"time"

	"github.com/fengyfei/gu/libs/helper"
)

const (
	dateFormat    = "20060102"
	registerTable = "register"
)

// Register represents the sign in information.
type Register struct {
	Id           int32
	UserID       int32
	Name         string `xorm:"varchar(30) not null unique"`
	Registered   bool
	RegisteredAt time.Time
	LeaveAt      time.Time
	StayTime     time.Duration
	CreatedDate  int32
}

// RegisterOverview represents the overview of sign in information.
type RegisterOverview struct {
	RegisterAt time.Time
	StayTime   time.Duration
}

// TableName returns table name in database.
func (Register) TableName() string {
	return registerTable
}

// Register represents the first sign in.
func (sp *serviceProvider) Register(uid *int32) error {
	staff := &Staff{}
	today := todayToInt()

	_, err := Engine.ID(*uid).Get(&staff)
	if err != nil {
		return err
	}

	register := Register{
		UserID:       staff.Id,
		Name:         staff.Name,
		Registered:   true,
		RegisteredAt: time.Now(),
		CreatedDate:  today,
	}

	_, err = Engine.Insert(register)

	return err
}

// IsRegistered decide whether to sign in or not.
func (sp *serviceProvider) IsRegistered(uid *int32) (*RegisterOverview, bool, error) {
	register := &Register{}
	today := todayToInt()

	_, err := Engine.Where("id=? AND createddate=?", *uid, today).Get(register)
	if err != nil {
		return nil, false, err
	}

	r := &RegisterOverview{
		RegisterAt: register.RegisteredAt,
		StayTime:   register.StayTime,
	}

	return r, register.Registered, nil
}

// RegisterAgain sign in again.
func (sp *serviceProvider) RegisterAgain(uid *int32) error {
	today := todayToInt()

	updater := &Register{
		RegisteredAt: time.Now(),
		Registered:   true,
	}

	_, err := Engine.Where("id=? AND createddate=?", *uid, today).Update(updater)
	if err != nil {
		return err
	}

	return nil
}

// LeaveOffice represents sign out.
func (sp *serviceProvider) LeaveOffice(uid *int32, r *RegisterOverview) error {
	today := todayToInt()

	d := time.Now().Sub(r.RegisterAt) / time.Minute
	staytime := d + r.StayTime

	updater := &Register{
		LeaveAt:    time.Now(),
		Registered: false,
		StayTime:   staytime,
	}

	_, err := Engine.Where("id=? AND createddate=?", *uid, today).Update(updater)
	if err != nil {
		return err
	}

	return nil
}

func todayToInt() int32 {
	today := time.Now().Format(dateFormat)

	i, _ := helper.StrToInt32(today)

	return i
}
