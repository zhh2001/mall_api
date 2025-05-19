package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mall_api/userop_web/global"
	"mall_api/userop_web/proto"
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

	userOpConn, err := grpc.NewClient(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			consulInfo.Host,
			consulInfo.Port,
			global.ServerConfig.UserOpSrvInfo.Name,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【用户操作服务】失败")
		return
	}

	global.UserFavClient = proto.NewUserFavClient(userOpConn)
	global.MessageClient = proto.NewMessageClient(userOpConn)
	global.AddressClient = proto.NewAddressClient(userOpConn)
}
