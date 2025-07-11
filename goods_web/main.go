package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mall_api/goods_web/global"
	"mall_api/goods_web/initialize"
	"mall_api/goods_web/utils"
	"mall_api/goods_web/utils/register/consul"
)

func main() {
	//1.初始化logger
	initialize.InitLogger()

	//2.初始化配置文件
	initialize.InitConfig()

	//3.初始化routers
	Router := initialize.Routers()

	//4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5.初始化srv的连接
	initialize.InitSrvConn()

	//6.初始化Sentinel
	initialize.InitSentinel()

	viper.AutomaticEnv()
	// 如果是本地开发环境，端口号固定
	debug := viper.GetBool("MALL_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	registerClient := consul.NewRegistryClient(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
	)
	serviceId := uuid.NewV4().String()
	err := registerClient.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceId,
	)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}

	/*
		1. S()可以获取一个全局的Sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug，info，warn，error，fetal
		3. S函数和L函数很有用，提供了一个全局的安全访问logger的途径
	*/
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.Deregister(serviceId); err != nil {
		zap.S().Panic("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
