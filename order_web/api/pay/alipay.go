package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
)

func main() {
	appID := "app_id"
	privateKey := "private_key"
	aliPublicKey := "ali_public_key"
	var client, err = alipay.New(appID, privateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://10.120.221.149:8023/o/v1/pay/alipay/notify"
	p.ReturnURL = "http://10.120.221.149:8089/o/v1/pay/alipay/return"
	p.Subject = "MALL-订单支付"
	p.OutTradeNo = "zhang_sues"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(url.String())
}
