package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_api/order_web/global"
	"mall_api/order_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.GoodsSrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【商品服务】失败")
		return
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	orderConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.OrderSrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【订单服务】失败")
		return
	}

	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	invConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.InventorySrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【库存服务】失败")
		return
	}

	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}
