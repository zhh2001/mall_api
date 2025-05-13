package global

import (
	"mall_api/user_web/config"
	"mall_api/user_web/proto"
)

var (
	ServerConfig  = &config.ServerConfig{}
	NacosConfig   = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
