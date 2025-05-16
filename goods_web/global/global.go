package global

import (
	ut "github.com/go-playground/universal-translator"
	"mall_api/goods_web/config"
	"mall_api/goods_web/proto"
)

var (
	Trans          ut.Translator
	ServerConfig   = &config.ServerConfig{}
	NacosConfig    = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
)
