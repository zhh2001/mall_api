package main

import (
	"fmt"
	
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func main() {
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "xxxx", "xxx")
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = "16600008888"
	request.QueryParams["SignName"] = "MALL"
	request.QueryParams["TemplateCode"] = "SMS_123456789"
	request.QueryParams["TemplateParam"] = `{"code":"777777"}`
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
