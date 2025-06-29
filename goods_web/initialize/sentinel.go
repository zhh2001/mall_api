package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化 Sentinel 异常：%v", err)
	}

	// 配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 匀速通过
			MaxQueueingTimeMs:      100,             // 匀速排队的最大等待时间，该字段仅仅对 `Throttling` ControlBehavior生效
			Threshold:              100,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		zap.S().Fatalf("加载规则失败：%v", err)
		return
	}
}
