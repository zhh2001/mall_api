package api

import (
	"context"
	"fmt"
	"mall_api/user_web/global"
	"mall_api/user_web/global/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"mall_api/user_web/proto"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	//连接用户grpc服务器
	userConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			Nickname: value.NickName,
			//Birthday: time.Unix(int64(value.Birthday), 0).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		//data["id"] = value.GetId()
		//data["name"] = value.GetNickName()
		//data["birthday"] = value.GetBirthday()
		//data["gender"] = value.GetGender()
		//data["mobile"] = value.GetMobile()

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

}
