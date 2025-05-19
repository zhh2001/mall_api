package router

import (
	"github.com/gin-gonic/gin"

	"mall_api/order_web/api/order"
	"mall_api/user_web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("goods")
	{
		OrderRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), order.List) // 订单列表
		OrderRouter.POST("", middlewares.JWTAuth(), order.New)                            // 新建订单
		OrderRouter.GET("/:id", middlewares.JWTAuth(), order.Detail)                      // 订单详情
	}
}
