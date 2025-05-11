package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"mall_api/user_web/forms"
	"mall_api/user_web/global"
)

func GenerateSmsCode(width int) string {
	//生成width长度的短信验证码
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing",
		global.ServerConfig.AliSmsInfo.ApiKey,
		global.ServerConfig.AliSmsInfo.ApiSecret,
	)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile
	request.QueryParams["SignName"] = "MALL"
	request.QueryParams["TemplateCode"] = "SMS_123456789"
	request.QueryParams["TemplateParam"] = `{"code":"` + smsCode + `"}`
	// 测试环境先不真的发送
	if false {
		response, err := client.ProcessCommonRequest(request)
		fmt.Print(client.DoAction(request, response))
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	// 将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port,
		),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Second*time.Duration(global.ServerConfig.AliSmsInfo.Expire))
	fmt.Println(smsCode)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
