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

	"github.com/fengyfei/gu/applications/echo/core/orm"
)

// ContactInfo is the more detail of one particular staff.
type ContactInfo struct {
	Id       int32
	Name     string
	RealName string
	Mobile   string
	Email    string
	HireAt   time.Time
	Male     bool
}

// ContactOverview is the more detail of one particular staff.
type ContactOverview struct {
	UserID   int32
	RealName string
}

// OverviewList list all on the job staff.
func (sp *serviceProvider) OverviewList() ([]ContactOverview, error) {
	slist := []Staff{}
	clist := []ContactOverview{}

	_, err := orm.Engine.Where("resigned=?", false).Get(&slist)
	if err != nil {
		return clist, err
	}

	for _, s := range slist {
		c := ContactOverview{
			UserID:   s.Id,
			RealName: s.RealName,
		}

		clist = append(clist, c)
	}

	return clist, nil
}

// InfoList get all staff detail information.
func (sp *serviceProvider) InfoList() ([]ContactInfo, error) {
	slist := []Staff{}
	clist := &ContactInfo{}

	_, err := orm.Engine.Where("resigned=?", false).Get(&slist)
	if err != nil {
		return clist, err
	}

	for _, s := range slist {
		c := ContactInfo{
			Id:       s.Id,
			Name:     s.Name,
			RealName: s.RealName,
			Mobile:   s.Mobile,
			Email:    s.Email,
			HireAt:   s.HireAt,
			Male:     s.Male,
		}

		clist = append(clist, c)
	}

	return clist, nil
}

// GetByID get one staff detail information.
func (sp *serviceProvider) GetByID(uid *int32) (*ContactInfo, error) {
	staff := &Staff{}
	contact := &ContactInfo{}

	_, err := orm.Engine.ID(*uid).Get(staff)
	if err != nil {
		return nil, err
	}

	contact = &ContactInfo{
		Name:     staff.Name,
		RealName: staff.RealName,
		Mobile:   staff.Mobile,
		Email:    staff.Email,
		HireAt:   staff.HireAt,
		Male:     staff.Male,
	}

	return contact, nil
}
