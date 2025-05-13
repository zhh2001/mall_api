package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"mall_api/user_web/config"
)

func main() {
	sc := []constant.ServerConfig{
		{
			IpAddr:      "10.120.21.77",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "73bb8093-63eb-4bf3-bd81-84b18fe259d6",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev"},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	serverConfig := config.ServerConfig{}
	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件产生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	time.Sleep(3000 * time.Second)
}
