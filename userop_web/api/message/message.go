package message

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mall_api/userop_web/api"
	"mall_api/userop_web/forms"
	"mall_api/userop_web/global"
	"mall_api/userop_web/models"
	"mall_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	request := &proto.MessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.MessageClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.GetTotal(),
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.GetData() {
		reMap := make(map[string]interface{})
		reMap["id"] = value.GetId()
		reMap["user_id"] = value.GetUserId()
		reMap["type"] = value.GetMessageType()
		reMap["subject"] = value.GetSubject()
		reMap["message"] = value.GetMessage()
		reMap["file"] = value.GetFile()

		result = append(result, reMap)
	}
	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.GetId(),
	})
}
