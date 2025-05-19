package user_fav

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mall_api/userop_web/api"
	"mall_api/userop_web/forms"
	"mall_api/userop_web/global"
	"mall_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	userFavRsp, err := global.UserFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("获取收藏列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range userFavRsp.GetData() {
		ids = append(ids, item.GetGoodsId())
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	//请求商品服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": userFavRsp.GetData(),
	}

	goodsList := make([]interface{}, 0)
	for _, item := range userFavRsp.GetData() {
		data := gin.H{
			"id": item.GetGoodsId(),
		}

		for _, good := range goods.GetData() {
			if item.GetGoodsId() == good.GetId() {
				data["name"] = good.GetName()
				data["shop_price"] = good.GetShopPrice()
			}
		}

		goodsList = append(goodsList, data)
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userFavForm := forms.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	//缺少一步， 这个时候应该去商品服务查询一下这个是否存在
	userId, _ := ctx.Get("userId")
	_, err := global.UserFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})

	if err != nil {
		zap.S().Errorw("添加收藏记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除收藏记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Detail(ctx *gin.Context) {
	goodsId := ctx.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.UserFavClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("查询收藏状态失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
