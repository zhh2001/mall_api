package router

import (
	"github.com/gin-gonic/gin"

	"mall_api/goods_web/api/goods"
	"mall_api/goods_web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) // 该接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                             // 获取商品详情
	}
}
