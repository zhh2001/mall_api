package initialize

import (
	"github.com/gin-gonic/gin"

	"mall_api/user_web/middlewares"
	router2 "mall_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(ApiGroup)
	router2.InitBaseRouter(ApiGroup)

	return Router
}
