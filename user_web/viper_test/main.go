package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func main() {
	debug := GetEnvInfo("MALL_DEBUG")
	configFilePrefix := "config"
	configFilename := fmt.Sprintf("./user_web/viper_test/%s-prod.yaml", configFilePrefix)
	if debug {
		configFilename = fmt.Sprintf("./user_web/viper_test/%s-dev.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFilename)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	//动态监控配置变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&serverConfig)
		fmt.Println(serverConfig)
	})

	time.Sleep(1000 * time.Second)
}
