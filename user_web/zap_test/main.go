package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction() // 生产环境
	//logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	url := "https://zhh2001.github.io/"
	sugar := logger.Sugar()
	sugar.Infow(
		"failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	logger.Info(
		"failed to fetch URL",
		zap.String("url", url),
		zap.Int("nums", 3),
		zap.Duration("backoff", time.Second),
	)
}
