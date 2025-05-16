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
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)          // 该接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                      // 获取商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) // 删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                               // 获取商品的库存

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
}
