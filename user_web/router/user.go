package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall_api/user_web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
