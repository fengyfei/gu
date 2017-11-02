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
	"errors"
	"time"

	"github.com/go-xorm/xorm"

	"github.com/fengyfei/gu/libs/orm"
	"github.com/fengyfei/gu/libs/security"
)

const (
	staffTable = "staff"
)

type Staff struct {
	Id        int32
	Name      string    `xorm:"varchar(30) notnull unique"`
	Pwd       string    `xorm:"varchar(128) notnull"`
	RealName  string    `xorm:"realname varchar(256) notnull unique"`
	Mobile    string    `xorm:"unique"`
	Email     string    `xorm:"varchar(80) unique"`
	CreatedAt time.Time `xorm:"created"`
	ResignAt  time.Time `xorm:"resignat"`
	Male      bool
	Active    bool
	Resigned  bool
}

// TableName returns table name in database.
func (Staff) TableName() string {
	return staffTable
}

type serviceProvider struct{}

var (
	errLoginFailed = errors.New("invalid username or password.")
	errPwdNotMatch = errors.New("old password not match.")

	// Service handles operations on model Staff.
	Service *serviceProvider
)

func init() {
	Service = &serviceProvider{}
}

// Login return user id and nil if login success.
func (sp *serviceProvider) Login(conn orm.Connection, name, pwd *string) (int32, error) {
	staff := &Staff{}

	_, err := conn.(*xorm.Engine).Where("name=?", *name).Get(staff)
	if err != nil {
		return 0, err
	}

	if !security.SaltHashCompare([]byte(staff.Pwd), pwd) {
		return 0, errLoginFailed
	}

	return staff.Id, nil
}

// Create create a new staff account.
func (sp *serviceProvider) Create(conn orm.Connection, name, pwd, realname, mobile, email *string, male bool) error {
	salt, err := security.SaltHashGenerate(pwd)
	if err != nil {
		return err
	}

	staff := &Staff{
		Name:     *name,
		Pwd:      string(salt),
		RealName: *realname,
		Mobile:   *mobile,
		Email:    *email,
		Male:     male,
		Active:   true,
	}

	_, err = conn.(*xorm.Engine).Insert(staff)

	return err
}

// Modify modify staff information.
func (sp *serviceProvider) Modify(conn orm.Connection, uid *int32, name, mobile, email *string) error {
	staff := &Staff{
		Name:   *name,
		Mobile: *mobile,
		Email:  *email,
	}

	_, err := conn.(*xorm.Engine).ID(*uid).Update(staff)

	return err
}

// ModifyPwd modify staff password.
func (sp *serviceProvider) ModifyPwd(conn orm.Connection, uid *int32, oldpwd, newpwd *string) error {
	staff := &Staff{}

	_, err := conn.(*xorm.Engine).ID(*uid).Get(staff)
	if err != nil {
		return err
	}

	if !security.SaltHashCompare([]byte(staff.Pwd), oldpwd) {
		return errPwdNotMatch
	}

	salt, err := security.SaltHashGenerate(newpwd)
	if err != nil {
		return err
	}

	update := &Staff{
		Pwd: string(salt),
	}

	_, err = conn.(*xorm.Engine).ID(*uid).Update(update)

	return err
}

// ModifyMobile modify staff mobile.
func (sp *serviceProvider) ModifyMobile(conn orm.Connection, uid *int32, mobile *string) error {
	staff := &Staff{
		Mobile: *mobile,
	}

	_, err := conn.(*xorm.Engine).ID(*uid).Update(staff)

	return err
}

// ModifyActive modify staff status.
func (sp *serviceProvider) ModifyActive(conn orm.Connection, uid *int32, active *bool) error {
	staff := &Staff{
		Active: *active,
	}

	_, err := conn.(*xorm.Engine).ID(*uid).Update(staff)

	return err
}

// Dismiss modify staff active to false and dismiss to true.
func (sp *serviceProvider) Dismiss(conn orm.Connection, uid *int32) error {
	staff := &Staff{
		Active:   false,
		Resigned: true,
		ResignAt: time.Now(),
	}

	_, err := conn.(*xorm.Engine).ID(*uid).Update(staff)

	return err
}

//IsActive return staff.Active and nil if query success
func (sp *serviceProvider) IsActive(conn orm.Connection, uid *int32) (bool, error) {
	staff := &Staff{}

	_, err := conn.(*xorm.Engine).ID(*uid).Get(staff)

	return staff.Active, err
}

// List list all on the job staff.
func (sp *serviceProvider) List(conn orm.Connection) ([]Staff, error) {
	list := []Staff{}

	_, err := conn.(*xorm.Engine).Where("resigned=?", false).Get(&list)
	if err != nil {
		return list, err
	}

	return list, nil
}

// GetByID get one staff detail information.
func (sp *serviceProvider) GetByID(conn orm.Connection, uid *int32) (*Staff, error) {
	staff := &Staff{}

	_, err := conn.(*xorm.Engine).ID(*uid).Get(staff)
	if err != nil {
		return nil, err
	}

	return staff, nil
}
