package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"mall_api/order_web/api"
	"mall_api/order_web/forms"
	"mall_api/order_web/global"
	"mall_api/order_web/models"
	"mall_api/order_web/proto"
)

func List(ctx *gin.Context) {
	// 订单的列表
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	// 如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	request.Pages = int32(pagesInt)
	request.PagePerNums = int32(perNumsInt)

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	/*
		{
			"total":100,
			"data":[
				{
					"":
				}
			]
		}
	*/
	reMap := gin.H{
		"total": rsp.GetTotal(),
	}
	orderList := make([]interface{}, 0)

	for _, item := range rsp.GetData() {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}
	userId, _ := ctx.Get("userId")
	orderRequest := proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	}
	rsp, err := global.OrderSrvClient.CreateOrder(context.Background(), &orderRequest)
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成支付宝的支付URL
	client, err := alipay.New(
		global.ServerConfig.AlipayInfo.AppID,
		global.ServerConfig.AlipayInfo.PrivateKey,
		false,
	)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AlipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "MALL-订单-" + rsp.GetOrderSn()
	p.OutTradeNo = rsp.GetOrderSn()
	p.TotalAmount = strconv.FormatFloat(float64(rsp.GetTotal()), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.GetId(),
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	// 如果是管理员用户则返回所有的订单
	request := proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	orderInfo := rsp.GetOrderInfo()
	reMap := gin.H{}
	reMap["id"] = orderInfo.GetId()
	reMap["status"] = orderInfo.GetStatus()
	reMap["user"] = orderInfo.GetUserId()
	reMap["post"] = orderInfo.GetPost()
	reMap["total"] = orderInfo.GetTotal()
	reMap["address"] = orderInfo.GetAddress()
	reMap["name"] = orderInfo.GetName()
	reMap["mobile"] = orderInfo.GetMobile()
	reMap["pay_type"] = orderInfo.GetPayType()
	reMap["order_sn"] = orderInfo.GetOrderSn()

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.GetGoods() {
		tmpMap := gin.H{
			"id":    item.GetGoodsId(),
			"name":  item.GetGoodsName(),
			"image": item.GetGoodsImage(),
			"price": item.GetGoodsPrice(),
			"nums":  item.GetNums(),
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	// 生成支付宝的支付URL
	client, err := alipay.New(
		global.ServerConfig.AlipayInfo.AppID,
		global.ServerConfig.AlipayInfo.PrivateKey,
		false,
	)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AlipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "MALL订单-" + orderInfo.GetOrderSn()
	p.OutTradeNo = orderInfo.GetOrderSn()
	p.TotalAmount = strconv.FormatFloat(float64(orderInfo.GetTotal()), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付URL失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	reMap["alipay_url"] = url.String()

	ctx.JSON(http.StatusOK, reMap)
}
