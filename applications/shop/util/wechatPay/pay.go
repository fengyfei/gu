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

package wechatPay

import (
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"time"
	"math/rand"
	"bytes"
	"crypto/md5"
	"io"
	"fmt"
	"strconv"
)

var (
	clt            *core.Client
	increasementId = 0
	tradeType      = "JSAPI"
	notifyURL      = ""
	appId          = ""
	mchId          = ""
	apiKey         = ""
)

func init() {
	clt = core.NewClient(appId, mchId, apiKey, nil)
}

func OnPay(openId, body, billId, IP string, totalFee int64) (string, error) {
	var (
		req      *pay.UnifiedOrderRequest
		resp     *pay.UnifiedOrderResponse
		err      error
		nonceStr string
	)

	nonceStr = RandomString(30)
	req.OpenId = openId
	req.Body = body
	req.OutTradeNo = billId
	req.SpbillCreateIP = IP
	req.NotifyURL = notifyURL
	req.TradeType = tradeType
	req.TotalFee = totalFee

	resp, err = pay.UnifiedOrder2(clt, req)
	if err != nil {
		return "", err
	}

	now := string(time.Now().Nanosecond())

	str := fmt.Sprintf("appId=%s&nonceStr=%s&package=prepay_id=%s&signType=MD5&timeStamp=%s&key=%s", appId, nonceStr, resp.PrepayId, now, apiKey)
	w := md5.New()
	io.WriteString(w, str)
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))

	return md5str2, nil
}

func RandomString(randLength int) (result string) {
	var num string = "0123456789"
	var lower string = "abcdefghijklmnopqrstuvwxyz"
	var upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := bytes.Buffer{}
	b.WriteString(num)
	b.WriteString(lower)
	b.WriteString(upper)
	var str = b.String()
	var strLen = len(str)
	if strLen == 0 {
		result = ""
		return
	}

	rand.Seed(time.Now().UnixNano())
	b = bytes.Buffer{}
	for i := 0; i < randLength; i++ {
		b.WriteByte(str[rand.Intn(strLen)])
	}
	result = b.String()
	return
}

func GenerateBillID() string{
	str := "shop"
	now := strconv.Itoa(time.Now().Nanosecond())
	str += now
	str += strconv.Itoa(10000 + increasementId)
	increasementId++
	fmt.Println(str, now, increasementId)

	return str
}
