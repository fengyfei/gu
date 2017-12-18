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
 *     Initial: 2017/12/17        ShiChao
 */

package security

import (
	"testing"
	"net"
	"github.com/fengyfei/gu/libs/logger"
)

func TestFindArray(t *testing.T) {
	arr1 := []string{"10.0.0.1", "192.0.2.0/24"}

	ip, err := NewIP(arr1, false)
	if err != nil {
		t.Error("parse err: ", err)
	}
	logger.Info("ip: %+v \n", *ip.whiteListsIPs[0])

	arr2 := []string{"12.54.66.68.55"}
	ip, err = NewIP(arr2, false)
	if err != nil {
		t.Error("parse err: ", err)
	}

}
func TestIP_Contains(t *testing.T) {
	arr := []string{"10.0.0.1", "192.0.2.0/24", "12.54.66.68.55"}
	ip, err := NewIP(arr, false)
	if err != nil {
		t.Error(err)
		return
	}

	ipaddr1 := "10.0.0.1"
	ipaddr2 := "10.0.0.2"
	ipaddr3 := "192.0.2.4"

	exist, _, _ := ip.Contains(ipaddr1)
	if !exist {
		t.Error("10.0.0.1 should be contained")
	}

	exist, _, _ = ip.Contains(ipaddr2)
	if exist {
		t.Error("10.0.0.2 shouldn't be contained")
	}

	exist, _, _ = ip.Contains(ipaddr3)
	if !exist {
		t.Error("192.0.2.4 should be contained")
	}
}

func TestIP_ContainsIP(t *testing.T) {
	arr := []string{"10.0.0.1", "192.0.2.0/24"}
	ip, err := NewIP(arr, false)
	if err != nil {
		t.Error(err)
		return
	}

	ipaddr1 := net.ParseIP("10.0.0.1")
	ipaddr2 := net.ParseIP("10.0.0.2")
	ipaddr3 := net.ParseIP("192.0.2.4")
	ipaddr4 := net.ParseIP("10.0.0.2333")

	exist := ip.ContainsIP(ipaddr1)
	if !exist {
		t.Error("10.0.0.1 should be contained")
	}

	exist = ip.ContainsIP(ipaddr2)
	if exist {
		t.Error("10.0.0.1 shouldn't be contained")
	}

	exist = ip.ContainsIP(ipaddr3)
	if !exist {
		t.Error("192.0.2.4 should be contained")
	}

	exist = ip.ContainsIP(ipaddr4)
	if exist {
		t.Error("10.0.0.2333 shouldn't be contained")
	}
}
