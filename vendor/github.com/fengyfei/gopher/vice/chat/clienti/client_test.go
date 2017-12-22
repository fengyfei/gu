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

package client_test

import (
	"testing"
	"time"
	"fmt"

	"github.com/fengyfei/gopher/vice/chat/clienti"
	"github.com/fengyfei/gopher/vice/chat/common"
)

func BenchmarkSendToII(b *testing.B) {
	msg := message.Msg{
		From:    "I",
		To:      "II",
		Content: "send to II from I",
	}

	count := 0
	timer := time.NewTimer(time.Second)

loop:
	for i := 0; i < b.N; i++ {
		select {
		case <-timer.C:
			break loop
		default:
		}

		go client.SendToII(msg)
		count = count + 1
	}

	fmt.Println("count --->>>", count)
}
