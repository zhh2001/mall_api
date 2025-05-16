package global

import (
	ut "github.com/go-playground/universal-translator"

	"mall_api/oss_web/config"
)

var (
	Trans        ut.Translator
	ServerConfig = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}
)
