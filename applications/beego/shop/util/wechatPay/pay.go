package wechatPay

import (
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"gopkg.in/chanxuehong/wechat.v2/mch/pay"
	"time"
	"math/rand"
	"bytes"
	"crypto/md5"
	"io"
	"fmt"
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
	now := string(time.Now().Nanosecond())
	str += now
	str += string(10000 + increasementId)

	return str
}
