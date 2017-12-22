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
 *     Initial: 2017/09/21        Jia Chenhui
 */

package command_test

import (
	"testing"

	"github.com/fengyfei/gopher/command"
)

func TestPingMacOS(t *testing.T) {
	var (
		addr string = "www.baidu.com"
		ps   command.PingSt
	)

	ps = command.PingMacOS(addr)

	if ps.SendPk != "20" || ps.RecvPk != "20" || ps.LossPk != "0.0" {
		t.Errorf("PingMacOS(%s): got ps.SendPk = %s, ps.RecvPk = %s, ps.LossPk = %s, "+
			"want ps.SendPk = %s, ps.RecvPk = %s, ps.LossPk = %s",
			addr, ps.SendPk, ps.RecvPk, ps.LossPk, "20", "20", "0.0")
	}
}
