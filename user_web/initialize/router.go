package initialize

import (
	"github.com/gin-gonic/gin"
	router2 "mall_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(ApiGroup)
	return Router
}
